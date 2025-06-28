package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Exists(ctx context.Context, key string) bool
}

type RedisClient struct {
	client *redis.Client
}

var RedisNotFound error = fmt.Errorf("Redis key not found")

func InitRedis(host, port, password string) *RedisClient {
	log.Printf("Connecting to Redis at %s:%s with password=%s", host, port, password)

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis connection err: ", err)
	}

	log.Println("Connected to redis")
	return &RedisClient{
		client: rdb,
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", RedisNotFound
	} else if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RedisClient) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Exists(ctx context.Context, key string) bool {
	res := r.client.Exists(ctx, key).Val()
	return res == 1
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
