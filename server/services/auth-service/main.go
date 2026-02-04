package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/payflow/auth-service/api"
	"github.com/k1ngalph0x/payflow/auth-service/config"
	"github.com/k1ngalph0x/payflow/auth-service/db"
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

	handler := api.NewHandler(conn, cfg)

	//1. Do move to the routes file
	//Routes config
	router := gin.Default()
	router.Use(gin.Logger())

	auth := router.Group("/auth")
	auth.POST("/signup", handler.SignUp)
	auth.POST("/signin", handler.SignIn)
	//////////


	router.Run(":8080")

	fmt.Println("Running auth-service")
}