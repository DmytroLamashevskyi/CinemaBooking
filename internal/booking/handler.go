package booking

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/DmytroLamashevskyi/CinemaBooking/internal/utils"
)

type handler struct {
	svc *Service
}

func NewHandler(svc *Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings, err := h.svc.ListBookings(movieID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list seats"})
		return
	}

	seats := make([]SeatsInfo, 0, len(bookings))
	for _, booking := range bookings {
		seats = append(seats, SeatsInfo{
			SeatID:    booking.SeatID,
			UserID:    booking.UserID,
			Booked:    true,
			Confirmed: booking.Status == "confirmed",
		})
	}

	utils.WriteJSON(w, http.StatusOK, seats)
}

func (h *handler) HoldSeat(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")

	type holdRequest struct {
		UserID string `json:"user_id"`
	}

	var req holdRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	booking, err := h.svc.HoldSeat(Booking{
		MovieID: movieID,
		SeatID:  seatID,
		UserID:  req.UserID,
	})
	if err != nil {
		if errors.Is(err, ErrSeatAlreadyBooked) {
			utils.WriteJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hold seat"})
		return
	}

	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movie_id"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}

	utils.WriteJSON(w, http.StatusOK, holdResponse{
		SessionID: booking.ID,
		MovieID:   booking.MovieID,
		SeatID:    booking.SeatID,
		ExpiresAt: booking.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	err := h.svc.ConfirmSession(sessionID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrSeatNotFound) {
			utils.WriteJSON(w, http.StatusGone, map[string]string{"error": err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to confirm session"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *handler) CancelSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.PathValue("sessionID")

	err := h.svc.CancelSession(sessionID)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) || errors.Is(err, ErrSeatNotFound) {
			utils.WriteJSON(w, http.StatusGone, map[string]string{"error": err.Error()})
			return
		}
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to cancel session"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
