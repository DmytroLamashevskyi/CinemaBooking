package booking

type Service struct {
	store BookingStore
}

func (s *Service) CancelSeat(movieID string, seatID string) error {
	_, err := s.store.Cancel(movieID, seatID)
	return err
}

func (s *Service) ConfirmSeat(movieID string, seatID string) error {
	_, err := s.store.Confirm(movieID, seatID)
	return err
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

func (s *Service) Book(booking Booking) error {
	_, err := s.store.Book(booking)
	return err
}

func (s *Service) ListBookings(movieID string) ([]Booking, error) {
	return s.store.ListBookings(movieID)
}
