# Implementation Plan: Observability Stack

## Phase: Helm Chart Preparation
- [ ] Task: Scaffold `kube-prometheus-stack` values
    - [ ] Create `helm/observability/values.yaml` configuration file for the stack
    - [ ] Add instructions for the Prometheus Community Helm repo to `helm/README.md`
- [ ] Task: Conductor - User Manual Verification 'Helm Chart Preparation' (Protocol in workflow.md)

## Phase: ServiceMonitor Configuration
- [ ] Task: Define `order-service` ServiceMonitor
    - [ ] Update `helm/order-service/templates/servicemonitor.yaml` or create a standalone manifest
- [ ] Task: Define `inventory-service` ServiceMonitor
    - [ ] Update `helm/inventory-service/templates/servicemonitor.yaml` or create a standalone manifest
- [ ] Task: Define `kafka-cluster` ServiceMonitor
    - [ ] Configure `values.yaml` for Bitnami Kafka to enable JMX exporter and `metrics.serviceMonitor.enabled=true`
- [ ] Task: Conductor - User Manual Verification 'ServiceMonitor Configuration' (Protocol in workflow.md)

## Phase: Dashboard as Code
- [ ] Task: Create Grafana Dashboard ConfigMap
    - [ ] Define JSON dashboard structure for RPS, Latency, CPU, Memory, and Lag
    - [ ] Wrap JSON in a Kubernetes ConfigMap labeled for Grafana auto-discovery (`grafana_dashboard: "1"`)
- [ ] Task: Validate Manifests
    - [ ] Run `helm lint` and template checks to validate the new configuration
- [ ] Task: Conductor - User Manual Verification 'Dashboard as Code' (Protocol in workflow.md)
