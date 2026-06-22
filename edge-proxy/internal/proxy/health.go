package proxy

import (
	"context"
	"net"
	"sync"
	"time"
)

// EndpointStatus represents the health status of a downstream endpoint.
type EndpointStatus struct {
	Addr   string
	Alive  bool
	LastOk time.Time
}

// HealthChecker manages active probing of downstream endpoints.
type HealthChecker struct {
	timeout  time.Duration
	interval time.Duration
	
	mu        sync.RWMutex
	endpoints map[string]*EndpointStatus
}

// NewHealthChecker creates a HealthChecker.
func NewHealthChecker(timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		timeout:   timeout,
		interval:  5 * time.Second, // Default interval
		endpoints: make(map[string]*EndpointStatus),
	}
}

// Probe attempts to establish a TCP connection to the addr.
// Returns true if successful.
func (hc *HealthChecker) Probe(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, hc.timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// AddEndpoint adds an endpoint to be continuously monitored.
func (hc *HealthChecker) AddEndpoint(addr string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	if _, exists := hc.endpoints[addr]; !exists {
		hc.endpoints[addr] = &EndpointStatus{Addr: addr, Alive: false}
	}
}

// IsAlive returns the last known status of the endpoint.
func (hc *HealthChecker) IsAlive(addr string) bool {
	hc.mu.RLock()
	defer hc.mu.RUnlock()
	if status, exists := hc.endpoints[addr]; exists {
		return status.Alive
	}
	return false
}

// Start begins the background worker pool for active health checks.
func (hc *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	// Initial probe
	hc.probeAll()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hc.probeAll()
		}
	}
}

func (hc *HealthChecker) probeAll() {
	hc.mu.RLock()
	addrs := make([]string, 0, len(hc.endpoints))
	for addr := range hc.endpoints {
		addrs = append(addrs, addr)
	}
	hc.mu.RUnlock()

	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(a string) {
			defer wg.Done()
			alive := hc.Probe(a)
			
			hc.mu.Lock()
			if status, exists := hc.endpoints[a]; exists {
				status.Alive = alive
				if alive {
					status.LastOk = time.Now()
				}
			}
			hc.mu.Unlock()
		}(addr)
	}
	wg.Wait()
}
