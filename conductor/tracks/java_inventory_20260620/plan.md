# Implementation Plan: Migrate Inventory Service to Java (Quarkus)

## Phase: Project Scaffolding
- [x] Task: Initialize Quarkus Maven project for `inventory-service` [e10148f]
    - [x] Generate project structure with required extensions (RESTEasy, Micrometer/Prometheus, Kafka Client)
    - [x] Remove old Python files (`app.py`, `consumer.py`, `requirements.txt`)
- [x] Task: Conductor - User Manual Verification 'Project Scaffolding' (Protocol in workflow.md) [checkpoint: 110f569]

## Phase: Application Logic Implementation
- [ ] Task: Implement Health Check Endpoint
    - [ ] Write tests for `/healthz` endpoint
    - [ ] Implement REST resource returning 200 OK
- [ ] Task: Implement Kafka Consumer
    - [ ] Create Kafka consumer service using Native Apache Kafka Client
    - [ ] Parse incoming JSON order payloads
    - [ ] Implement structured JSON logging
- [ ] Task: Implement Prometheus Metrics
    - [ ] Expose `/metrics` endpoint using Micrometer Prometheus extension
    - [ ] Increment `inventory_events_total` counter upon processing orders
- [ ] Task: Conductor - User Manual Verification 'Application Logic Implementation' (Protocol in workflow.md)

## Phase: Dockerization and Infrastructure Updates
- [ ] Task: Write multi-stage Dockerfile for Quarkus Maven build
- [ ] Task: Update `docker-compose.yml` to reflect new container build and ports
- [ ] Task: Conductor - User Manual Verification 'Dockerization and Infrastructure Updates' (Protocol in workflow.md)
