package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

func Connect() *redis.Client {
	addr := fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})
	return rdb
}

func Set(ctx context.Context, key string, value string, expiration time.Duration) {
	rdb := Connect()
	defer rdb.Close()
	_ = rdb.Set(ctx, key, value, expiration)
}

func Get(ctx context.Context, key string) string {
	rdb := Connect()
	defer rdb.Close()
	res := rdb.Get(ctx, key)
	return res.Val()
}

func Del(ctx context.Context, state_key string, idUser string) {
	rdb := Connect()
	defer rdb.Close()
	_ = rdb.HDel(ctx, state_key, idUser)
}
