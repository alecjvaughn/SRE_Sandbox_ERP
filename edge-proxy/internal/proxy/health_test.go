package proxy

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestHealthChecker_Probe(t *testing.T) {
	// 1. Start a mock server to act as a healthy downstream
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock downstream: %v", err)
	}
	defer l.Close()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close() // Just accept and close to prove it's alive
		}
	}()

	// 2. Test probing a healthy endpoint
	hc := NewHealthChecker(1 * time.Second)
	alive := hc.Probe(l.Addr().String())
	if !alive {
		t.Error("Expected healthy downstream to return true")
	}

	// 3. Test probing an invalid/dead endpoint
	// Use an unassigned port or unreachable address
	dead := hc.Probe("127.0.0.1:1") // Port 1 is typically blocked or unused
	if dead {
		t.Error("Expected dead downstream to return false")
	}
}

func TestHealthChecker_WorkerPool(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock downstream: %v", err)
	}
	addr := l.Addr().String()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	hc := NewHealthChecker(1 * time.Millisecond)
	hc.interval = 10 * time.Millisecond // Fast interval for testing
	hc.AddEndpoint(addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initially should be false before first probe
	if hc.IsAlive(addr) {
		t.Error("Expected endpoint to initially be false")
	}

	go hc.Start(ctx)

	// Wait for a few probe intervals
	time.Sleep(50 * time.Millisecond)

	if !hc.IsAlive(addr) {
		t.Error("Expected endpoint to be alive after probing")
	}

	// Close listener, simulating failure
	l.Close()
	time.Sleep(50 * time.Millisecond)

	if hc.IsAlive(addr) {
		t.Error("Expected endpoint to be dead after listener closed")
	}
}
