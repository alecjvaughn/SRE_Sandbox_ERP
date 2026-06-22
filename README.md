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
- **Phase 8**: Infrastructure as Code & GitOps (In Progress).

## 🛠 Prerequisites

- [Go](https://golang.org/doc/install) (1.22+)
- [Java 21+](https://jdk.java.net/) & Maven
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Terraform](https://www.terraform.io/)
- [Kind](https://kind.sigs.k8s.io/)

## 📂 Project Structure

```
.
├── PLAN_OVERVIEW.md    # High-level roadmap and architectural goals
├── argocd/             # GitOps definitions (App of Apps)
├── conductor/          # Tracks, specification, and agent workflow plans
├── inventory-service/  # Java Quarkus Kafka consumer and metrics server
├── order-service/      # Go REST API for handling orders
├── terraform/          # Infrastructure as Code (Terraform configs)
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

## 🏗 Infrastructure as Code (IaC) & GitOps

The local Kubernetes environment is managed via a combination of Terraform and ArgoCD. Terraform bootstraps the local `Kind` cluster and installs ArgoCD. ArgoCD then takes over via the "App of Apps" pattern to deploy the Confluent Operator, Kafka cluster, and Microservices.

### Tearing Down and Rebuilding the Environment

To tear down and rebuild the environment from a clean slate, use Terraform from the `terraform/` directory:

1. **Destroy the environment:**
   ```bash
   cd terraform
   terraform destroy -auto-approve
   ```
2. **Rebuild the environment:**
   ```bash
   terraform apply -auto-approve
   ```
   *(Note: Once Terraform finishes, ArgoCD will automatically begin syncing the GitOps definitions from this repository to deploy the workloads).*

### Confluent CLI Integration

To interact with the deployed Kafka KRaft cluster running inside Kubernetes, you can use the Confluent CLI.

1. **Port-forward the Kafka bootstrap server:**
   ```bash
   kubectl port-forward svc/kafka-cluster-kafka-bootstrap -n default 9092:9092
   ```

2. **Check cluster status using Confluent CLI:**
   ```bash
   # List topics
   confluent kafka topic list --url http://localhost:9092

   # Consume messages from the orders topic
   confluent kafka topic consume orders --from-beginning --url http://localhost:9092
   ```

## 🌩 Chaos Engineering Observation

When running the advanced Kafka network chaos experiments (latency, packet loss, and bandwidth throttling), you can observe the network disruption via the **Chaos Mesh Dashboard** (`http://localhost:2333`) or by checking the application logs for timeouts and retries.

*(Note: Raw `container_network_*` PromQL queries like packet drops or bandwidth throughput will return "No Data" in this specific Minikube sandbox because the default `cAdvisor` configuration does not attach `pod` labels to network interface metrics, tracking only the root node `id="/"`).*

- **Observe CPU Stress:**
  ```promql
  rate(container_cpu_usage_seconds_total{namespace="default", pod=~"order-service.*|inventory-service.*"}[1m])
  ```
- **Observe Memory Stress:**
  ```promql
  container_memory_working_set_bytes{namespace="default", pod=~"order-service.*|inventory-service.*"}
  ```
*(Note: Application latency should be observed via HTTP/Kafka application-level metrics or timeouts in logs, rather than raw container network metrics).*

## ⚠️ Known Issues

### Grafana Dashboards "No Data" in Minikube
When running the `kube-prometheus-stack` locally in Minikube, the default Grafana dashboard "Kubernetes / Compute Resources / Namespace (Workloads)" may show "No Data" for CPU and Memory metrics.
- **Root Cause**: The default Prometheus recording rules require `image!=""` and `container!=""` labels. However, Minikube's `containerd` cAdvisor drops these labels for root pod cgroups. The dashboard's hardcoded JSON filters cause the metrics to be filtered out entirely. Furthermore, Prometheus fails to scrape Minikube's `kubelet` by default due to self-signed TLS certificates.
- **Workaround**: Apply a patch to the `kubelet.serviceMonitor` via your `values.yaml` to skip TLS verification and explicitly inject the missing labels using `cAdvisorMetricRelabelings`:
  ```yaml
  kubelet:
    serviceMonitor:
      insecureSkipVerify: true
      cAdvisorMetricRelabelings:
        - sourceLabels: [pod]
          targetLabel: container
          action: replace
        - sourceLabels: [pod]
          targetLabel: image
          action: replace
  ```
