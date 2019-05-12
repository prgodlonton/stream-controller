package internal

import "github.com/go-redis/redis"

// Store records the streams being watched by users
type Store interface {
	Add(userID, streamID string) error
	Remove(userID, streamID string) error
}

// RedisStore records the streams in a Redis server cluster
type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) Store {
	return &RedisStore{
		client: client,
	}
}

func (rs *RedisStore) Add(userID, streamID string) error {
	return nil
}

func (rs *RedisStore) Remove(userID, streamID string) error {
	return nil
}
