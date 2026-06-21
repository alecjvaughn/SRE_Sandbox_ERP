# Specification: Kubernetes Active-Active Infrastructure

## 1. Overview
Deploy the SRE Sandbox ERP architecture into a DigitalOcean Kubernetes (DOKS) cluster using Helm Charts. The deployment will enforce an Active-Active topology using replica anti-affinity rules to ensure high availability.

## 2. Functional Requirements
- **Go Order Service Deployment**: Create a Helm chart to deploy the Go REST API.
- **Java Inventory Service Deployment**: Create a Helm chart to deploy the Java (Quarkus) event consumer.
- **Kafka Cluster Deployment**: Deploy a Kafka broker cluster using Helm (e.g., bitnami/kafka).
- **Ingress & Load Balancing**: Configure a DigitalOcean Load Balancer and Ingress controller (e.g., ingress-nginx) to route external traffic to the Order Service.

## 3. Non-Functional Requirements
- **High Availability (Anti-Affinity)**:
  - **Applications (`order-service`, `inventory-service`)**: Use Soft Anti-Affinity (`preferredDuringSchedulingIgnoredDuringExecution`) to encourage distribution across worker nodes.
  - **Stateful Infrastructure (Kafka)**: Use Hard Anti-Affinity (`requiredDuringSchedulingIgnoredDuringExecution`) to strictly enforce node separation for Kafka brokers.
- **Resource Constraints**: Define CPU and Memory `requests` and `limits` in all Helm charts.
- **Health Probes**: Configure liveness and readiness probes pointing to `/healthz` for all application pods.

## 4. Acceptance Criteria
- [ ] Helm charts are created and version-controlled for all application components.
- [ ] Kafka is deployable via Helm with hard anti-affinity enabled.
- [ ] `order-service` and `inventory-service` are deployable via Helm with soft anti-affinity enabled.
- [ ] DigitalOcean Load Balancer and Ingress definitions are included and correctly route HTTP traffic.
- [ ] `helm lint` and `helm template` successfully validate all charts.

## 5. Out of Scope
- Implementation of the Observability Stack (Prometheus/Grafana).
- Implementation of Chaos Engineering experiments.
- Automated CI/CD pipeline integration (e.g., GitHub Actions, ArgoCD).
