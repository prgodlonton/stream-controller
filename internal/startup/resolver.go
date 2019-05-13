package startup

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/prgodlonton/stream-controller/internal"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// Resolver resolves all dependencies and handles dependency injection
type Resolver struct {
	config *Config

	// singletons
	client *redis.Client
	logger *zap.SugaredLogger
	server *http.Server
}

// NewResolver returns a new resolver
func NewResolver(config *Config) *Resolver {
	resolver := &Resolver{
		config: config,
	}
	resolver.resolveEager()
	return resolver
}

func (r *Resolver) resolveEager() {
	r.ResolveLogger()
	r.ResolveRedisClient()
	r.ResolveServer()
}

func (r *Resolver) ResolveLogger() *zap.SugaredLogger {
	if r.logger == nil {
		logger, _ := zap.NewDevelopment()
		r.logger = logger.Sugar()
	}
	return r.logger
}

func (r *Resolver) ResolveRedisClient() *redis.Client {
	if r.client == nil {
		r.client = redis.NewClient(
			&redis.Options{
				Addr:     r.config.Redis.Address,
				Password: r.config.Redis.Password,
				DB:       r.config.Redis.DB,
			},
		)
		if _, err := r.client.Ping().Result(); err != nil {
			panic(errors.Wrap(err, "resolver: failed to ping redis server"))
		}
	}
	return r.client
}

func (r *Resolver) ResolveRouter() http.Handler {
	return internal.NewRouter(
		r.ResolveLogger(),
		r.ResolveStore(),
	)
}

func (r *Resolver) ResolveServer() *http.Server {
	if r.server == nil {
		r.server = &http.Server{
			Addr:         r.config.Server.Address,
			Handler:      r.ResolveRouter(),
			IdleTimeout:  time.Second * 60,
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
		}
	}
	return r.server
}

func (r *Resolver) ResolveStore() internal.Store {
	return internal.NewRedisStore(
		r.ResolveRedisClient(),
	)
}
