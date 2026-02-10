package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/wallet-service/api"
	"github.com/k1ngalph0x/payflow/wallet-service/config"
	"github.com/k1ngalph0x/payflow/wallet-service/db"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.LoadConfig()


	if err != nil {
		log.Fatalf("Error loading config: ", err)
	}

	conn, err := db.ConnectDB()

	if err != nil {
		log.Fatalf("Error connecting to database: ", err)
	}

	defer conn.Close()	

	walletHandler := api.NewWalletHandler(conn, cfg)

	go func() {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatal(err)
		}

		grpcServer := grpc.NewServer()
		walletpb.RegisterWalletServiceServer(grpcServer, walletHandler)

		fmt.Println("Wallet-Service running on :50051")
		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatal(err)
		}
	}()

	httpHandler := api.NewHTTPHandler(walletHandler)
	authMiddleware := middleware.NewAuthMiddleware(cfg.TOKEN.JwtKey)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(authMiddleware.RequireAuth())

	router.GET("/wallet/balance", httpHandler.GetBalance)
	router.GET("/wallet/transactions", httpHandler.GetTransactions)

	router.Run(":50052")
}