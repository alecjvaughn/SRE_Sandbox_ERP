# SRE Sandbox ERP

This repository contains the core transactional and event-driven foundation of the SRE Sandbox. It simulates an ERP architecture consisting of a Go-based order service and a Python-based inventory service communicating via a Kafka event bus.

## 🏗 Architecture

- **Go Order Service**: Active-active REST gateway handling incoming orders.
- **Python Inventory Service**: Background event processing consumer listening to Kafka.
- **Kafka**: Confluent Kafka broker managing the event bus.
- **Observability**: Prometheus scraping and Grafana dashboards.
- **Resilience**: Chaos Mesh for injecting latency, failures, and evictions.

## 🚀 Current Status

- **Phase 1**: Implementing Go Order Service endpoints and TDD scaffolding (In Progress).

## 🛠 Prerequisites

- [Go](https://golang.org/doc/install) (1.22+)
- [Python](https://www.python.org/downloads/) (3.11+)
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose

## 📂 Project Structure

```
.
├── PLAN_OVERVIEW.md    # High-level roadmap and architectural goals
├── conductor/          # Tracks, specification, and agent workflow plans
├── order-service/      # Go REST API for handling orders
└── README.md           # This evolving project documentation
```

## 💻 Getting Started

*(Instructions will be added here as the services are dockerized and the local sandbox is finalized.)*
