# Payflow — Payment Processing Platform

A microservices-based payment processing platform built with Go and Flutter, demonstrating distributed systems concepts, event-driven architecture, and idempotent payment processing.

---

## Tech Stack

**Backend:** Go · PostgreSQL · RabbitMQ · gRPC · Protocol Buffers · GORM

**Frontend:** Flutter · GetX

---

## Architecture

Payflow consists of four independent microservices communicating over gRPC (internally) and REST (externally).

```
Mobile/Browser
      │
      │ HTTP/JSON
      ▼
┌─────────────────────────────────────────┐
│  Identity Service  :8080                │
│  Merchant Service  :8082                │
│  Payment Service   :8081                │
│  Wallet HTTP       :8083                │
└─────────────────────────────────────────┘
      │
      │ gRPC
      ▼
┌─────────────────────────────────────────┐
│  Wallet Service    :50051               │
└─────────────────────────────────────────┘
      │
      │ RabbitMQ
      ▼
┌─────────────────────────────────────────┐
│  Settlement Worker   (payment.created)  │
│  Merchant Worker     (payment.captured) │
└─────────────────────────────────────────┘
```

### Services

| Service          | Port                       | Responsibility                                  |
| ---------------- | -------------------------- | ----------------------------------------------- |
| Identity Service | 8080                       | Auth, JWT, refresh tokens                       |
| Merchant Service | 8082                       | Merchant onboarding and profiles                |
| Payment Service  | 8081                       | Payment creation, idempotency, event publishing |
| Wallet Service   | 50051 (gRPC) / 8083 (HTTP) | Wallet balances and transaction history         |

---

## Payment Settlement Flow

```
POST /payments
    → payment row created (status: CREATED)
    → event published to payment.created queue

Settlement Worker (payment.created)
    → status: PROCESSING
    → debit user wallet (gRPC)
    → credit platform wallet (gRPC)
    → status: FUNDS_CAPTURED
    → event published to payment.captured queue

Merchant Worker (payment.captured)
    → debit platform wallet (gRPC)
    → credit merchant wallet (gRPC)
    → status: SETTLED
```

---

## Key Features

- **JWT authentication** with 15-minute access tokens and 30-day refresh tokens
- **Role-based access control** — user and merchant roles with middleware enforcement
- **Idempotent payments** — `Idempotency-Key` header with request hash validation prevents duplicate charges on retries
- **Async settlement** via RabbitMQ keeps payment creation fast and decoupled from wallet operations
- **Transaction-level idempotency** — wallet debit/credit operations are idempotent via reference field deduplication
- **Row-level locking** on wallet operations prevents race conditions during concurrent transactions
- **GORM AutoMigrate** — schema managed in code, no manual SQL migrations

---

## Prerequisites

- Go 1.21+
- PostgreSQL 15+
- RabbitMQ
- Flutter SDK

---

## Environment Variables

Each service has its own `.env` file. Below are the required variables per service.

**Identity Service**

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payflow_auth
DB_USERNAME=postgres
DB_PASSWORD=yourpassword
JwtKey=your-secret-key
WALLET_CLIENT=localhost:50051
```

**Merchant Service**

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payflow_auth
DB_USERNAME=postgres
DB_PASSWORD=yourpassword
JwtKey=your-secret-key
WALLET_CLIENT=localhost:50051
```

**Payment Service**

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payflow_auth
DB_USERNAME=postgres
DB_PASSWORD=yourpassword
JwtKey=your-secret-key
WALLET_CLIENT=localhost:50051
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
PLATFORM_USER_ID=your-platform-wallet-user-id
```

**Wallet Service**

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payflow_auth
DB_USERNAME=postgres
DB_PASSWORD=yourpassword
JwtKey=your-secret-key
WALLET_PORT=50051
```

> **Note:** All services share a single PostgreSQL database in this setup. The platform user wallet must be pre-created manually before processing payments.

---

## Running Locally

**1. Start infrastructure**

```bash
# PostgreSQL and RabbitMQ must be running
```

**2. Start services in order** (wallet service must be first — others depend on it via gRPC)

```bash
# Terminal 1
cd server/services/wallet-service && go run main.go

# Terminal 2
cd server/services/identity-service && go run main.go

# Terminal 3
cd server/services/merchant-service && go run main.go

# Terminal 4
cd server/services/payment-service && go run main.go
```

Or using Make:

```bash
cd server && make dev
```

**3. Run Flutter app**

```bash
cd client
flutter pub get
flutter run
```

---

## API Reference

### Identity Service — `localhost:8080`

| Method | Endpoint        | Auth | Description                     |
| ------ | --------------- | ---- | ------------------------------- |
| POST   | `/auth/signup`  | No   | Register user or merchant       |
| POST   | `/auth/signin`  | No   | Sign in, returns JWT            |
| POST   | `/auth/refresh` | No   | Refresh access token via cookie |

### Merchant Service — `localhost:8082`

| Method | Endpoint                      | Auth           | Description             |
| ------ | ----------------------------- | -------------- | ----------------------- |
| POST   | `/merchant/onboard`           | Yes (merchant) | Onboard merchant        |
| GET    | `/merchant/onboarding/status` | Yes (merchant) | Check onboarding status |
| GET    | `/merchant/list`              | Yes            | List active merchants   |

### Payment Service — `localhost:8081`

| Method | Endpoint                      | Auth | Description                                        |
| ------ | ----------------------------- | ---- | -------------------------------------------------- |
| POST   | `/payments`                   | Yes  | Create payment — requires `Idempotency-Key` header |
| GET    | `/payments/status?reference=` | Yes  | Poll payment status                                |
| GET    | `/payments/history`           | Yes  | Paginated payment history                          |

### Wallet Service — `localhost:8083`

| Method | Endpoint               | Auth | Description                   |
| ------ | ---------------------- | ---- | ----------------------------- |
| GET    | `/wallet/balance`      | Yes  | Get wallet balance            |
| GET    | `/wallet/transactions` | Yes  | Paginated transaction history |

---

## Database Schema

All services share one PostgreSQL instance with the following tables:

| Table              | Service  | Description           |
| ------------------ | -------- | --------------------- |
| `users`            | Identity | User accounts         |
| `refresh_tokens`   | Identity | Refresh token store   |
| `merchants`        | Merchant | Merchant profiles     |
| `payments`         | Payment  | Payment records       |
| `idempotency_keys` | Payment  | Idempotency key store |
| `wallets`          | Wallet   | Wallet balances       |
| `transactions`     | Wallet   | Transaction ledger    |

---

## Author

Built as a learning project to demonstrate microservices architecture, distributed systems concepts, and payment processing workflows.
