package main

import (
	"context"
	"fmt"
	"github.com/pgodlonton/stream-controller/internal/startup"
	"net/http"
	"os"
	"os/signal"
)

// main entry point
func main() {
	config := startup.ReadConfiguration()
	resolver := startup.NewResolver(config)

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	fmt.Println("starting...")

	// start server
	server := resolver.ResolveHTTPServer()
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("http server stopped due to %v\n", err.Error())
		}
	}()

	// listen for interrupt/kill signal
	s := <-signals
	fmt.Printf("caught %v: stopping...\n", s.String())

	// shutdown server
	ctx, cfn := context.WithTimeout(context.Background(), config.Server.ShutdownTimeout)
	defer cfn()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("unclean http server shutdown due to %v\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("stopped successfully")
}
