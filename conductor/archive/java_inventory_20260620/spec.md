# Specification: Migrate Inventory Service to Java (Quarkus)

## Overview
This track replaces the existing Python-based `inventory-service` with a Java-based implementation using the Quarkus framework. The new service must maintain feature parity with the previous version, including Kafka event consumption, Prometheus metrics, and K8s health checks.

## Functional Requirements
- Rewrite the `inventory-service` using Java and Quarkus.
- Build tool: Maven.
- Implement a Kafka consumer using the Native Apache Kafka Client to consume from the `orders` topic.
- Maintain stock decrease logic based on incoming order payloads.
- Expose a `/healthz` endpoint returning HTTP 200 OK.
- Expose a `/metrics` endpoint with Prometheus metrics, tracking the `inventory_events_total` counter.

## Non-Functional Requirements
- Update the Dockerfile to use a multi-stage Maven build suitable for Quarkus.
- Maintain compatibility with the existing `docker-compose.yml` infrastructure.
- Adhere to the general product guidelines (structured JSON logging).

## Out of Scope
- No database integration (stock updates are purely mocked via metrics).
- No changes to `order-service`.
