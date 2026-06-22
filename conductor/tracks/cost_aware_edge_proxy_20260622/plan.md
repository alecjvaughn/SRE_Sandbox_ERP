# Implementation Plan: Custom Go Edge Proxy & Egress Cost Analysis

## Phase 1: Go Edge Proxy Core Network Engine [checkpoint: 460e955]
- [x] Task: Set up the proxy project scaffold in Go [a3d9b69]
    - [x] Initialize `go mod` in a new `edge-proxy` directory.
    - [x] Define the project structure (`cmd/`, `internal/proxy/`, `internal/metrics/`).
- [x] Task: Implement core TCP listening and connection handling (TDD) [e1ab397]
    - [x] Write tests for TCP listener and graceful shutdown.
    - [x] Implement TCP listener binding to a configured port.
- [x] Task: Implement Layer 4 Pass-through and Routing [3948233]
    - [x] Write unit tests for connection copying (bidirectional byte streaming).
    - [x] Implement the proxy tunneling logic to forward bytes from client to downstream broker.

## Phase 2: Dynamic Routing and Failover Logic [checkpoint: a5a85d4]
- [x] Task: Implement downstream health checking [bfdc716]
    - [x] Write tests for active TCP probing of downstream endpoints.
    - [x] Implement a worker pool to maintain the status of multiple Kafka broker endpoints.
- [x] Task: Implement AZ-aware routing and failover [0d1dfd7]
    - [x] Write tests for routing logic (prefer healthy, least-latency, or specific AZ).
    - [x] Update the proxy to select the best downstream connection upon client connect.
    - [x] Implement automatic connection retry/failover if the downstream fails during the handshake.

## Phase 3: Cost-Centric Observability & Metrics [checkpoint: cab77dd]
- [x] Task: Instrument proxy with Prometheus metrics [1d1844d]
    - [x] Add `prometheus/client_golang` dependency.
    - [x] Define custom counters: `proxy_bytes_in_total`, `proxy_bytes_out_total`, `proxy_active_connections`.
    - [x] Expose an HTTP `/metrics` endpoint on a separate admin port.
- [x] Task: Integrate bytes tracking into the TCP tunnel [1d1844d]
    - [x] Implement a custom `io.Reader` / `io.Writer` wrapper to count bytes passing through the proxy.
    - [x] Update the Prometheus counters in real-time as data streams.

## Phase 4: Infrastructure & Chaos Deployment [checkpoint: 6110bbe]
- [x] Task: Containerize and deploy the Edge Proxy [34aef02]
    - [x] Write a `Dockerfile` for the `edge-proxy`.
    - [x] Update `Makefile` to build and deploy the proxy image to AKS.
    - [x] Create Kubernetes manifests (`Deployment`, `Service`, `ServiceMonitor`) and add them to ArgoCD.
- [x] Task: Define and execute NetworkChaos experiments [fefa5ee]
    - [x] Create `NetworkChaos` YAML manifests targeting the proxy to inject 500ms latency and 20% packet loss.
    - [x] Deploy the chaos manifests to the cluster and validate proxy failover resilience.
- [x] Task: Configure Grafana Egress Cost Dashboard [e149c12]
    - [x] Create a Grafana dashboard JSON model plotting the proxy's network I/O.
    - [x] Add a panel calculating simulated costs (e.g., $0.01 per GB for cross-AZ traffic).
    - [x] Add dashboard config map to ArgoCD observability manifests.
- [ ] Task: Conductor - User Manual Verification 'Infrastructure & Chaos Deployment' (Protocol in workflow.md)
