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

func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	
	b.ID = id
	b.Status = "held"
	b.ExpiresAt = now.Add(defaultTTL)
	
	val, err := json.Marshal(b)
	if err != nil {
		return Booking{}, err
	}

	cmd := s.client.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultTTL,
	})
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			return Booking{}, ErrSeatAlreadyBooked
		}
		return Booking{}, err
	}
	if cmd.Val() != "OK" {
		return Booking{}, ErrSeatAlreadyBooked
	}

	if err := s.client.Set(ctx, sessionKey(id), key, defaultTTL).Err(); err != nil {
		s.client.Del(ctx, key)
		return Booking{}, err
	}

	return b, nil
}

func (s *RedisStore) Hold(b Booking) (Booking, error) {
	return s.hold(b)
}

func (s *RedisStore) ConfirmSession(sessionID string) (Booking, error) {
	ctx := context.Background()
	seatKey, err := s.client.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return Booking{}, ErrSessionNotFound
		}
		return Booking{}, err
	}

	val, err := s.client.Get(ctx, seatKey).Result()
	if err != nil {
		if err == redis.Nil {
			return Booking{}, ErrSeatNotFound
		}
		return Booking{}, err
	}

	var b Booking
	if err := json.Unmarshal([]byte(val), &b); err != nil {
		return Booking{}, err
	}

	b.Status = "confirmed"
	newVal, err := json.Marshal(b)
	if err != nil {
		return Booking{}, err
	}
	err = s.client.Set(ctx, seatKey, newVal, 0).Err() // persist
	if err != nil {
		return Booking{}, err
	}

	s.client.Del(ctx, sessionKey(sessionID))
	return b, nil
}

func (s *RedisStore) CancelSession(sessionID string) (Booking, error) {
	ctx := context.Background()
	seatKey, err := s.client.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return Booking{}, ErrSessionNotFound
		}
		return Booking{}, err
	}

	val, err := s.client.Get(ctx, seatKey).Result()
	if err != nil {
		if err == redis.Nil {
			return Booking{}, ErrSeatNotFound
		}
		return Booking{}, err
	}

	var b Booking
	if err := json.Unmarshal([]byte(val), &b); err != nil {
		return Booking{}, err
	}

	if err := s.client.Del(ctx, seatKey).Err(); err != nil {
		return Booking{}, err
	}
	if err := s.client.Del(ctx, sessionKey(sessionID)).Err(); err != nil {
		return Booking{}, err
	}
	return b, nil
}

func (s *RedisStore) ListBookings(movieID string) ([]Booking, error) {
	pattern := fmt.Sprintf("seat:%s:*", movieID)
	ctx := context.Background()
	keys, err := s.client.Keys(ctx, pattern).Result()
	if err != nil {
		return []Booking{}, err
	}

	bookings := make([]Booking, 0, len(keys))
	for _, key := range keys {
		val, err := s.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var b Booking
		if err := json.Unmarshal([]byte(val), &b); err != nil {
			continue
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}
