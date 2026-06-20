# Implementation Plan: Implement Go Order and Python Inventory Services with Kafka Event Bus

This plan outlines the step-by-step TDD tasks to build the core services and verify their event-driven integration.

## Phase 1: Go Order Service

- [x] Task: Implement HTTP Endpoints and Basic Validation (TDD) [932256a]
    - [ ] Write unit tests for POST `/order` payload validation and GET `/healthz` endpoint
    - [ ] Implement order-service routing, request parsing, health checks, and mock response
    - [ ] Run test suite to verify tests pass
- [ ] Task: Implement Kafka Producer (TDD)
    - [ ] Write unit tests for Kafka event serialization and transmission
    - [ ] Implement Kafka producer setup using `franz-go` and integrate it into POST `/order` handler
    - [ ] Run test suite to verify event production works
- [ ] Task: Add Observability and Docker Containerization
    - [ ] Write unit tests for `/metrics` endpoint and custom metrics incrementing
    - [ ] Implement Prometheus metrics collection for `erp_orders_total`
    - [ ] Write a multi-stage Dockerfile for `order-service`
    - [ ] Run test suite and verify docker build completes successfully
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Go Order Service' (Protocol in workflow.md)

## Phase 2: Python Inventory Service

- [ ] Task: Implement Health Server and Metric Server (TDD)
    - [ ] Write unit tests for GET `/healthz` check on port 8081
    - [ ] Implement HTTP server routing for health check and start Prometheus server on port 8000
    - [ ] Run test suite to verify server behavior
- [ ] Task: Implement Kafka Consumer and Event Handlers (TDD)
    - [ ] Write unit/integration tests for Kafka message consuming and stock decrement logs
    - [ ] Implement consumer thread using `kafka-python` or `confluent-kafka` and event parsing logic
    - [ ] Run test suite to verify consumer reads and updates metric counters
- [ ] Task: Python Containerization
    - [ ] Write a Dockerfile for `inventory-service`
    - [ ] Verify docker build completes successfully
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Python Inventory Service' (Protocol in workflow.md)

## Phase 3: Integration and Docker Compose Local Sandbox

- [ ] Task: Orchestrate Sandbox with Docker Compose
    - [ ] Write a `docker-compose.yml` defining services for order-service, inventory-service, zookeeper, and kafka
    - [ ] Configure network links and ports between services
- [ ] Task: End-to-End Pipeline Verification
    - [ ] Start the compose cluster and wait for brokers to initialize
    - [ ] Send sample curl POST requests to order-service and inspect inventory consumer logs for stock mutations
    - [ ] Verify Prometheus scrape targets and `/metrics` output
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Integration and Docker Compose Local Sandbox' (Protocol in workflow.md)
