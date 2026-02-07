package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/merchant-service/api"
	"github.com/k1ngalph0x/payflow/merchant-service/config"
	"github.com/k1ngalph0x/payflow/merchant-service/db"
	grpcclient "github.com/k1ngalph0x/payflow/wallet-service/grpc"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil{
		log.Fatalf("Error loading config: ", err)
	}

	//Connect to db
	conn, err := db.ConnectDB()

	if err != nil{
		log.Fatalf("Error connecting to database: ", err)
	}

	defer conn.Close()

	walletClient, err := grpcclient.NewWalletClient("localhost:50051")
	if err != nil{
		log.Fatal(err)
	}
	handler := api.NewMerchantHandler(conn, cfg, walletClient)
	authMiddleware := middleware.NewAuthMiddleware(cfg.TOKEN.JwtKey)


	//Routes config
	router := gin.Default()
	router.Use(gin.Logger())

	router.Use(authMiddleware.RequireAuth())

	router.POST("/merchant/onboard", handler.Onboard)

	router.Run(":8082")

	fmt.Println("Running merchant-service")
}

