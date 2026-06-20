# Track Specification: Implement Go Order and Python Inventory Services with Kafka Event Bus

## 1. Overview
The goal of this track is to build the core transactional and event-driven foundation of the SRE Sandbox. This includes writing the Go `order-service` API, the Python `inventory-service` consumer, containerizing both services, and testing their end-to-end event loop locally using a single-node Kafka broker under Docker Compose.

## 2. Component Specifications

### 2.1 Go Order Service (`order-service`)
- **Language**: Go 1.22+
- **REST Endpoints**:
  - `POST /order`:
    - Request Body: JSON-encoded, containing:
      - `order_id` (string, required, UUID or sequence number)
      - `item` (string, required, e.g. "Widget")
      - `qty` (integer, required, > 0)
    - Responses:
      - `202 Accepted` on success: `{"status": "Order Transmitted", "order_id": "<order_id>"}`
      - `400 Bad Request` if validations fail: `{"error": "<error message>"}`
      - `500 Internal Server Error` if Kafka production fails.
  - `GET /healthz`:
    - Returns `200 OK` (plain text "OK")
  - `GET /metrics`:
    - Serves Prometheus client metrics (standard golang metrics + custom metric).
- **Custom Metrics**:
  - `erp_orders_total` (Counter): Counts total orders processed. Labels: `status` ("success" or "failed").
- **Kafka Integration**:
  - Emits JSON payloads to the `orders` topic on successful requests.
  - Payload format:
    ```json
    {
      "order_id": "string",
      "item": "string",
      "qty": 1,
      "timestamp": "ISO-8601 string"
    }
    ```
  - Driver: `github.com/twmb/franz-go`

### 2.2 Python Inventory Service (`inventory-service`)
- **Language**: Python 3.11+
- **Kafka Consumer**:
  - Subscribes to the `orders` topic.
  - Joins the consumer group `inventory-group`.
  - Deserializes JSON order messages.
  - Simulates a stock decrease (logs the item name, order ID, and mutated stock level).
- **Observability Ports**:
  - Health endpoint `GET /healthz` exposed on port `8081` returning `200 OK`.
  - Prometheus metrics server started on port `8000` using `prometheus_client`.
- **Custom Metrics**:
  - `erp_stock_updates_total` (Counter): Cumulative count of stock changes processed.

### 2.3 Kafka broker
- Local single-node Kafka running in Docker Compose on port `9092`.
- Topic `orders` created with 2 partitions.

## 3. Local Orchestration (Docker Compose)
A unified `docker-compose.yml` configures:
1. `zookeeper`: For Kafka coordination.
2. `kafka`: Kafka broker listening on localhost:9092.
3. `order-service`: Runs Go web service listening on port `8080`.
4. `inventory-service`: Runs Python consumer serving metrics on port `8000` and health on port `8081`.
