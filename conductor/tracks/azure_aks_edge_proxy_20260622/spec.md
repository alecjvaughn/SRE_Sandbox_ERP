# Specification: Migrate to Azure AKS and Build Edge Proxy

## Overview
This track pivots the infrastructure from DigitalOcean to Azure Kubernetes Service (AKS) to utilize enterprise-grade networking primitives (VNets, NSGs). It introduces a Custom Go Edge Proxy to route Kafka traffic efficiently. **Crucially, this track is designed for an environment with no upfront credits, prioritizing an ultra-low-cost, ephemeral architecture.**

## Functional Requirements
- Provision an AKS cluster using Terraform configured for extreme cost efficiency.
- **Cost Monitoring & Documentation**: 
  - Generate an "Azure Hourly Cost Optimization" architecture document detailing exact hourly pricing, spot instance strategies, and bandwidth estimations.
  - Implement Azure Cost Management budgets and billing alerts via Terraform/CLI to strictly cap spending and advise on limit-setting.
- Configure Azure VNets, Subnets, and Network Security Groups (NSGs) for secure, isolated routing.
- Implement a custom Layer 4 TCP Edge Proxy in Go using the `net` package to intelligently route Kafka traffic based on real-time AZ health and cost metrics.
- Utilize Chaos Mesh `NetworkChaos` to inject latency and packet loss to simulate cross-AZ link drops, triggering the proxy's failover logic.
- Emit network I/O and simulated egress cost metrics to Prometheus.

## Non-Functional Requirements
- **Zero-to-Minimal Cost Enforcement**: 
  - The AKS node pool MUST use ultra-low-cost instances (e.g., Azure Spot instances or `Standard_B2s` on ephemeral deployments).
  - The infrastructure is strictly **Ephemeral**; the track must include highly reliable teardown automation (`terraform destroy`) to ensure the cluster only runs for a few hours during active development/testing.
- The Go Edge Proxy must gracefully handle connection drops and orchestrate quick reconnects to secondary AZs during a simulated failure.

## Out of Scope
- Modifying the core business logic of the Java inventory service or the mock ERP functionality.
