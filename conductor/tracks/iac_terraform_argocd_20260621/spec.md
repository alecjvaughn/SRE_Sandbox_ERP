# Specification: Define IaC using Terraform and ArgoCD

## Overview
This track defines the Infrastructure as Code (IaC) foundation for the SRE Sandbox project. The goal is to allow deploying the entire project from a clean slate using Terraform and ArgoCD, focusing exclusively on a local development environment (Kind).

## Functional Requirements
- **Terraform Automation**: Terraform will be used to bootstrap the local environment. It will deploy ArgoCD into the cluster.
- **ArgoCD GitOps**: ArgoCD will be configured using the "App of Apps" pattern to establish a scalable GitOps structure for managing all cluster workloads.
- **Workload Management**: ArgoCD (via the App of Apps) will be responsible for deploying the Confluent Kafka Operator, Kafka CRDs, and the microservices (`order-service`, `inventory-service`).
- **Confluent CLI Integration**: Provide instructions and commands for using the Confluent CLI ad-hoc during testing and monitoring.
- **Local State Management**: Terraform state will be stored locally (`terraform.tfstate`).
- **Clean Slate Redeployment**: The configuration must be capable of rebuilding the entire infrastructure and application stack from an empty Kind cluster.
- **Makefile Automation**: Provide a Makefile containing atomic commands for managing Docker images, Kubernetes interactions, and full environment teardown (nuke) commands.

## Non-Functional Requirements
- **Scalability**: While the target is Kind, the "App of Apps" ArgoCD structure should be designed so that it can easily be adapted to a cloud environment in the future.
- **Reproducibility**: The IaC must be fully declarative and reproducible.

## Acceptance Criteria
- [ ] A `terraform` directory is created containing the main Terraform configuration files.
- [ ] Terraform successfully bootstraps ArgoCD into the local Kind cluster.
- [ ] An ArgoCD "App of Apps" manifest structure is established.
- [ ] ArgoCD successfully deploys the Confluent Kafka operator, KRaft cluster, and microservices via GitOps.
- [ ] Confluent CLI testing and monitoring instructions are documented.
- [ ] The entire environment can be destroyed and redeployed consistently.
- [ ] A Makefile is created with atomic commands for Docker and Kubernetes management, including environment nuke commands.

## Out of Scope
- Cloud provider cluster provisioning (AWS/GCP/Azure).
- Remote Terraform state management.
