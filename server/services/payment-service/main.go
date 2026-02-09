package main

import (
	"log"

	"github.com/gin-gonic/gin"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/payment-service/api"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	"github.com/k1ngalph0x/payflow/payment-service/db"
	"github.com/k1ngalph0x/payflow/payment-service/internal/events"
	"github.com/k1ngalph0x/payflow/payment-service/internal/worker"
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
	
	publisher, err := events.NewPublisher(cfg.PLATFORM.RabbitMQURL)
	if err != nil{
		log.Fatal(err)
	}
	defer publisher.Conn.Close()
	defer publisher.Channel.Close()

	go worker.StartSettlementWorker(conn, cfg, walletClient, cfg.PLATFORM.RabbitMQURL)

	handler := api.NewPaymentHandler(conn, cfg,  publisher)
	authMiddleware := middleware.NewAuthMiddleware(cfg.TOKEN.JwtKey)

	router := gin.Default()
	router.Use(gin.Logger())

	router.Use(authMiddleware.RequireAuth())

	router.POST("/payments", handler.CreatePayment)
	//router.POST("/payments/:reference/settle", handler.SettlePayment) 

	router.Run(":8081")
}