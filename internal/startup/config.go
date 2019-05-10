package startup

import "time"

type Config struct {
	Server Server `json:"server"`
}

type Server struct {
	Addr            string        `json:"addr"`
	ShutdownTimeout time.Duration `json:"shutdown-timeout"`
}

func ReadConfiguration() *Config {
	return &Config{
		Server: Server{
			Addr:            "0.0.0.0:8080",
			ShutdownTimeout: 5 * time.Second,
		},
	}
}
