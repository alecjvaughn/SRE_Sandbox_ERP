package proxy

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/alecjvaughn/sre-sandbox/edge-proxy/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/testutil"
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

func TestServer_ProxyConnection(t *testing.T) {
	// 1. Start a mock downstream server (like a Kafka broker)
	downstream, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start downstream server: %v", err)
	}
	defer downstream.Close()

	// Mock downstream behavior: echo back whatever it receives
	go func() {
		for {
			conn, err := downstream.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				buf := make([]byte, 1024)
				n, _ := c.Read(buf)
				c.Write(buf[:n])
			}(conn)
		}
	}()

	// 2. Start the proxy server pointing to the downstream
	srv := NewServer("127.0.0.1:0")
	srv.SetDownstream(downstream.Addr().String())

	go srv.Start()
	time.Sleep(500 * time.Millisecond)
	defer srv.Shutdown(context.Background())

	// 3. Connect to proxy and send data
	conn, err := net.Dial("tcp", srv.ListenAddr())
	if err != nil {
		t.Fatalf("Failed to connect to proxy: %v", err)
	}
	defer conn.Close()

	msg := []byte("hello kafka")
	if _, err := conn.Write(msg); err != nil {
		t.Fatalf("Failed to write to proxy: %v", err)
	}

	// 4. Verify proxy forwarded the data back (echoed from downstream)
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from proxy: %v", err)
	}

	if string(buf[:n]) != string(msg) {
		t.Errorf("Expected %q, got %q", msg, buf[:n])
	}

	// 5. Verify metrics incremented
	time.Sleep(50 * time.Millisecond) // Wait for go-routines to record metrics
	
	bytesIn := testutil.ToFloat64(metrics.BytesInTotal.WithLabelValues("default"))
	if bytesIn < float64(len(msg)) {
		t.Errorf("Expected BytesInTotal >= %d, got %v", len(msg), bytesIn)
	}
}
