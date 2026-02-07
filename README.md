# Drone Delivery Management Backend

A robust, production-grade backend for managing autonomous drone deliveries. Built with **Go**, **Gin**, and **PostgreSQL**.

## 🚀 Features

- **JWT Authentication**: Secure role-based access for Admins, Users, and Drones.
- **Drone Management**: Real-time location tracking and status updates.
- **Order Lifecycle**: ACID-compliant state transitions (Pending -> Reserved -> Delivered).
- **Dispatcher System**: Atomic job assignment to idle drones.
- **Observability**: Full tracing and metrics with **OpenTelemetry**, **Jaeger**, and **Prometheus**.

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL 16
- **Observability**: OpenTelemetry, Jaeger, Prometheus, Grafana
- **Infrastructure**: Docker Compose

## ⚡ Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Make (optional, but recommended)

### 1. Start Infrastructure
Start PostgreSQL, Jaeger, Prometheus, Grafana and App:
```bash
make docker-up
```

### 2. Run Database Migrations
Initialize the database schema:
```bash
make migrate-up

```

## 🧪 Verification & Testing

We include a script to simulate a complete delivery flow (Register -> Heartbeat -> Create Order -> Reserve -> Deliver):

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

### Authentication
- `POST /auth/token` - Login (Admin/User/Drone)

### Drones
- `POST /api/v1/drones` - Register new drone
- `POST /api/v1/drones/location` - Update heartbeat/location
- `POST /api/v1/drones/jobs/reserve` - Reserve pending order

### Orders
- `POST /api/v1/orders` - Create new order
- `POST /api/v1/orders/:id/status` - Update order status