package booking

type MemoryStore struct {
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) (Booking, error) {
	if _, exists := s.bookings[b.SeatID]; exists {
		return b, ErrSeatAlreadyBooked
	}
	s.bookings[b.SeatID] = b
	return b, nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var result []Booking
	for _, booking := range s.bookings {
		if booking.MovieID == movieID {
			result = append(result, booking)
		}
	}
	return result
}
