package internal

// Store records the streams being watched by users
type Store interface {
	Add(userID, streamID string) error
	Remove(userID, streamID string) error
}
