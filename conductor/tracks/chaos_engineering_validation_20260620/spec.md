# Specification: Chaos Engineering Validation

## Overview
This track introduces chaos engineering experiments using Chaos Mesh to validate the resilience and self-healing capabilities of the SRE Sandbox ERP. The experiments will be executed on the local `minikube` sandbox.

## Functional Requirements
- **Experiment Scope**: Configure and run the following Chaos Mesh experiments:
  - **Pod Eviction**: Randomly terminate pods in `order-service` and `inventory-service` to test Kubernetes ReplicaSet recovery.
  - **Network Latency**: Inject latency between `order-service`, `inventory-service`, and the Kafka broker to validate timeouts and retry mechanisms.
  - **CPU/Memory Stress**: Simulate resource exhaustion on pods to test Kubernetes resource limits and OOM recovery.
- **Success Criteria**: 
  - Automated health checks (`/healthz`) must eventually return 200 OK without manual intervention after failures.
  - System recovery must be visually observable via the Grafana/Prometheus dashboard metrics.

## Non-Functional Requirements
- Experiments must be defined declaratively as Kubernetes manifests and committed to version control.

## Out of Scope
- Network partition (isolation) experiments.
- Testing on Cloud Staging or Production environments.
