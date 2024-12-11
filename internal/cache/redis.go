package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"stone-api/internal/config"
)

func Init() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Get().Cache.Host, config.Get().Cache.Port),
		Password: "", // no password
		DB:       0,
	})
}
