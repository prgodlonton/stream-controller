package internal

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var exceededStreamsQuota = errors.New("user has exceeded streaming quota")

const (
	// *atomic* lua script to add strings to the given set if it has less than 3 elements
	condSetAdd = `if redis.call("SCARD",KEYS[1]) < 3 then return redis.call("SADD",KEYS[1],ARGV[1]) else return -1 end`
)

// Store records the streams being watched by users
type Store interface {
	AddStream(userID, streamID string) error
	GetStreams(userID string) ([]string, error)
	RemoveStream(userID, streamID string) error
}

// RedisStore a Redis-backed store
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new Redis-backed store
func NewRedisStore(client *redis.Client) Store {
	return &RedisStore{
		client: client,
	}
}

// Adds records a user as watching a stream
func (rs *RedisStore) AddStream(userID, streamID string) error {
	cmd := rs.client.Eval(condSetAdd, []string{userID}, streamID)
	val, err := cmd.Result()
	if err != nil {
		return errors.Wrap(err, "failed to add element to list")
	}

	added, ok := val.(int64)
	if !ok {
		return errors.New("cannot convert redis eval return value to int64")
	}
	if added == -1 {
		return exceededStreamsQuota
	}
	return nil
}

// Get returns all stream being watched by a single user
func (rs *RedisStore) GetStreams(userID string) ([]string, error) {
	cmd := rs.client.SMembers(userID)
	elements, err := cmd.Result()
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to get list elements")
	}
	return elements, nil
}

// Remove removes the record of a user watching a stream
func (rs *RedisStore) RemoveStream(userID, streamID string) error {
	cmd := rs.client.SRem(userID, streamID)
	if _, err := cmd.Result(); err != nil {
		return errors.Wrap(err, "failed to remove element from list")
	}
	return nil
}
