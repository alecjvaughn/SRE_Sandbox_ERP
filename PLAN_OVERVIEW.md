# Plan Overview: SRE Sandbox with ERP

This document outlines the high-level roadmap, current active track, and status of tasks for the SRE Sandbox with ERP.

## Architectural Diagram

```mermaid
graph TD
    LB[DO Load Balancer] -->|Round-Robin| Node1[DOKS Node 1]
    LB -->|Round-Robin| Node2[DOKS Node 2]
    Node1 -->|Active-Active| Go1[order-service-a Pod]
    Node2 -->|Active-Active| Go2[order-service-b Pod]
    Go1 -->|Produce event| Kafka[Confluent Kafka Broker]
    Go2 -->|Produce event| Kafka
    Kafka -->|Consume event| Java1[inventory-java-a Pod]
    Kafka -->|Consume event| Java2[inventory-java-b Pod]
    
    Chaos[Chaos Mesh] -->|Inject Failures| Java1
    Chaos -->|Inject Failures| Java2
    
    Prom[Prometheus] -->|Scrape /metrics| Go1
    Prom -->|Scrape /metrics| Go2
    Prom -->|Scrape /metrics| Java1
    Prom -->|Scrape /metrics| Java2
    Prom -->|Scrape /metrics| Kafka
    
    Grafana[Grafana Dashboard] -->|Query| Prom
```

## Requirements & Solutions

### Core Requirements
- **Active-Active REST Gateway**: Go `order-service` with anti-affinity running across multiple nodes.
- **Background Event Processing**: Java `inventory-service` (Quarkus) consuming events via Kafka.
- **Chaos Resilience**: Pod eviction recovery, latency, and memory limit handling.
- **Full Observability**: Live Prometheus scraping and Grafana metrics dashboard.

### Implementation Goals
1. Implement clean REST API handlers in Go and structured event handlers in Python.
2. Ensure strict test-driven development (TDD) cycle for all components.
3. Package application components as microservice containers.
4. Establish local development environment via Docker Compose to validate flows before cloud deploy.

---

## Completed Tracks
### Track: Implement Go Order and Python Inventory Services with Kafka Event Bus
Status: `Completed`
- [x] Phase 1: Go Order Service
- [x] Phase 2: Python Inventory Service
- [x] Phase 3: Integration and Docker Compose Local Sandbox

### Track: Migrate Inventory Service to Java (Quarkus)
Status: `Completed`
- [x] Phase 1: Project Scaffolding
- [x] Phase 2: Application Logic Implementation
- [x] Phase 3: Dockerization and Infrastructure Updates
