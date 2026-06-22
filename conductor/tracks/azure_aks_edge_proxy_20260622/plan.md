# Implementation Plan: Migrate to Azure AKS and Build Edge Proxy

## Phase: Initial Azure CLI Setup & Ephemeral Infrastructure
- [x] Task: Write Terraform scripts for an AKS cluster and Azure VNet networking (Subnets, NSGs)
- [x] Task: Create Terraform outputs for Cost Management budgets and billing alerts
- [x] Task: Write `make` commands for automated environment teardown to enforce ephemeral lifecycle [50e04cd]
- [x] Task: Run Terraform to provision the Azure environment
- [ ] Task: Conductor - User Manual Verification 'Initial Azure CLI Setup & Ephemeral Infrastructure' (Protocol in workflow.md)

## Phase: Ingress & Cloudflare DNS Integration
- [ ] Task: Deploy an Ingress Controller (e.g., NGINX) to the AKS cluster
- [ ] Task: Configure Cloudflare DNS records to point `aleclabs.us` to the Azure Load Balancer Public IP
- [ ] Task: Create Ingress rules to route `aleclabs.us` (or subdomains) to the Grafana dashboards
- [ ] Task: Conductor - User Manual Verification 'Ingress & Cloudflare DNS Integration' (Protocol in workflow.md)

## Phase: Cost Analysis Documentation
- [ ] Task: Generate `azure-cost-optimization.md` outlining the hourly AKS costs, bandwidth estimations, and limits
- [ ] Task: Review the document against the initial constraints
- [ ] Task: Conductor - User Manual Verification 'Cost Analysis Documentation' (Protocol in workflow.md)
