package storage

import "time"

type Cacher interface {
	Put(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
}
