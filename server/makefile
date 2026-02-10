.PHONY: help dev stop clean build deps

help: 
	@echo "Payflow - Available Commands:"
	@echo "  dev     - Start all services"
	@echo "  stop    - Stop all services"
	@echo "  build   - Build all services"
	@echo "  clean   - Clean build artifacts"
	@echo "  deps    - Download dependencies"

dev:
	@echo "Starting all Payflow services..."
	@echo "Starting Wallet Service..."
	@start cmd /k "cd services/wallet-service && go run main.go"
	@timeout /t 2 /nobreak > nul
	@echo "Starting Identity Service..."
	@start cmd /k "cd services/identity-service && go run main.go"
	@timeout /t 1 /nobreak > nul
	@echo "Starting Merchant Service..."
	@start cmd /k "cd services/merchant-service && go run main.go"
	@timeout /t 1 /nobreak > nul
	@echo "Starting Payment Service..."
	@start cmd /k "cd services/payment-service && go run main.go"
	@echo ""
	@echo "All services started!"
	@echo "Identity Service:  http://localhost:8080"
	@echo "Payment Service:   http://localhost:8081"
	@echo "Merchant Service:  http://localhost:8082"
	@echo "Wallet Service:    gRPC :50051, HTTP :50052"

stop:
	@echo "Stopping all services..."
	@taskkill /F /IM go.exe /T 2>nul || echo No Go processes found
	@echo "All services stopped"

build:
	@echo "Building all services..."
	@cd services/identity-service && go build -o identity-service.exe
	@cd services/merchant-service && go build -o merchant-service.exe
	@cd services/payment-service && go build -o payment-service.exe
	@cd services/wallet-service && go build -o wallet-service.exe
	@echo "Build complete"

clean:
	@echo "Cleaning build artifacts..."
	@del /Q services\identity-service\identity-service.exe 2>nul || echo ""
	@del /Q services\merchant-service\merchant-service.exe 2>nul || echo ""
	@del /Q services\payment-service\payment-service.exe 2>nul || echo ""
	@del /Q services\wallet-service\wallet-service.exe 2>nul || echo ""
	@echo "Cleanup complete"

deps:
	@echo "Downloading dependencies..."
	@cd services/identity-service && go mod download
	@cd services/merchant-service && go mod download
	@cd services/payment-service && go mod download
	@cd services/wallet-service && go mod download
	@cd client && go mod download
	@echo "Dependencies downloaded"