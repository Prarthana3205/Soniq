package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() {
	// Default local Redis
	addr := "localhost:6379"
	password := ""

	// Use Railway REDIS_URL if available (format: redis://default:password@host:port)
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			fmt.Println("Failed to parse REDIS_URL:", err)
		} else {
			rdb = redis.NewClient(opt)
			fmt.Println("Redis connected via REDIS_URL")
			return
		}
	}

	// Fallback to REDIS_HOST env var
	if os.Getenv("REDIS_HOST") != "" {
		addr = os.Getenv("REDIS_HOST")
		password = os.Getenv("REDIS_PASSWORD")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis connection error:", err)
	} else {
		fmt.Println("Redis connected to", addr)
	}
}

func PublishMessage(msg string) {
	err := rdb.Publish(ctx, "soniq-messages", msg).Err()
	if err != nil {
		fmt.Println("Redis publish error:", err)
	}
}

func Subscribe(callback func(string)) {
	sub := rdb.Subscribe(ctx, "soniq-messages")
	ch := sub.Channel()

	go func() {
		for msg := range ch {
			callback(msg.Payload)
		}
	}()
}
