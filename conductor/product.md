# Initial Concept

Active-Active SRE Sandbox with Go Mock ERP, Java Inventory Service, and Confluent Kafka on K8s (using Obsidian Design Doc)

# Product Guide: SRE Sandbox with ERP

## Vision
The SRE Sandbox with ERP is a lightweight, cost-effective simulation platform that mirrors cloud-native enterprise enterprise resource planning (ERP) architectures. By decoupling transaction processing (Go mock service) from state/inventory mutation (Java mock service) via a Kafka event stream, this project provides a safe, highly observable environment to test infrastructure resiliency, chaos engineering, and automated self-healing mechanisms.

## Core Features
1. **Active-Active Go Order Service**: High-throughput REST API that accepts transaction payloads and produces events to Kafka, running with replica anti-affinity across multiple Kubernetes worker nodes.
2. **Java Inventory Service**: Event consumer subscribing to Kafka, processing stock mutations, and exposing service metrics via Quarkus.
3. **Kafka Event Fabric**: Confluent-based partitioned topic streaming to handle concurrent microservice communications.
4. **Chaos Engineering Engine**: Chaos Mesh integrations for scheduled pod evictions, network latency injections, and memory spikes.
5. **Observability Stack**: Prometheus metrics collection and real-time Grafana dashboards displaying latency, throughput, compute saturation, and self-healing events.

## Target Audience
- **SRE & DevOps Engineers**: To test and demo cloud-native resilience practices.
- **Software Engineers**: To understand decoupled event-driven architectures and API design under stress.
- **Interviewers / Recruiters**: Serving as an advanced, production-grade systems architecture portfolio piece.
