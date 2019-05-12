package startup

import (
	"github.com/pgodlonton/stream-controller/internal"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Resolver struct {
	config *Config

	// singletons
	logger *zap.SugaredLogger
	server *http.Server
	store  internal.Store
}

func NewResolver(config *Config) *Resolver {
	resolver := &Resolver{
		config: config,
	}
	resolver.resolveEager()
	return resolver
}

func (r *Resolver) resolveEager() {
	r.ResolveLogger()
	r.ResolveServer()
	r.ResolveStore()
}

func (r *Resolver) ResolveLogger() *zap.SugaredLogger {
	if r.logger == nil {
		logger, _ := zap.NewDevelopment()
		r.logger = logger.Sugar()
	}
	return r.logger
}

func (r *Resolver) ResolveRouter() http.Handler {
	return internal.NewRouter(
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
	if r.store == nil {
		r.store = &internal.RedisStore{}
	}
	return r.store
}
