# Implementation Plan: Chaos Engineering Validation

## Phase: Setup Chaos Mesh [checkpoint: 049e84f]
- [x] Task: Deploy Chaos Mesh to Local Minikube 8f14a91
    - [x] Add Chaos Mesh Helm repository
    - [x] Install Chaos Mesh via Helm
    - [x] Verify Chaos Mesh pods are running
- [x] Task: Conductor - User Manual Verification 'Setup Chaos Mesh' (Protocol in workflow.md) 049e84f

## Phase: Chaos Experiments Implementation [checkpoint: 8e9573e]
- [x] Task: Define Pod Eviction Experiment 5a0d882
    - [x] Create `PodChaos` manifest for randomly terminating `order-service` pods
    - [x] Create `PodChaos` manifest for randomly terminating `inventory-service` pods
- [x] Task: Define Network Latency Experiment a4f5f6f
    - [x] Create `NetworkChaos` manifest to inject latency between services and Kafka
- [x] Task: Define CPU/Memory Stress Experiment 400ee37
    - [x] Create `StressChaos` manifest to simulate resource exhaustion on application pods
- [x] Task: Conductor - User Manual Verification 'Chaos Experiments Implementation' (Protocol in workflow.md) 8e9573e

## Phase: Validation and Observability [checkpoint: 706c72c]
- [x] Task: Execute and Monitor Local Experiments a54910c
    - [x] Apply chaos manifests to the minikube cluster
    - [x] Monitor `/healthz` endpoints for recovery
    - [x] Validate system recovery via Grafana/Prometheus dashboards
- [x] Task: Conductor - User Manual Verification 'Validation and Observability' (Protocol in workflow.md) 706c72c

## Phase: UI Dashboard Metrics Validation [checkpoint: d34c7e8]
- [x] Task: Fix cAdvisor label filtering for Grafana d34c7e8
    - [x] Identify "No Data" issue caused by missing container labels in Minikube
    - [x] Configure Prometheus `cAdvisorMetricRelabelings` override via Helm values
    - [x] Validate CPU/Memory Chaos telemetry appears correctly on default dashboards
    - [x] Document workaround in `README.md`
- [x] Task: Conductor - User Manual Verification 'UI Dashboard Validation' d34c7e8

## Phase: Advanced Kafka Network Chaos [checkpoint: 728b9d0]
- [x] Task: Re-activate chaos schedulers and define advanced network faults 728b9d0
    - [x] Update `NetworkChaos` manifests with `Schedule` CRDs or cron schedules for randomized attacks
    - [x] Inject packet loss rules between `order-service` and Kafka
    - [x] Inject bandwidth throttling rules between `order-service` and Kafka
    - [x] Validate 500ms latency rules
- [x] Task: Conductor - User Manual Verification 'Advanced Kafka Network Chaos' (Protocol in workflow.md) 728b9d0

## Phase: Post-Migration Verification [checkpoint: TBD]
- [ ] Task: Re-verify system recovery against new Confluent Kafka cluster
    - [ ] Run test orders to ensure UI and metrics are displaying recovery behavior
- [ ] Task: Conductor - User Manual Verification 'Post-Migration Verification'
