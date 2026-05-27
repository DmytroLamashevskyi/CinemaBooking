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
	log.Println("Connected to Redis at ", client.Options().Addr)
	return client
}
