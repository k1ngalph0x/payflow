package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/wallet-service/api"
	"github.com/k1ngalph0x/payflow/wallet-service/config"
	"github.com/k1ngalph0x/payflow/wallet-service/db"
	"github.com/k1ngalph0x/payflow/wallet-service/models"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	"google.golang.org/grpc"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	
	err = conn.AutoMigrate(&models.Wallet{}, &models.Transaction{})
	if err != nil {
		log.Fatalf("failed to run migrations: %w", err)
	}

	go func() {
		//listener, err := net.Listen("tcp", ":"+config.WALLET.Port)
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("gRPC listener failed: %v", err)
		}

		grpcServer := grpc.NewServer()
		walletpb.RegisterWalletServiceServer(grpcServer, api.NewWalletHandler(conn, config))
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	httpHandler := api.NewWalletHTTPHandler(conn)
	authMiddleware := middleware.NewAuthMiddleware(config.TOKEN.JwtKey)

	router := gin.Default()
	router.Use(authMiddleware.RequireAuth())

	router.GET("/wallet/balance", httpHandler.GetBalance)
	router.GET("/wallet/transactions", httpHandler.GetTransactions)

	err = router.Run(":8083")
	if err != nil{
		log.Fatalf("HTTP server failed: %v", err)
	}
}