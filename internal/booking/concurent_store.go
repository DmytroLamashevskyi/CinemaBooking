package booking

import "sync"

type ConcurentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurentStore() *ConcurentStore {
	return &ConcurentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurentStore) Book(b Booking) (Booking, error) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.bookings[b.SeatID]; exists {
		return b, ErrSeatAlreadyBooked
	}
	s.bookings[b.SeatID] = b
	return b, nil
}

func (s *ConcurentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock()
	var result []Booking
	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			result = append(result, booking)
		}
	}
	return result
}
