package proxy

import (
	"errors"
	"math/rand"
	"sync"
)

var ErrNoHealthyEndpoints = errors.New("no healthy endpoints available")

// RouteEndpoint contains routing metadata for an endpoint.
type RouteEndpoint struct {
	Addr  string
	AZ    string
	Alive bool
}

// Router handles AZ-aware connection routing.
type Router struct {
	mu        sync.RWMutex
	endpoints map[string]*RouteEndpoint
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{
		endpoints: make(map[string]*RouteEndpoint),
	}
}

// AddEndpoint registers a downstream endpoint with its AZ.
func (r *Router) AddEndpoint(addr, az string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.endpoints[addr]; !exists {
		r.endpoints[addr] = &RouteEndpoint{
			Addr:  addr,
			AZ:    az,
			Alive: false, // Default to false until health checked
		}
	}
}

// UpdateHealth updates the alive status of an endpoint.
func (r *Router) UpdateHealth(addr string, alive bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ep, exists := r.endpoints[addr]; exists {
		ep.Alive = alive
	}
}

// SelectBest chooses the best available endpoint for the given client AZ.
// It prioritizes healthy endpoints in the same AZ, then falls back to any healthy endpoint.
func (r *Router) SelectBest(clientAZ string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var sameAZ []string
	var crossAZ []string

	for _, ep := range r.endpoints {
		if !ep.Alive {
			continue
		}
		if ep.AZ == clientAZ && clientAZ != "" {
			sameAZ = append(sameAZ, ep.Addr)
		} else {
			crossAZ = append(crossAZ, ep.Addr)
		}
	}

	if len(sameAZ) > 0 {
		// Randomly select one from the same AZ
		return sameAZ[rand.Intn(len(sameAZ))], nil
	}

	if len(crossAZ) > 0 {
		// Randomly select one from other AZs
		return crossAZ[rand.Intn(len(crossAZ))], nil
	}

	return "", ErrNoHealthyEndpoints
}
