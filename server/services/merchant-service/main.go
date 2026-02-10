package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/merchant-service/api"
	"github.com/k1ngalph0x/payflow/merchant-service/config"
	"github.com/k1ngalph0x/payflow/merchant-service/db"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil{
		log.Fatalf("Error loading config: ", err)
	}
	
	conn, err := db.ConnectDB()

	if err != nil{
		log.Fatalf("Error connecting to database: ", err)
	}

	defer conn.Close()

	walletClient, err := walletclient.NewWalletClient(cfg.PLATFORM.WalletClient)
	if err != nil{
		log.Fatal(err)
	}
	handler := api.NewMerchantHandler(conn, cfg, walletClient)
	authMiddleware := middleware.NewAuthMiddleware(cfg.TOKEN.JwtKey)



	router := gin.Default()
	router.Use(gin.Logger())

	router.Use(authMiddleware.RequireAuth())
	router.GET("/merchant/list", handler.GetMerchants) 
	router.POST("/merchant/onboard", handler.Onboard)
	router.GET("/merchant/onboarding/status", handler.OnboardingStatus)

	router.Run(":8082")

	fmt.Println("Running merchant-service")
}

