# Implementation Plan: Migrate Inventory Service to Java (Quarkus)

## Phase: Project Scaffolding
- [x] Task: Initialize Quarkus Maven project for `inventory-service` [e10148f]
    - [x] Generate project structure with required extensions (RESTEasy, Micrometer/Prometheus, Kafka Client)
    - [x] Remove old Python files (`app.py`, `consumer.py`, `requirements.txt`)
- [x] Task: Conductor - User Manual Verification 'Project Scaffolding' (Protocol in workflow.md) [checkpoint: 110f569]

## Phase: Application Logic Implementation
- [x] Task: Implement Health Check Endpoint [8025de4]
    - [x] Write tests for `/healthz` endpoint
    - [x] Implement REST resource returning 200 OK
- [x] Task: Implement Kafka Consumer [179e92a]
    - [x] Create Kafka consumer service using Native Apache Kafka Client
    - [x] Parse incoming JSON order payloads
    - [x] Implement structured JSON logging
- [x] Task: Implement Prometheus Metrics [6e6ab50]
    - [x] Expose `/metrics` endpoint using Micrometer Prometheus extension
    - [x] Increment `inventory_events_total` counter upon processing orders
    - [x] Increment `inventory_errors_total` on parsing/processing failures
- [x] Task: Conductor - User Manual Verification 'Application Logic Implementation' (Protocol in workflow.md) [checkpoint: 2590631]

## Phase: Dockerization and Infrastructure Updates
- [x] Task: Write multi-stage Dockerfile for Quarkus Maven build [4100ea6]
    - [x] Create `Dockerfile` in `inventory-service` to build and package application
- [x] Task: Update `docker-compose.yml` to reflect new container build and ports [0aa53b8]
- [x] Task: Conductor - User Manual Verification 'Dockerization and Infrastructure Updates' (Protocol in workflow.md) [checkpoint: 32974db]
