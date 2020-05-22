package redis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func New(address string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	return r.client.Get(key).Result()
}

func (r *RedisCache) Put(key string, value interface{}, ttl time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(key, json, time.Duration(ttl)*time.Second).Err()
}
