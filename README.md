# Drone Delivery Management Backend

A robust, production-grade backend for managing autonomous drone deliveries. Built with **Go**, **Gin**, and **PostgreSQL**.

## 🚀 Features

- **gRPC Streaming**: High-performance real-time location updates via bidirectional streams.
- **Redis Caching**: Sub-millisecond drone location lookups and reduced DB load.
- **RabbitMQ Async Dispatching**: Decoupled order processing and automated job assignment.
- **Offline Drone Detection**: Automated heartbeat monitoring via Redis TTL and background recovery.
- **Atomic Order Reservation**: Race-condition-free job assignment using Postgres `FOR UPDATE SKIP LOCKED`.
- **Observability**: Full tracing and metrics with **OpenTelemetry**, **Jaeger**, and **Prometheus**.

## 🛠️ Tech Stack

- **Language**: Go 1.22+
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **gRPC**: [Protobuf](https://github.com/protocolbuffers/protobuf) for real-time streaming
- **Database**: PostgreSQL 16
- **Caching**: Redis 7
- **Message Broker**: RabbitMQ 3
- **Observability**: OpenTelemetry, Jaeger, Prometheus, Grafana
- **Infrastructure**: Docker Compose

## ⚙️ Configuration
The application is configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | Postgres DSN | `postgres://user:password@localhost:5433/drone_delivery?sslmode=disable` |
| `REDIS_URL` | Redis endpoint | `localhost:6379` |
| `RABBITMQ_URL` | RabbitMQ AMQP URL | `amqp://user:password@localhost:5672/` |
| `PORT` | HTTP Server Port | `8081` |
| `OTEL_COLLECTOR_URL` | OTel Collector Endpoint | `localhost:4317` |

## 🧪 Verification & Testing

We include a script to simulate a complete delivery flow:
```bash
chmod +x scripts/test_flow.sh
./scripts/test_flow.sh
```

## 📊 Observability
Access the dashboards to monitor the system:

| Service | URL | Credentials |
|---------|-----|-------------|
| **Jaeger** (Traces) | http://localhost:16686 | N/A |
| **Prometheus** (Metrics) | http://localhost:9090 | N/A |
| **Grafana** (Dashboards) | http://localhost:3000 | admin / admin |

## 📝 API Reference

### HTTP API (REST)
- `POST /auth/token` - Login (Admin/User/Drone)
- `POST /api/v1/drones` - Register drone
- `POST /api/v1/orders` - Create order (Asynchronous via RabbitMQ)

### gRPC API (Streaming)
- `rpc ReportLocation(stream LocationRequest) returns (stream LocationResponse)`
  - Used by drones for high-frequency location updates.
  - Updates are cached in Redis and persisted to Postgres.