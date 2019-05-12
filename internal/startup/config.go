package startup

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"os"
)

const (
	// ConsulAddr environment variable holding the consul address
	ConsulAddr = "CONSUL_ADDRESS"

	// ConsulKey environment variable holding the key where the config is stored
	ConsulKey = "CONSUL_KEY"
)

// Config holds all configuration
type Config struct {
	Redis  Redis  `json:"redis"`
	Server Server `json:"server"`
}

// Redis holds redis server configuration
type Redis struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// Server holds server-specific configuration
type Server struct {
	Address         string `json:"address"`
	ShutdownTimeout int    `json:"shutdown-timeout"`
}

func ReadConfiguration() *Config {
	consulAddress := getEnvValue(ConsulAddr, "http://localhost:8500")
	consulKey := getEnvValue(ConsulKey, "services/stream-control")

	kv := getConsulKV(consulAddress)
	return getConsulConfig(kv, consulKey)
}

func getEnvValue(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getConsulKV(address string) *api.KV {
	client, err := api.NewClient(
		&api.Config{
			Address: address,
		},
	)
	if err != nil {
		panic(err)
	}
	return client.KV()
}

func getConsulConfig(store *api.KV, consulKey string) *Config {
	pair, _, err := store.Get(consulKey, &api.QueryOptions{})
	if err != nil {
		panic(err)
	}
	var config Config
	if err := json.Unmarshal(pair.Value, &config); err != nil {
		panic(err)
	}
	return &config
}
