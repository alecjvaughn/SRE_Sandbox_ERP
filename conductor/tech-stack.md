# Technology Stack: Cost-Aware Traffic Plane & Egress Optimizer

## 1. Programming Languages & Runtime Environments
- **Go (Golang) 1.22+**: Used for the Custom Edge Traffic Proxy and the high-performance transaction API (`order-service`).
- **Java 21+**: Used for the background event consumer and inventory manager (`inventory-service`).

## 2. Frameworks & Libraries
- **Custom Go Edge Proxy**:
  - Raw `net` package for Layer 4 TCP proxying and routing.
  - Custom extensions for active health-checking and least-cost AZ routing.
- **Go Order Service**:
  - `net/http` (Standard Library) for REST API.
  - `github.com/twmb/franz-go` for optimized, high-throughput Kafka producer communication.
  - `github.com/prometheus/client_golang/prometheus/promhttp` for native Prometheus metric exposing.
- **Java (Quarkus) Inventory Service**:
  - `Quarkus RESTEasy` for health checking and light administration.
  - `Native Apache Kafka Client` for Kafka consumer messaging.
  - `Micrometer Metrics` for Java runtime and application metrics.

## 3. Infrastructure & Orchestration (Azure Pay-As-You-Go)
- **Local Dev Sandbox**: `Docker Compose` or `Kind (Kubernetes in Docker)` for zero-cost rapid local prototyping.
- **Cloud Cluster**: Azure Kubernetes Service (AKS) with 3x worker nodes deployed strictly within the Azure free tier / $200 credit limit constraints.
- **Advanced Networking**: Azure Virtual Networks (VNets), User Defined Routes (UDRs), Network Security Groups (NSGs), and Azure Private Endpoints to simulate Confluent Cloud's isolated tenant data streams.
- **Load Balancing**: Azure Standard Load Balancer for ingress traffic distribution.
- **Configuration Management**: `Helm` for managing Kubernetes manifests, applying soft/hard anti-affinity rules, and defining workload resource constraints.

## 4. Message Broker (Event-Driven Fabric)
- **Apache Kafka**: Confluent for Kubernetes (CFK) Operator deploying a KRaft-enabled Kafka architecture.

## 5. Chaos & Cost-Centric Observability
- **Chaos Injections**: `Chaos Mesh` Custom Resource Definitions (CRDs), specifically `NetworkChaos`, targeting packet loss, bandwidth throttling, and 500ms latency delays to simulate cross-AZ link drops.
- **Monitoring**: `kube-prometheus-stack` Helm chart (Prometheus Operator & Grafana instance).
- **Cost & Traffic Visualization**: Custom Grafana Dashboards tracking Network I/O (Bytes In/Out), TCP Retransmission rates, and real-time calculated Egress Cost ($) based on simulated cross-AZ routing.
