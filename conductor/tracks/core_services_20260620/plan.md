# Implementation Plan: Implement Go Order and Python Inventory Services with Kafka Event Bus

This plan outlines the step-by-step TDD tasks to build the core services and verify their event-driven integration.

## Phase 1: Go Order Service [checkpoint: a863b8d]

- [x] Task: Implement HTTP Endpoints and Basic Validation (TDD) [932256a]
    - [x] Write unit tests for POST `/order` payload validation and GET `/healthz` endpoint
    - [x] Implement order-service routing, request parsing, health checks, and mock response
    - [x] Run test suite to verify tests pass
- [x] Task: Implement Kafka Producer (TDD) [01c5b13]
    - [x] Write unit tests for Kafka event serialization and transmission
    - [x] Implement Kafka producer setup using `franz-go` and integrate it into POST `/order` handler
    - [x] Run test suite to verify event production works
- [x] Task: Add Observability and Docker Containerization [1e5c6fb]
    - [x] Write unit tests for `/metrics` endpoint and custom metrics incrementing
    - [x] Implement Prometheus metrics collection for `erp_orders_total`
    - [x] Write a multi-stage Dockerfile for `order-service`
    - [x] Run test suite and verify docker build completes successfully
- [x] Task: Conductor - User Manual Verification 'Phase 1: Go Order Service' (Protocol in workflow.md) [a863b8d]

## Phase 2: Python Inventory Service [checkpoint: 7217850]

- [x] Task: Implement Health Server and Metric Server (TDD) [9037c29]
    - [x] Write unit tests for GET `/healthz` check on port 8081
    - [x] Implement HTTP server routing for health check and start Prometheus server on port 8000
    - [x] Run test suite to verify server behavior
- [x] Task: Implement Kafka Consumer and Event Handlers (TDD) [97dbb53]
    - [x] Write unit/integration tests for Kafka message consuming and stock decrement logs
    - [x] Implement consumer thread using `kafka-python` or `confluent-kafka` and event parsing logic
    - [x] Run test suite to verify consumer reads and updates metric counters
- [x] Task: Python Containerization [4a9d51c]
    - [x] Write a Dockerfile for `inventory-service`
    - [x] Verify docker build completes successfully
- [x] Task: Conductor - User Manual Verification 'Phase 2: Python Inventory Service' (Protocol in workflow.md) [7217850]

## Phase 3: Integration and Docker Compose Local Sandbox

- [ ] Task: Orchestrate Sandbox with Docker Compose
    - [ ] Write a `docker-compose.yml` defining services for order-service, inventory-service, zookeeper, and kafka
    - [ ] Configure network links and ports between services
- [ ] Task: End-to-End Pipeline Verification
    - [ ] Start the compose cluster and wait for brokers to initialize
    - [ ] Send sample curl POST requests to order-service and inspect inventory consumer logs for stock mutations
    - [ ] Verify Prometheus scrape targets and `/metrics` output
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Integration and Docker Compose Local Sandbox' (Protocol in workflow.md)

## Phase: Review Fixes
- [x] Task: Apply review suggestions [f5789e0]
