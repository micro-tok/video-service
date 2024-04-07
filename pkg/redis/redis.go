package redis

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/micro-tok/video-service/pkg/config"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(cfg *config.Config) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		panic(err)
	}

	return &RedisService{
		client: client,
	}
}

func (s RedisService) Set(key string, value interface{}) error {
	return s.client.Set(s.client.Context(), key, value, time.Hour).Err()
}

func (s RedisService) Get(key string) (string, error) {
	return s.client.Get(s.client.Context(), key).Result()
}
