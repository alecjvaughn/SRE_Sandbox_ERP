FROM golang:alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the proxy
RUN CGO_ENABLED=0 GOOS=linux go build -o /edge-proxy ./cmd/proxy

# Run step
FROM alpine:3.19
WORKDIR /
COPY --from=builder /edge-proxy /edge-proxy

EXPOSE 8080 9090

ENTRYPOINT ["/edge-proxy"]
