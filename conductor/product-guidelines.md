# Product Guidelines: SRE Sandbox with ERP

These guidelines establish standards for coding, configuration, documentation, and operational workflows inside this repository.

## 1. Reliability & Self-Healing Guidelines
- **Health Checking**: Every microservice must expose a `/healthz` HTTP endpoint. The response must return `200 OK` under normal operations.
- **Resource Constraints**: All Kubernetes deployment templates must explicitly define CPU and memory `requests` and `limits`.
- **Fault Isolation**: No database or shared state is permitted between microservices directly. All communications must go through the Kafka event stream or public API endpoints.
- **Replica Anti-Affinity**: In production mode, critical services must maintain a minimum replica count of 2 with anti-affinity rules to prevent scheduling on the same physical node.

## 2. Telemetry & Observability Guidelines
- **Metrics Standard**: All services must expose metrics on a `/metrics` endpoint scrapable by Prometheus.
- **Log Formatting**: Application logs should be output in structured JSON format (containing level, timestamp, message, and caller metadata).
- **Dashboard Consistency**: Real-time dashboards in Grafana must use PromQL queries mapping usage percentage against hard resource limits.

## 3. Configuration & Deployment Guidelines
- **Infrastructure as Code (IaC)**: Kubernetes manifests and Helm value overrides must be kept declarative and checked into version control.
- **Secrets Management**: No API keys, credentials, or TLS certificates may be committed to this repository. Use Kubernetes secrets or environment variable placeholders.
- **Port Standards**:
  - Go Order Service API: port `8080` (container) / port `80` (load balancer)
  - Python Inventory Service: port `8000` (metrics) / port `8081` (health/control)
  - Kafka Broker: port `9092`
