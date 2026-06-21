# SRE Sandbox ERP

This repository contains the core transactional and event-driven foundation of the SRE Sandbox. It simulates an ERP architecture consisting of a Go-based order service and a Java-based inventory service communicating via a Kafka event bus.

## 🏗 Architecture

- **Go Order Service**: Active-active REST gateway handling incoming orders.
- **Java (Quarkus) Inventory Service**: Background event processing consumer listening to Kafka.
- **Kafka**: Confluent Kafka broker managing the event bus.
- **Observability**: Prometheus scraping and Grafana dashboards.
- **Resilience**: Chaos Mesh for injecting latency, failures, and evictions.

## 🚀 Current Status

- **Phase 1-4**: Core Microservices & Local Sandbox (Completed).
- **Phase 5**: Kubernetes Active-Active Infrastructure (Completed).
- **Phase 6**: Observability Stack (Completed).
- **Phase 7**: Chaos Engineering & Testing (Completed).

## 🛠 Prerequisites

- [Go](https://golang.org/doc/install) (1.22+)
- [Java 21+](https://jdk.java.net/) & Maven
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose

## 📂 Project Structure

```
.
├── PLAN_OVERVIEW.md    # High-level roadmap and architectural goals
├── conductor/          # Tracks, specification, and agent workflow plans
├── inventory-service/  # Java Quarkus Kafka consumer and metrics server
├── order-service/      # Go REST API for handling orders
└── README.md           # This evolving project documentation
```

## 💻 Getting Started

The local sandbox runs completely via Docker Compose, utilizing Apache Kafka in KRaft mode (no Zookeeper required).

1. **Start the local sandbox:**
   ```bash
   docker compose up --build -d
   ```

2. **Verify the services:**
   - Go Order Service API: `http://localhost:8080/order`
   - Inventory Service Health Check: `http://localhost:8081/healthz`
   - Order Service Metrics: `http://localhost:8080/metrics`
   - Inventory Service Metrics: `http://localhost:8000/metrics`

3. **Send a test order:**
   ```bash
   curl -X POST http://localhost:8080/order -d '{"order_id": "test-1", "item": "Widget", "qty": 5}'
   ```

## ⚠️ Known Issues

### Grafana Dashboards "No Data" in Minikube
When running the `kube-prometheus-stack` locally in Minikube, the default Grafana dashboard "Kubernetes / Compute Resources / Namespace (Workloads)" may show "No Data" for CPU and Memory metrics.
- **Root Cause**: The default Prometheus recording rules require `image!=""` and `container!=""` labels. However, Minikube's `containerd` cAdvisor drops these labels for root pod cgroups. The dashboard's hardcoded JSON filters cause the metrics to be filtered out entirely. Furthermore, Prometheus fails to scrape Minikube's `kubelet` by default due to self-signed TLS certificates.
- **Workaround 1 (TLS Bypass)**: Apply a patch to the `kubelet.serviceMonitor` to set `insecureSkipVerify: true` (e.g. via `helm/observability/values-minikube.yaml`).
- **Workaround 2 (Raw Queries)**: Instead of relying on the default dashboards, use the **Grafana Explore** tab to query the raw metrics directly. For example:
  - CPU Usage: `rate(container_cpu_usage_seconds_total{namespace="default"}[1m])`
  - Memory Usage: `container_memory_working_set_bytes{namespace="default"}`
