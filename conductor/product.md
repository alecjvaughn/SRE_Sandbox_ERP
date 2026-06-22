# Initial Concept

Cost-Aware Traffic Plane & Egress Optimizer for Kafka

# Product Guide: Cost-Aware Traffic Plane & Egress Optimizer

## Vision
This project is an advanced cloud-native data plane simulation targeting the optimization of ingress and egress traffic for Confluent Kafka workloads. It pivots away from standard application development to focus entirely on network routing, failover mechanisms, and bandwidth cost efficiency. Built on Microsoft Azure to utilize enterprise-grade networking primitives (VNets, NSGs, Private Links), the project provides a reproducible sandbox to prove how intelligent traffic shaping can maintain 99.9% uptime while actively minimizing cloud networking costs across Availability Zones.

## Core Features
1. **Custom Go Edge Proxy**: A lightweight Layer 4 TCP proxy built in Go that sits between the internet and the Kafka cluster. It inspects incoming traffic and intelligently routes it to the nearest or cheapest Kafka broker based on real-time AZ health.
2. **Simulated Hybrid Network Failovers**: Orchestrating critical network failovers during simulated hybrid-cloud link drops. The Go proxy dynamically re-routes Kafka consumer traffic to secondary pathways without dropping TCP connections.
3. **Network-Targeted Chaos Engineering**: Utilizing Chaos Mesh `NetworkChaos` CRDs to inject packet loss, 500ms latency delays, and bandwidth throttling between the transaction engine and Kafka brokers, simulating severe network degradation.
4. **Cost-Centric Observability**: Extending Prometheus and Grafana to track not just throughput, but explicit Network I/O (Bytes In/Out), TCP Retransmission rates, and real-time simulated Egress Costs ($) across Availability Zones.
5. **Azure Networking Primitives**: Leveraging Azure Virtual Networks (VNets), User Defined Routes (UDRs), Network Security Groups (NSGs), and Azure Private Endpoints to create a highly secure, isolated traffic fabric simulating a true Confluent Cloud deployment.

## Target Audience
- **Interviewers / Recruiters (IBM Confluent Cloud Traffic Team)**: Serving as a highly targeted, undeniably relevant senior-level architectural portfolio piece.
- **SRE & Cloud Network Engineers**: To test complex traffic routing, private link isolation, and cost-optimization practices under simulated stress.
