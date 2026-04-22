package cache

import (
	"context"
	"fmt"

	"github.com/Onebluesky882/my-chat-app/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connect error: %w", err)
	}
	fmt.Println("connected to redis")
	return rdb, nil
}
