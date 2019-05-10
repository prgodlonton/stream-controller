package main

import (
	"context"
	"fmt"
	"github.com/pgodlonton/stream-controller/internal/startup"
	"os"
	"os/signal"
	"time"
)

// main entry point
func main() {
	_, cfn := context.WithCancel(context.Background())
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)

	fmt.Println("starting...")

	// start server

	// listen for interrupt/kill signal
	s := <-signals
	fmt.Printf("caught %v: stopping...\n", s.String())

	// shutdown server
	cfn()
	go func() {
		time.AfterFunc(startup.ServerShutdownTimeout, func() {
			fmt.Println("took too long, forcibly stopping")
			os.Exit(1)
		})
	}()

	fmt.Println("stopped successfully")
}
