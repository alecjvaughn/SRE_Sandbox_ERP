# SRE Sandbox: Cost-Aware Traffic Plane & Egress Optimizer (Azure AKS)

This repository has evolved from an ERP simulation into an advanced cloud-native data plane simulation targeting the optimization of ingress and egress traffic for Confluent Kafka workloads. Built for Azure Kubernetes Service (AKS), the project focuses on network routing, failover mechanisms, and bandwidth cost efficiency across Availability Zones.

## 🏗 Architecture

- **Custom Go Edge Proxy**: Intelligent Layer 4 TCP proxy for dynamically routing Kafka traffic based on AZ health and cost.
- **Go Order Service**: Active-active REST gateway handling incoming messages.
- **Java (Quarkus) Inventory Service**: Background event processing consumer listening to Kafka.
- **Kafka**: Confluent Kafka broker managing the event bus.
- **Observability**: Prometheus scraping and Grafana dashboards tracking Network I/O and Egress Costs.
- **Resilience**: Chaos Mesh for injecting network degradation, latency, and failovers.

## 🚀 Current Status

- **Phase 1-7**: Core Microservices, Local Sandbox, Observability, Chaos (Completed).
- **Phase 8**: Infrastructure as Code (Azure AKS Provisioning) (Completed).
- **Phase 9**: Custom Go Edge Proxy & Egress Cost Analysis (Completed).

## 🛠 Prerequisites

- [Go](https://golang.org/doc/install) (1.24+)
- [Java 21+](https://jdk.java.net/) & Maven
- [Docker](https://docs.docker.com/get-docker/) & Docker Hub Account
- [Terraform](https://www.terraform.io/)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)

## 📂 Project Structure

```
.
├── PLAN_OVERVIEW.md    # High-level roadmap and architectural goals
├── argocd/             # GitOps definitions (App of Apps)
├── conductor/          # Tracks, specification, and agent workflow plans
├── edge-proxy/         # Custom Go Layer 4 TCP Traffic Proxy
├── inventory-service/  # Java Quarkus Kafka consumer and metrics server
├── order-service/      # Go REST API for handling orders
├── terraform/          # Infrastructure as Code (Azure configs)
└── README.md           # This evolving project documentation
```

## 💻 Getting Started

The environment is built to be an ephemeral Azure AKS cluster managed via Terraform and ArgoCD GitOps.

1. **Deploy the infrastructure:**
   ```bash
   make azure-apply
   ```

2. **Wait for ArgoCD to sync applications:**
   ArgoCD will automatically start syncing all microservices. You can monitor the progress with:
   ```bash
   kubectl get applications -n argocd
   kubectl get pods -n default
   ```

3. **Verify the services:**
   You can port-forward the services to your local machine for testing:
   ```bash
   kubectl port-forward svc/edge-proxy -n default 8080:8080 &
   kubectl port-forward svc/inventory-service -n default 8081:8080 &
   ```
   - Edge Proxy (routing to Order Service): `http://localhost:8080/order`
   - Inventory Service Health Check: `http://localhost:8081/healthz`

4. **Send a test order:**
   ```bash
   curl -X POST http://localhost:8080/order -d '{"order_id": "test-1", "item": "Widget", "qty": 5}'
   ```

## 🏗 Infrastructure as Code (IaC) & GitOps

The Azure Kubernetes environment is managed via a combination of Terraform and ArgoCD. Terraform bootstraps the AKS cluster and installs ArgoCD. ArgoCD then takes over via the "App of Apps" pattern to deploy the Confluent Operator, Kafka cluster, and Microservices.

### Tearing Down and Rebuilding the Environment

To tear down and rebuild the environment from a clean slate, use the provided `make` targets:

1. **Destroy the environment:**
   ```bash
   make nuke
   ```
2. **Rebuild the environment:**
   ```bash
   make azure-apply
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
