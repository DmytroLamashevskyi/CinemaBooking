package redis

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr string) *goredis.Client {
	client := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
		return nil
	}
	log.Printf("Connected to Redis at %s\n", client.Options().Addr)
	return client
}
