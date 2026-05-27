package main

import (
	"log"
	"net/http"

	"github.com/DmytroLamashevskyi/CinemaBooking/internal/utils"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", listMoviesHandler)
	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	if err := http.ListenAndServe(":8088", mux); err != nil {
		log.Fatal(err)
	}
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seatsPerRow"`
}

var movies = []movieResponse{
	{ID: "matrix", Title: "The Matrix", Rows: 5, SeatsPerRow: 10},
	{ID: "interstellar", Title: "Interstellar", Rows: 5, SeatsPerRow: 8},
	{ID: "parasite", Title: "Parasite", Rows: 6, SeatsPerRow: 8},
}

func listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

func bookSeatHandler(w http.ResponseWriter, r *http.Request) {

}
