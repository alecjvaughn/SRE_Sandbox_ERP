package main

import (
	"context"
	"flag"
	"log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecjvaughn/sre-sandbox/edge-proxy/internal/proxy"
)

func main() {
	bindAddr := flag.String("bind", ":8080", "Address to bind the proxy to")
	adminAddr := flag.String("admin", ":9090", "Address for admin HTTP server (metrics)")
	downstreamAddr := flag.String("downstream", "google.com:80", "Temporary single downstream address to forward to")
	flag.Parse()

	// Start admin server for metrics
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Printf("Starting admin server on %s", *adminAddr)
		if err := http.ListenAndServe(*adminAddr, mux); err != nil {
			log.Fatalf("Admin server failed: %v", err)
		}
	}()

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
