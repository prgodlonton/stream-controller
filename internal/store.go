package internal

// Store records the streams being watched by users
type Store interface {
	Add(userID, streamID string) error
	Remove(userID, streamID string) error
}

// RedisStore records the streams in a Redis server cluster
type RedisStore struct {
}

func (rs *RedisStore) Add(userID, streamID string) error {
	return nil
}

func (rs *RedisStore) Remove(userID, streamID string) error {
	return nil
}
