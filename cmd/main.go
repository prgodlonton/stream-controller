package main

import (
	"context"
	"github.com/pgodlonton/stream-controller/internal/startup"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// main entry point
func main() {
	config := startup.ReadConfiguration()
	resolver := startup.NewResolver(config)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	logger := resolver.ResolveLogger()
	logger.Info("starting...")

	// start server
	server := resolver.ResolveServer()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorw("unexpected http server listen error", "error", err)
		}
	}()

	// listen for interrupt/kill signal
	s := <-signals
	logger.Infow("caught signal: stopping...", "signal", s)

	// shutdown server
	waitTime := time.Duration(config.Server.ShutdownTimeout) * time.Second
	ctx, cfn := context.WithTimeout(context.Background(), waitTime)
	defer cfn()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorw("unclean http server shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("stopped gracefully")
}
