# Technology Stack: SRE Sandbox with ERP

## 1. Programming Languages & Runtime Environments
- **Go (Golang) 1.22+**: Used for the high-performance transaction API (`order-service`).
- **Java 21+**: Used for the background event consumer and inventory manager (`inventory-service`).

## 2. Frameworks & Libraries
- **Go Order Service**:
  - `net/http` (Standard Library) for REST API.
  - `github.com/twmb/franz-go` for optimized, high-throughput Kafka producer communication.
  - `github.com/prometheus/client_golang/prometheus/promhttp` for native Prometheus metric exposing.
- **Java (Quarkus) Inventory Service**:
  - `Quarkus RESTEasy` for health checking and light administration.
  - `Native Apache Kafka Client` for Kafka consumer messaging.
  - `Micrometer Metrics` for Java runtime and application metrics.

## 3. Infrastructure & Orchestration
- **Local Dev Sandbox**: `Docker Compose` or `Kind (Kubernetes in Docker)` for zero-cost rapid local prototyping.
- **Cloud Cluster**: DigitalOcean Kubernetes (DOKS) with 3x worker nodes (2 vCPU, 4GB RAM) for production-grade SRE testing.
- **Load Balancing**: DigitalOcean Managed Load Balancer for public API traffic distribution.
- **Configuration Management**: `Helm` for managing Kubernetes manifests, applying soft/hard anti-affinity rules, and defining workload resource constraints.

## 4. Message Broker (Event-Driven Fabric)
- **Apache Kafka**: Single or multi-node broker deployment utilizing KRaft mode (no Zookeeper) within Kubernetes (using `apache/kafka` or standard Helm charts).

## 5. Chaos & Observability
- **Chaos Injections**: `Chaos Mesh` Custom Resource Definitions (CRDs) targeting cluster workloads (pod kills, network partitions).
- **Monitoring**: `kube-prometheus-stack` Helm chart (Prometheus Operator & Grafana instance).
- **Visualization**: Grafana Dashboards visualizing throughput, CPU/Memory limits, and pod eviction counts.
