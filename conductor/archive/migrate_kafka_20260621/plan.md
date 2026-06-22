# Implementation Plan: Migrate to Confluent Kafka

## Phase: Pre-Migration Cleanup [checkpoint: 0cdaab7]
- [x] Task: Uninstall the existing Bitnami Kafka Helm release 2ff0f96
- [x] Task: Delete any lingering PVCs, Services, or ConfigMaps associated with the old Kafka cluster 0f5b895
- [x] Task: Conductor - User Manual Verification 'Pre-Migration Cleanup' (Protocol in workflow.md) 0cdaab7

## Phase: Confluent Kafka Deployment [checkpoint: IaC-Track]
- [x] Task: Add the Confluent CFK Helm repository and configure it for a 3-node KRaft topology 9800675
- [x] Task: Write updated CFK CRDs for the new Kafka cluster 86d0546
- [x] Task: Deploy the CFK Operator and Confluent Kafka cluster to the Minikube environment c9d6829
- [x] Task: Wait for pods to become healthy and verify cluster formation (Completed via IaC Track)
- [x] Task: Conductor - User Manual Verification 'Confluent Kafka Deployment' (Completed via IaC Track)

## Phase: Application Integration & IaC Alignment [checkpoint: TBD]
- [x] Task: Identify the correct Kafka bootstrap server URLs for the new Confluent cluster 037b30c
- [x] Task: Inject the Kafka URLs as environment variables into `order-service` and `inventory-service` Helm charts (`values.yaml` / `deployment.yaml`) 037b30c
- [x] Task: Commit the changes to trigger an ArgoCD GitOps sync and verify pod rollout 037b30c
- [x] Task: Run test requests to ensure the microservices successfully produce and consume events against the new cluster ab0951f
- [ ] Task: Conductor - User Manual Verification 'Application Integration & IaC Alignment' (Protocol in workflow.md)
