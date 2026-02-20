package main

import (
	"log"

	"github.com/gin-gonic/gin"

	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/merchant-service/api"
	"github.com/k1ngalph0x/payflow/merchant-service/config"
	"github.com/k1ngalph0x/payflow/merchant-service/db"
	"github.com/k1ngalph0x/payflow/merchant-service/models"
)

func main() {

	config, err := config.LoadConfig()
	if err != nil{
		log.Fatalf("Error loading config: %v", err)
	}
	
	conn, err := db.ConnectDB()
	if err != nil{
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = conn.AutoMigrate(&models.Merchant{}); 
	if err != nil {
		log.Fatalf("failed to run migrations: %w", err)
	}

	walletClient, err := walletclient.NewWalletClient(config.PLATFORM.WalletClient)
	if err != nil{
		log.Fatalf("Error creating wallet client: %v", err)
	}
	handler := api.NewMerchantHandler(conn, config, walletClient)
	authMiddleware := middleware.NewAuthMiddleware(config.TOKEN.JwtKey)

	router := gin.Default()
	router.Use(gin.Logger())

	router.Use(authMiddleware.RequireAuth())
	router.GET("/merchant/list", handler.GetMerchants) 
	router.POST("/merchant/onboard", handler.Onboard)
	router.GET("/merchant/onboarding/status", handler.OnboardingStatus)

	err = router.Run(":8082")
	if err != nil{
		log.Fatalf("Server failed: %v", err)
	}
}

