package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Manager struct {
	client *redis.Client
}

func New(client *redis.Client) *Manager {
	return &Manager{client: client}
}

func (m *Manager) Get(ctx context.Context, key string) (string, error) {
	return m.client.Get(ctx, key).Result()
}

func (m *Manager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.client.Set(ctx, key, value, expiration).Err()
}
