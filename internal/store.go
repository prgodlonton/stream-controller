package internal

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// Store records the streams being watched by users
type Store interface {
	Add(userID, streamID string) error
	Get(userID string) ([]string, error)
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
	cmd := rs.client.RPush(userID, streamID)
	if _, err := cmd.Result(); err != nil {
		return errors.Wrap(err, "failed to add element to list")
	}
	return nil
}

func (rs *RedisStore) Get(userID string) ([]string, error) {
	cmd := rs.client.LRange(userID, 0, -1)
	elements, err := cmd.Result()
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get list elements")
	}
	return elements, nil
}

func (rs *RedisStore) Remove(userID, streamID string) error {
	cmd := rs.client.LRem(userID, 1, streamID)
	if _, err := cmd.Result(); err != nil {
		return errors.Wrap(err, "failed to remove element from list")
	}
	return nil
}
