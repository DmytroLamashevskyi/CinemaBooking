package main

import (
	"log"
	"net/http"

	"github.com/DmytroLamashevskyi/CinemaBooking/internal/adapters/redis"
	"github.com/DmytroLamashevskyi/CinemaBooking/internal/booking"
	"github.com/DmytroLamashevskyi/CinemaBooking/internal/utils"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", listMoviesHandler)

	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)
	bookingHandler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/confirm", bookingHandler.ConfirmSeat)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/cancel", bookingHandler.CancelSeat)
	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.CancelSession)

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	if err := http.ListenAndServe(":8088", mux); err != nil {
		log.Fatal(err)
	}
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

var movies = []movieResponse{
	{ID: "matrix", Title: "The Matrix", Rows: 5, SeatsPerRow: 10},
	{ID: "interstellar", Title: "Interstellar", Rows: 5, SeatsPerRow: 8},
	{ID: "parasite", Title: "Parasite", Rows: 6, SeatsPerRow: 8},
}

func listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}
