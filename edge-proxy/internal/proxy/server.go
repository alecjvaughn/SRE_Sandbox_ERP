package proxy

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// Server represents the TCP edge proxy server.
type Server struct {
	bindAddr string
	listener net.Listener
	wg       sync.WaitGroup
	quit     chan struct{}

	router        *Router
	healthChecker *HealthChecker
	localAZ       string // The AZ this proxy instance is running in
}

// NewServer creates a new proxy Server bound to the given address.
func NewServer(bindAddr string) *Server {
	return &Server{
		bindAddr:      bindAddr,
		quit:          make(chan struct{}),
		router:        NewRouter(),
		healthChecker: NewHealthChecker(1 * time.Second),
		localAZ:       "default",
	}
}

// SetLocalAZ configures the AZ for this proxy instance to prefer local routes.
func (s *Server) SetLocalAZ(az string) {
	s.localAZ = az
}

// AddDownstream registers a downstream address with an AZ and starts health checking.
func (s *Server) AddDownstream(addr, az string) {
	s.router.AddEndpoint(addr, az)
	s.healthChecker.AddEndpoint(addr)
}

// SetDownstream sets the downstream address. Temporary for Phase 1.
func (s *Server) SetDownstream(addr string) {
	s.AddDownstream(addr, s.localAZ)
}

// Start begins listening for TCP connections.
// It blocks until Shutdown is called or a listener error occurs.
func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Link HealthChecker to Router
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.healthChecker.mu.RLock()
				for epAddr, status := range s.healthChecker.endpoints {
					s.router.UpdateHealth(epAddr, status.Alive)
				}
				s.healthChecker.mu.RUnlock()
			}
		}
	}()

	go s.healthChecker.Start(ctx)

	l, err := net.Listen("tcp", s.bindAddr)
	if err != nil {
		return err
	}
	s.listener = l

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				// Shutdown was called
				return net.ErrClosed
			default:
				// Return other listener errors
				return err
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(clientConn net.Conn) {
	defer s.wg.Done()
	defer clientConn.Close()

	var downstreamConn net.Conn
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		var addr string
		addr, err = s.router.SelectBest(s.localAZ)
		if err != nil {
			log.Printf("No healthy downstream found: %v\n", err)
			return
		}

		downstreamConn, err = net.Dial("tcp", addr)
		if err == nil {
			break
		}
		
		log.Printf("Failed to connect to %s, marking as dead and retrying...\n", addr)
		s.router.UpdateHealth(addr, false) // Mark dead immediately to pick another
	}

	if downstreamConn == nil {
		log.Println("Exhausted retries connecting to downstream")
		return
	}
	defer downstreamConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(downstreamConn, clientConn)
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientConn, downstreamConn)
	}()

	wg.Wait()
}

// ListenAddr returns the actual network address the server is listening on.
// This is useful if the server was started on port 0.
func (s *Server) ListenAddr() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return ""
}

// Shutdown gracefully stops the server, waiting for active connections to finish.
func (s *Server) Shutdown(ctx context.Context) error {
	close(s.quit)

	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return err
		}
	}

	// Wait for active connections to finish or context timeout
	c := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		return nil // Success
	case <-ctx.Done():
		return errors.New("shutdown timeout")
	}
}
