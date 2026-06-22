package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecjvaughn/sre-sandbox/edge-proxy/internal/proxy"
)

func main() {
	bindAddr := flag.String("bind", ":8080", "Address to bind the proxy to")
	downstreamAddr := flag.String("downstream", "google.com:80", "Temporary single downstream address to forward to")
	flag.Parse()

	srv := proxy.NewServer(*bindAddr)
	srv.SetDownstream(*downstreamAddr)
	
	errCh := make(chan error, 1)
	go func() {
		log.Printf("Starting proxy server on %s forwarding to %s", *bindAddr, *downstreamAddr)
		errCh <- srv.Start()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Fatalf("Server failed: %v", err)
	case sig := <-sigCh:
		log.Printf("Received signal %v, shutting down...", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v", err)
	}
	log.Println("Server stopped gracefully.")
}
