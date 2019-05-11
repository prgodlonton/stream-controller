package startup

import (
	"github.com/pgodlonton/stream-controller/internal/handlers"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Resolver struct {
	config *Config

	// singletons
	logger *zap.SugaredLogger
	server *http.Server
}

func NewResolver(config *Config) *Resolver {
	resolver := &Resolver{
		config: config,
	}
	return resolver
}

func (r *Resolver) ResolveLogger() *zap.SugaredLogger {
	if r.logger == nil {
		logger, _ := zap.NewDevelopment()
		r.logger = logger.Sugar()
	}
	return r.logger
}

func (r *Resolver) ResolveHTTPHandler() http.Handler {
	return handlers.NewHandler(
		r.ResolveLogger(),
	)
}

func (r *Resolver) ResolveHTTPServer() *http.Server {
	if r.server == nil {
		r.server = &http.Server{
			Addr:         r.config.Server.Address,
			Handler:      r.ResolveHTTPHandler(),
			IdleTimeout:  time.Second * 60,
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
		}
	}
	return r.server
}
