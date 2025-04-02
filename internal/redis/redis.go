package redisdb

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(addr string) {
	RDB = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("Redis connection failed")
	}

	log.Println("Connected to Redis successfully!")
}
