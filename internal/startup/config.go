package startup

import "time"

const (

	// ServerShutdownTimeout specifies the time to wait for the server to finish serving pending requests.
	ServerShutdownTimeout = 5 * time.Second
)


