package config

import (
	"context"
	"fmt"

	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_CONN_STR")
	if redisURL == "" {
		return nil, fmt.Errorf("failed to get redis connstr")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	rdb := redis.NewClient(opt)

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	fmt.Println("✅ Redis connected")

	return rdb, nil
}
