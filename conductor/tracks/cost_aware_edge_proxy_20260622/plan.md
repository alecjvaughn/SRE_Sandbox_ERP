# Implementation Plan: Custom Go Edge Proxy & Egress Cost Analysis

## Phase 1: Go Edge Proxy Core Network Engine
- [x] Task: Set up the proxy project scaffold in Go [a3d9b69]
    - [x] Initialize `go mod` in a new `edge-proxy` directory.
    - [x] Define the project structure (`cmd/`, `internal/proxy/`, `internal/metrics/`).
- [x] Task: Implement core TCP listening and connection handling (TDD) [e1ab397]
    - [x] Write tests for TCP listener and graceful shutdown.
    - [x] Implement TCP listener binding to a configured port.
- [ ] Task: Implement Layer 4 Pass-through and Routing
    - [ ] Write unit tests for connection copying (bidirectional byte streaming).
    - [ ] Implement the proxy tunneling logic to forward bytes from client to downstream broker.

## Phase 2: Dynamic Routing and Failover Logic
- [ ] Task: Implement downstream health checking
    - [ ] Write tests for active TCP probing of downstream endpoints.
    - [ ] Implement a worker pool to maintain the status of multiple Kafka broker endpoints.
- [ ] Task: Implement AZ-aware routing and failover
    - [ ] Write tests for routing logic (prefer healthy, least-latency, or specific AZ).
    - [ ] Update the proxy to select the best downstream connection upon client connect.
    - [ ] Implement automatic connection retry/failover if the downstream fails during the handshake.

## Phase 3: Cost-Centric Observability & Metrics
- [ ] Task: Instrument proxy with Prometheus metrics
    - [ ] Add `prometheus/client_golang` dependency.
    - [ ] Define custom counters: `proxy_bytes_in_total`, `proxy_bytes_out_total`, `proxy_active_connections`.
    - [ ] Expose an HTTP `/metrics` endpoint on a separate admin port.
- [ ] Task: Integrate bytes tracking into the TCP tunnel
    - [ ] Implement a custom `io.Reader` / `io.Writer` wrapper to count bytes passing through the proxy.
    - [ ] Update the Prometheus counters in real-time as data streams.

## Phase 4: Infrastructure & Chaos Deployment
- [ ] Task: Containerize and deploy the Edge Proxy
    - [ ] Write a `Dockerfile` for the `edge-proxy`.
    - [ ] Update `Makefile` to build and deploy the proxy image to AKS.
    - [ ] Create Kubernetes manifests (`Deployment`, `Service`, `ServiceMonitor`) and add them to ArgoCD.
- [ ] Task: Define and execute NetworkChaos experiments
    - [ ] Create `NetworkChaos` YAML manifests targeting the proxy to inject 500ms latency and 20% packet loss.
    - [ ] Deploy the chaos manifests to the cluster and validate proxy failover resilience.
- [ ] Task: Configure Grafana Egress Cost Dashboard
    - [ ] Create a Grafana dashboard JSON model plotting the proxy's network I/O.
    - [ ] Add a panel calculating simulated costs (e.g., $0.01 per GB for cross-AZ traffic).
- [ ] Task: Conductor - User Manual Verification 'Infrastructure & Chaos Deployment' (Protocol in workflow.md)
