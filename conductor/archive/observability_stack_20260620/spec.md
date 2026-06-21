# Specification: Observability Stack

## Overview
This track focuses on deploying a comprehensive observability stack to the Kubernetes cluster to provide real-time metrics and dashboards for the SRE Sandbox ERP project. This will be critical for monitoring application health and validating the upcoming Chaos Engineering experiments.

## Functional Requirements
- **Deployment Strategy**: Utilize the community `kube-prometheus-stack` Helm chart to deploy Prometheus, Prometheus Operator, and Grafana.
- **Metrics Collection**: Configure custom `ServiceMonitors` to scrape metrics from:
  - Go `order-service`
  - Java `inventory-service`
  - Confluent Kafka Cluster
- **Dashboards**: Provision a default Grafana dashboard as Code (via ConfigMap or Helm values) that prominently displays:
  - API Throughput (RPS) and Request Latency
  - CPU & Memory Resource Utilization for all components
  - Kafka Consumer Lag for event streaming health

## Non-Functional Requirements
- Dashboards must be provisioned declaratively via GitOps/Helm.
- The stack must be lightweight enough to run alongside the application within local environments (e.g., Minikube).

## Out of Scope
- Alertmanager configuration (e.g., PagerDuty, Slack alerts).
- Distributed tracing (Jaeger/Zipkin).
- Log aggregation (Fluentbit/Elasticsearch/Loki).
