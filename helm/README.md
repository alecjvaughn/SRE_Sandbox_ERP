# Infrastructure Helm Charts

This directory contains the Helm charts and configurations for deploying the SRE Sandbox ERP project.

## Kafka Cluster Installation
We use the Bitnami Kafka Helm chart in KRaft mode.

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install kafka-cluster bitnami/kafka -f helm/kafka-cluster/values.yaml --namespace default
```

## Observability Stack Installation
We use the `kube-prometheus-stack` to provision Prometheus and Grafana.

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install observability prometheus-community/kube-prometheus-stack -f helm/observability/values.yaml --namespace monitoring --create-namespace
```
