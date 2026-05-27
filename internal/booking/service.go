package booking

type Service struct {
	store BookingStore
}

func (s *Service) HoldSeat(b Booking) (Booking, error) {
	return s.store.Hold(b)
}

func (s *Service) ConfirmSession(sessionID string) error {
	_, err := s.store.ConfirmSession(sessionID)
	return err
}

func (s *Service) CancelSession(sessionID string) error {
	_, err := s.store.CancelSession(sessionID)
	return err
}

func NewService(store BookingStore) *Service {
	return &Service{store}
}

func (s *Service) ListBookings(movieID string) ([]Booking, error) {
	return s.store.ListBookings(movieID)
}
