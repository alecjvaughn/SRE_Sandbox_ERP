# Implementation Plan: Kubernetes Active-Active Infrastructure

## Phase: Helm Chart Scaffolding
- [ ] Task: Scaffold base Helm charts
    - [ ] Create `helm/order-service` directory structure
    - [ ] Create `helm/inventory-service` directory structure
- [ ] Task: Configure Kafka dependency
    - [ ] Define a `helm/kafka-cluster` configuration or document Bitnami Kafka installation steps
- [ ] Task: Conductor - User Manual Verification 'Helm Chart Scaffolding' (Protocol in workflow.md)

## Phase: Application Manifest Definitions
- [ ] Task: Implement `order-service` Helm templates
    - [ ] Define Deployment with Soft Anti-Affinity rules
    - [ ] Define CPU/Memory requests and limits
    - [ ] Define Service, liveness, and readiness probes (`/healthz`)
- [ ] Task: Implement `inventory-service` Helm templates
    - [ ] Define Deployment with Soft Anti-Affinity rules
    - [ ] Define CPU/Memory requests and limits
    - [ ] Define liveness and readiness probes (`/healthz`)
- [ ] Task: Conductor - User Manual Verification 'Application Manifest Definitions' (Protocol in workflow.md)

## Phase: Ingress and State Configuration
- [ ] Task: Configure Ingress and Load Balancer
    - [ ] Define Ingress resource for `order-service`
    - [ ] Provide DigitalOcean Load Balancer annotations/documentation
- [ ] Task: Configure Kafka Hard Anti-Affinity
    - [ ] Define custom `values.yaml` for `bitnami/kafka` applying `requiredDuringSchedulingIgnoredDuringExecution`
- [ ] Task: Validate Helm Charts
    - [ ] Run `helm lint` on application charts
    - [ ] Run `helm template` to verify output structure
- [ ] Task: Conductor - User Manual Verification 'Ingress and State Configuration' (Protocol in workflow.md)
