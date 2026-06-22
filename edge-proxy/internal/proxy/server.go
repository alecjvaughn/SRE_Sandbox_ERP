package proxy

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

// Server represents the TCP edge proxy server.
type Server struct {
	bindAddr       string
	downstreamAddr string // Temporary single downstream
	listener       net.Listener
	wg             sync.WaitGroup
	quit           chan struct{}
}

// NewServer creates a new proxy Server bound to the given address.
func NewServer(bindAddr string) *Server {
	return &Server{
		bindAddr: bindAddr,
		quit:     make(chan struct{}),
	}
}

// Start begins listening for TCP connections.
// It blocks until Shutdown is called or a listener error occurs.
func (s *Server) Start() error {
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

	if s.downstreamAddr == "" {
		log.Println("No downstream configured")
		return
	}

	downstreamConn, err := net.Dial("tcp", s.downstreamAddr)
	if err != nil {
		log.Printf("Failed to connect to downstream %s: %v\n", s.downstreamAddr, err)
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
