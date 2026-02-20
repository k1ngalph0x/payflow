.PHONY: help dev stop build clean deps test

help:
	@echo "Payflow - Available Commands:"
	@echo "  dev     - Start all services in order"
	@echo "  stop    - Stop all services"
	@echo "  build   - Build all services"
	@echo "  clean   - Clean build artifacts"
	@echo "  deps    - Download all dependencies"
	@echo "  test    - Run tests for all services"

dev:
	@echo "Starting Payflow services..."
	@echo "Starting Wallet Service..."
	@start cmd /k "cd services\wallet-service && go run main.go"
	@timeout /t 3 /nobreak > nul
	@echo "Starting Identity Service..."
	@start cmd /k "cd services\identity-service && go run main.go"
	@timeout /t 2 /nobreak > nul
	@echo "Starting Merchant Service..."
	@start cmd /k "cd services\merchant-service && go run main.go"
	@timeout /t 2 /nobreak > nul
	@echo "Starting Payment Service..."
	@start cmd /k "cd services\payment-service && go run main.go"
	@echo ""
	@echo "All services started!"
	@echo "  Identity Service:  http://localhost:8080"
	@echo "  Payment Service:   http://localhost:8081"
	@echo "  Merchant Service:  http://localhost:8082"
	@echo "  Wallet Service:    gRPC :50051 | HTTP :8083"

stop:
	@echo "Stopping all services..."
	@taskkill /F /IM go.exe /T 2>nul || echo No Go processes found
	@echo "All services stopped"

build:
	@echo "Building all services..."
	@cd services\identity-service && go build -o identity-service.exe .
	@cd services\merchant-service && go build -o merchant-service.exe .
	@cd services\payment-service && go build -o payment-service.exe .
	@cd services\wallet-service && go build -o wallet-service.exe .
	@echo "Build complete"

clean:
	@echo "Cleaning build artifacts..."
	@if exist services\identity-service\identity-service.exe del /Q services\identity-service\identity-service.exe
	@if exist services\merchant-service\merchant-service.exe del /Q services\merchant-service\merchant-service.exe
	@if exist services\payment-service\payment-service.exe  del /Q services\payment-service\payment-service.exe
	@if exist services\wallet-service\wallet-service.exe    del /Q services\wallet-service\wallet-service.exe
	@echo "Cleanup complete"

deps:
	@echo "Downloading dependencies..."
	@cd services\identity-service && go mod download
	@cd services\merchant-service && go mod download
	@cd services\payment-service  && go mod download
	@cd services\wallet-service   && go mod download
	@cd client                    && go mod download
	@echo "Dependencies downloaded"

test:
	@echo "Running tests..."
	@cd services\identity-service && go test ./... -v
	@cd services\merchant-service && go test ./... -v
	@cd services\payment-service  && go test ./... -v
	@cd services\wallet-service   && go test ./... -v
	@echo "Tests complete"