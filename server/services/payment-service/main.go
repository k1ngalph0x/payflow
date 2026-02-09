package main

import (
	"log"

	"github.com/gin-gonic/gin"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/payment-service/api"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	"github.com/k1ngalph0x/payflow/payment-service/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: ", err)
	}

	//Connect to db
	conn, err := db.ConnectDB()

	if err != nil {
		log.Fatalf("Error connecting to database: ", err)
	}

	defer conn.Close()

	walletClient, err := walletclient.NewWalletClient("localhost:50051")
	if err != nil{
		log.Fatal(err)
	}

	handler := api.NewPaymentHandler(conn, cfg, walletClient)
	authMiddleware := middleware.NewAuthMiddleware(cfg.TOKEN.JwtKey)

	router := gin.Default()
	router.Use(gin.Logger())

	router.Use(authMiddleware.RequireAuth())

	router.POST("/payments", handler.CreatePayment)
	router.POST("/payments/:reference/settle", handler.SettlePayment) 

	router.Run(":8081")
}