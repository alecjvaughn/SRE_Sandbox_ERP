# Implementation Plan: Migrate to Azure AKS and Build Edge Proxy

## Phase: Initial Azure CLI Setup & Ephemeral Infrastructure
- [ ] Task: Write Terraform scripts for an AKS cluster (`Standard_B2s`) and Azure VNet networking (Subnets, NSGs)
- [ ] Task: Create Terraform outputs for Cost Management budgets and billing alerts
- [ ] Task: Write `make` commands for automated environment teardown to enforce ephemeral lifecycle
- [ ] Task: Run Terraform to provision the Azure environment
- [ ] Task: Conductor - User Manual Verification 'Initial Azure CLI Setup & Ephemeral Infrastructure' (Protocol in workflow.md)

## Phase: Ingress & Cloudflare DNS Integration
- [ ] Task: Deploy an Ingress Controller (e.g., NGINX) to the AKS cluster
- [ ] Task: Configure Cloudflare DNS records to point `aleclabs.us` to the Azure Load Balancer Public IP
- [ ] Task: Create Ingress rules to route `aleclabs.us` (or subdomains) to the Grafana dashboards
- [ ] Task: Conductor - User Manual Verification 'Ingress & Cloudflare DNS Integration' (Protocol in workflow.md)

## Phase: Custom Go Edge Proxy Development
- [ ] Task: Create new directory `edge-proxy` and initialize Go module
- [ ] Task: Implement raw TCP listener and forwarder logic using Go `net` package
- [ ] Task: Add basic AZ health checking and simulated cost-aware routing metrics to the proxy
- [ ] Task: Expose Prometheus metrics endpoint (`/metrics`) for Network I/O and Retransmission rates
- [ ] Task: Conductor - User Manual Verification 'Custom Go Edge Proxy Development' (Protocol in workflow.md)

## Phase: Deployment & Chaos Engineering Integration
- [ ] Task: Dockerize the Go Edge Proxy and create Kubernetes deployment manifests
- [ ] Task: Deploy the proxy, order-service, inventory-service, and Kafka to AKS
- [ ] Task: Write `NetworkChaos` CRDs to inject latency and packet loss between proxy and Kafka
- [ ] Task: Execute chaos experiments and monitor failover/reconnection behavior
- [ ] Task: Conductor - User Manual Verification 'Deployment & Chaos Engineering Integration' (Protocol in workflow.md)

## Phase: Cost Analysis Documentation
- [ ] Task: Generate `azure-cost-optimization.md` outlining the hourly AKS costs, bandwidth estimations, and limits
- [ ] Task: Review the document against the initial constraints
- [ ] Task: Conductor - User Manual Verification 'Cost Analysis Documentation' (Protocol in workflow.md)
