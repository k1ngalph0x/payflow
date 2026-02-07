package main

import (
	"log"

	"github.com/gin-gonic/gin"
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

	handler := api.NewPaymentHandler(conn, cfg)

	router := gin.Default()
	router.Use(gin.Logger())

	router.POST("/payments", handler.CreatePayment)
}