package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultTTL = 2 * time.Minute

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (s *RedisStore) Book(b Booking) (Booking, error) {
	session, err := s.hold(b)
	if err != nil {
		return b, err
	}
	return session, nil
}

func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	b.ID = id
	val, _ := json.Marshal(b)

	result := s.client.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultTTL,
	})
	if result.Val() != "OK" {
		return Booking{}, ErrSeatAlreadyBooked
	}

	s.client.Set(ctx, sessionKey(id), key, defaultTTL)

	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultTTL),
	}, nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	return []Booking{}
}
