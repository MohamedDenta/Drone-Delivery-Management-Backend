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

## ⚡ Getting Started

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Go 1.22+](https://go.dev/dl/)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI

### Quick Start
Follow these steps to get the system up and running:

1. **Start Infrastructure**: Start PostgreSQL, Redis, RabbitMQ, and Observability tools:
   ```bash
   make docker-up
   ```
2. **Database Migrations**: Initialize the database schema:
   ```bash
   make migrate-up
   ```
3. **Run Application**: Start the backend server:
   ```bash
   make run
   ```

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
- `GET /api/v1/drones` - List all drones (Admin)
- `POST /api/v1/drones` - Register drone
- `POST /api/v1/drones/location` - Update location & heartbeat (REST fallback)
- `PATCH /api/v1/drones/:id/status` - Manually update drone status (e.g., BROKEN/IDLE)
- `POST /api/v1/drones/jobs/reserve` - Manually reserve the next pending order
- `GET /api/v1/orders` - List all orders (Admin)
- `POST /api/v1/orders` - Create order (Asynchronous via RabbitMQ)
- `GET /api/v1/orders/:id` - Fetch order details (Status, Location, ETA)
- `PATCH /api/v1/orders/:id` - Update order destination (Only if PENDING)
- `POST /api/v1/orders/:id/status` - Manually update order state
- `DELETE /api/v1/orders/:id` - Withdraw/Cancel order (Only if not yet picked up)

### gRPC API (Streaming)
- `rpc ReportLocation(stream LocationRequest) returns (stream LocationResponse)`
  - Used by drones for high-frequency location updates.
  - Updates are cached in Redis, persisted to Postgres, and refresh the drone's **Heartbeat** (30s TTL).

#### Testing with `grpcurl`
We've enabled gRPC Reflection for easy testing:
1. **Install grpcurl**: `make install-tools`
2. **List services**: `grpcurl -plaintext localhost:50051 list`
3. **Test Stream**:
   ```bash
   grpcurl -plaintext -d '{"drone_id": "drone-001", "latitude": 30.0, "longitude": 31.0}' \
     localhost:50051 drone.DroneService/ReportLocation
   ```

## ⚙️ Background Workers
The system runs background processes for automation and reliability:
- **Order Dispatcher**: Consumes `order.created` events from RabbitMQ and assigns them to idle drones using atomic SQL locks.
- **Heartbeat Monitor**: Periodically scans Redis for expired drone heartbeats (drones missing for >30s) and marks them as `OFFLINE`, triggering immediate order recovery.