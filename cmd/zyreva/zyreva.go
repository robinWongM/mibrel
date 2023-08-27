package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/zyreva/zyreva/internal/pkg/api"
	"github.com/zyreva/zyreva/internal/pkg/tasks"
)

func main() {
	go tasks.RunServer()

	rpcServer := api.NewServer()
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: rpcServer.Handler,
	}
	go server.ListenAndServe()

	// create a channel to subscribe ctrl+c/SIGINT event
	sigInterruptChannel := make(chan os.Signal, 1)
	signal.Notify(sigInterruptChannel, os.Interrupt)
	// block execution from continuing further until SIGINT comes
	<-sigInterruptChannel

	// create a context which will expire after 4 seconds of grace period
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	go server.Shutdown(ctx)

	// wait until ctx ends (which will happen after 4 seconds)
	<-ctx.Done()
}
