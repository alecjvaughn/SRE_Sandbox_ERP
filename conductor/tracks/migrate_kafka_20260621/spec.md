# Specification: Migrate to Confluent Kafka

## Overview
Migrate the local Minikube Kafka infrastructure from the deprecated/broken Bitnami Helm chart to the official Confluent Helm chart (`confluentinc/cp-kafka`). This is required because Bitnami has restricted their free image catalog on Docker Hub, causing `ErrImagePull` errors on all `bitnami/kafka` tags.

## Functional Requirements
- Remove the existing Bitnami Kafka Helm deployment and resources.
- Deploy a new Kafka cluster using the official Confluent Helm chart.
- The new cluster must use a 3-node topology with KRaft (Zookeeper-less) to mirror the previous architecture.
- Re-create necessary topics (e.g., `orders`) since the old cluster had no persistent data to migrate (it was failing at initialization).
- Ensure the Go microservices (`order-service`, `inventory-service`) are updated to point to the new Confluent Kafka broker URLs if the service names change.

## Out of Scope
- Migrating physical data/volumes (previous deployment failed before any data was written).
- Modifying the Chaos Mesh configuration (will be re-tested in a separate track phase).
