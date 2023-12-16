package redis

import (
	"encoding/json"

	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func New(address string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	_, err := client.Ping(nil).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	return r.client.Get(nil, key).Result()
}

func (r *RedisCache) Put(key string, value interface{}, ttl time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ret := r.client.Set(nil, key, json, time.Duration(ttl)*time.Second).Err()
	return ret
}
