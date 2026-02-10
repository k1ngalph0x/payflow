# Payflow - Payment Processing Platform

A microservices-based payment processing platform built with Go, demonstrating distributed systems concepts and event-driven architecture.

## Overview

Payflow is a payment platform that enables users to send and receive payments while merchants can accept payments from customers. The system is built using a microservices architecture with asynchronous payment settlement and idempotent operations.

## Tech Stack

**Backend:**

- Go 1.21
- PostgreSQL
- RabbitMQ
- gRPC
- Protocol Buffers

**Frontend:**

- Flutter
- GetX (State Management)

**Architecture:**

- Microservices Architecture
- Event-Driven Design
- RESTful APIs
- Inter-service Communication (gRPC)

## System Architecture

The platform consists of four independent microservices:

- **Identity Service** (Port 8080): Handles user authentication, JWT token generation, and refresh token management
- **Merchant Service** (Port 8082): Manages merchant onboarding and profiles
- **Payment Service** (Port 8081): Processes payment creation with idempotency support and publishes settlement events
- **Wallet Service** (gRPC :50051, HTTP :50052): Manages user/merchant wallet balances and transaction history

## Key Features

- User and merchant authentication with JWT and refresh tokens
- Idempotent payment processing with request hash validation
- Asynchronous payment settlement using RabbitMQ
- Wallet management with transaction-level idempotency
- Real-time balance and transaction history
- Role-based access control

## Project Structure

```
.
├── server/
│   ├── services/
│   │   ├── identity-service/
│   │   ├── merchant-service/
│   │   ├── payment-service/
│   │   └── wallet-service/
│   ├── client/
│   └── Makefile
└── client/ (Flutter app)
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- RabbitMQ
- Flutter SDK (for mobile app)
- Make (optional, for running commands)

## Environment Setup

Each service requires a `.env` file with the following variables:

```env
DB_URL="host=localhost port=5432 dbname=postgres user=postgres password=yourpassword sslmode=disable"
JwtKey="your-secret-key"
WALLET_SERVICE_ADDR="localhost:50051"
RABBITMQ_URL="amqp://guest:guest@localhost:5672/"
PLATFORM_USER_ID="your-platform-user-id"
```

## Running the Application

### Start Infrastructure

```bash
# Start PostgreSQL
# Start RabbitMQ
```

### Start Backend Services

**Using Makefile:**

```bash
cd server
make dev
```

**Manual start:**

```bash
# Terminal 1 - Wallet Service
cd server/services/wallet-service
go run main.go

# Terminal 2 - Identity Service
cd server/services/identity-service
go run main.go

# Terminal 3 - Merchant Service
cd server/services/merchant-service
go run main.go

# Terminal 4 - Payment Service
cd server/services/payment-service
go run main.go
```

### Run Flutter App

```bash
cd client
flutter pub get
flutter run
```

## API Endpoints

### Identity Service (8080)

- POST `/auth/signup` - User/Merchant registration
- POST `/auth/signin` - User login
- POST `/auth/refresh` - Refresh access token

### Merchant Service (8082)

- POST `/merchant/onboard` - Merchant onboarding
- GET `/merchant/onboarding/status` - Check onboarding status
- GET `/merchant/list` - List active merchants

### Payment Service (8081)

- POST `/payments` - Create payment (requires Idempotency-Key header)
- GET `/payments/status?reference={ref}` - Check payment status
- GET `/payments/history` - Payment history

### Wallet Service (50052)

- GET `/wallet/balance` - Get wallet balance
- GET `/wallet/transactions` - Transaction history

## Implementation Highlights

### Idempotency

Payment creation uses idempotency keys with request hash validation to prevent duplicate charges on network retries.

### Asynchronous Settlement

Payment settlement is handled asynchronously through RabbitMQ, ensuring the payment creation API remains fast and responsive.

### Transaction Safety

Wallet operations use database transactions with row-level locking to prevent race conditions during concurrent debits/credits.

### Reference-Based Idempotency

Wallet debit/credit operations check for existing transactions using the reference field, ensuring idempotent behavior at the transaction level.

## Database Schema

Each service maintains its own database with the following key tables:

- **Identity Service**: `payflow_auth`, `payflow_refresh_tokens`
- **Merchant Service**: `payflow_merchants`
- **Payment Service**: `payflow_payments`, `payflow_idempotency_keys`
- **Wallet Service**: `payflow_wallets`, `payflow_wallet_transactions`

## Development

### Build All Services

```bash
make build
```

### Download Dependencies

```bash
make deps
```

### Stop Services

```bash
make stop
```

### Clean Artifacts

```bash
make clean
```

## Testing

Run tests for all services:

```bash
make test
```

## License

MIT

## Author

Built as a learning project to demonstrate microservices architecture, distributed systems concepts, and payment processing workflows.
