# Infrastructure Helm Charts

This directory contains the Helm charts and configurations for deploying the SRE Sandbox ERP project.

## Kafka Cluster Installation
We use the Bitnami Kafka Helm chart in KRaft mode.

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install kafka-cluster bitnami/kafka -f helm/kafka-cluster/values.yaml --namespace default
```
