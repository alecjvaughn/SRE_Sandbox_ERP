package proxy

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestServer_StartAndStop(t *testing.T) {
	// 1. Create a new proxy server
	srv := NewServer("127.0.0.1:0") // Port 0 asks OS for random free port

	// 2. Start the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// 3. Verify it's listening by trying to connect
	addr := srv.ListenAddr()
	if addr == "" {
		t.Fatal("Expected server to have a listen address, got empty string")
	}

	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		t.Fatalf("Failed to connect to proxy server: %v", err)
	}
	conn.Close()

	// 4. Test graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Failed to shutdown gracefully: %v", err)
	}

	// 5. Verify start returned nil (or expected closed error)
	select {
	case err := <-errCh:
		if err != nil && err != net.ErrClosed {
			t.Errorf("Expected nil or net.ErrClosed from Start, got: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Error("Start() did not return after Shutdown()")
	}
}
