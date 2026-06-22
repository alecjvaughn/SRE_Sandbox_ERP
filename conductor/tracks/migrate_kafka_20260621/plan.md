# Implementation Plan: Migrate to Confluent Kafka

## Phase: Pre-Migration Cleanup
- [x] Task: Uninstall the existing Bitnami Kafka Helm release 2ff0f96
- [~] Task: Delete any lingering PVCs, Services, or ConfigMaps associated with the old Kafka cluster
- [ ] Task: Conductor - User Manual Verification 'Pre-Migration Cleanup' (Protocol in workflow.md)

## Phase: Confluent Kafka Deployment
- [ ] Task: Add the Confluent Helm repository and configure it for a 3-node KRaft topology
- [ ] Task: Write updated Helm `values.yaml` for the new Kafka cluster
- [ ] Task: Deploy the Confluent Kafka cluster to the Minikube environment
- [ ] Task: Wait for pods to become healthy and verify cluster formation
- [ ] Task: Conductor - User Manual Verification 'Confluent Kafka Deployment' (Protocol in workflow.md)

## Phase: Application Integration
- [ ] Task: Identify if the new Kafka bootstrap server URLs differ from the old ones
- [ ] Task: Update the `values.yaml` or ConfigMaps for `order-service` and `inventory-service` with the new Kafka URLs
- [ ] Task: Restart the microservices to pick up the new configuration
- [ ] Task: Run test requests to ensure the Go applications can successfully produce and consume events
- [ ] Task: Conductor - User Manual Verification 'Application Integration' (Protocol in workflow.md)
