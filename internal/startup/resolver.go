package startup

import (
	"github.com/pgodlonton/stream-controller/internal/handlers"
	"net/http"
	"time"
)

type Resolver struct {
	config *Config

	// singletons
	server *http.Server
}

func NewResolver(config *Config) *Resolver {
	resolver := &Resolver{
		config: config,
	}
	return resolver
}

func (r *Resolver) ResolveHTTPHandler() http.Handler {
	return &handlers.Handler{}
}

func (r *Resolver) ResolveHTTPServer() *http.Server {
	if r.server == nil {
		r.server = &http.Server{
			Addr:         r.config.Server.Addr,
			Handler:      r.ResolveHTTPHandler(),
			IdleTimeout:  time.Second * 60,
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
		}
	}
	return r.server
}
