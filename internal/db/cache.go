package db

import (
	"context"
	"time"

	redisdb "github.com/saurify/saurify-backend/internal/redis"
)

const CacheExpiration = 24 * time.Hour

func SaveURLToCache(shortCode, originalURL string) error {
	ctx := context.Background()
	return redisdb.RDB.Set(ctx, shortCode, originalURL, CacheExpiration).Err()
}

func GetURLToCache(shortCode string) (string, error) {
	ctx := context.Background()
	return redisdb.RDB.Get(ctx, shortCode).Result()
}
