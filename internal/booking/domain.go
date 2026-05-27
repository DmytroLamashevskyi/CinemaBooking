package booking

import (
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked = errors.New("seat already booked")
	ErrSessionNotFound   = errors.New("session not found or expired")
	ErrSeatNotFound      = errors.New("seat reservation not found or expired")
)

type SeatsInfo struct {
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Booked    bool   `json:"booked"`
	Confirmed bool   `json:"confirmed"`
}

// Booking represents a confirmed seat reservation.
type Booking struct {
	ID        string    `json:"id"`
	MovieID   string    `json:"movie_id"`
	SeatID    string    `json:"seat_id"`
	UserID    string    `json:"user_id"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}

type BookingStore interface {
	ListBookings(movieID string) ([]Booking, error)
	Hold(b Booking) (Booking, error)
	ConfirmSession(sessionID string) (Booking, error)
	CancelSession(sessionID string) (Booking, error)
}
