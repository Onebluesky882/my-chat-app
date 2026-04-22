package cache

import (
	"context"
	"testing"

	"github.com/Onebluesky882/my-chat-app/internal/config"
)

func TestRedisSetGet(t *testing.T) {
	cfg := config.LoadConfig()

	rdb, err := NewRedis(cfg)
	if err != nil {
		t.Fatalf("redis connect error: %v", err)
	}

	ctx := context.Background()

	err = rdb.Set(ctx, "test:key", "123", 0).Err()
	if err != nil {
		t.Fatalf("set error: %v", err)
	}

	val, err := rdb.Get(ctx, "test:key").Result()
	if err != nil {
		t.Fatalf("get error: %v", err)
	}

	if val != "123" {
		t.Fatalf("expected 123, got %s", val)
	}
}
