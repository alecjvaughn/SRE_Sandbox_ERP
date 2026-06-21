# Implementation Plan: Kubernetes Active-Active Infrastructure

## Phase: Helm Chart Scaffolding
- [x] Task: Scaffold base Helm charts
    - [x] Create `helm/order-service` directory structure
    - [x] Create `helm/inventory-service` directory structure
- [x] Task: Configure Kafka dependency
    - [x] Define a `helm/kafka-cluster` configuration or document Bitnami Kafka installation steps
- [x] Task: Conductor - User Manual Verification 'Helm Chart Scaffolding' (Protocol in workflow.md)

## Phase: Application Manifest Definitions
- [x] Task: Implement `order-service` Helm templates
    - [x] Define Deployment with Soft Anti-Affinity rules
    - [x] Define CPU/Memory requests and limits
    - [x] Define Service, liveness, and readiness probes (`/healthz`)
- [x] Task: Implement `inventory-service` Helm templates
    - [x] Define Deployment with Soft Anti-Affinity rules
    - [x] Define CPU/Memory requests and limits
    - [x] Define liveness and readiness probes (`/healthz`)
- [x] Task: Conductor - User Manual Verification 'Application Manifest Definitions' (Protocol in workflow.md)

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
