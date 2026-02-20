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
	"github.com/k1ngalph0x/payflow/payment-service/models"
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

	err = conn.AutoMigrate(&models.Payment{}, &models.IdempotencyKey{})
	if err != nil {
		log.Fatalf("failed to run migrations: %w", err)
	}


	walletClient, err := walletclient.NewWalletClient(config.PLATFORM.WalletClient)
	if err != nil {
		log.Fatalf("Error creating wallet client: %v", err)
	}
	
	publisher, err := events.NewPublisher(config.PLATFORM.RabbitMQURL)
	if err != nil{
		log.Fatalf("Error creating publisher: %v", err)
	}
	defer publisher.Conn.Close()
	defer publisher.Channel.Close()

	go worker.StartSettlementWorker(
		conn, 
		config, 
		walletClient, 
		config.PLATFORM.RabbitMQURL,
	)

	go worker.StartMerchantSettlementWorker(
		conn, 
		config, 
		walletClient, 
		config.PLATFORM.RabbitMQURL,
	)

	handler := api.NewPaymentHandler(conn, config, walletClient, publisher)
	authMiddleware := middleware.NewAuthMiddleware(config.TOKEN.JwtKey)

	router := gin.Default()
	router.Use(gin.Logger())

	router.Use(authMiddleware.RequireAuth())
	router.POST("/payments", handler.CreatePayment)
	router.GET("/payments/status", handler.GetPaymentStatus)    
	router.GET("/payments/history", handler.GetPaymentHistory) 
	//router.POST("/payments/:reference/settle", handler.SettlePayment) 
	err = router.Run(":8081")
	if err != nil{
		log.Fatalf("Server failed: %v", err)
	}
}