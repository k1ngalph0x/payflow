package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/api"
	"github.com/k1ngalph0x/payflow/identity-service/config"
	"github.com/k1ngalph0x/payflow/identity-service/db"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
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

	walletClient, err := walletclient.NewWalletClient("localhost:50051")
	if err != nil{
		log.Fatal(err)
	}
	handler := api.NewHandler(conn, cfg, walletClient)
	authMiddleware := middleware.NewAuthMiddleware(cfg.TOKEN.JwtKey)


	//Routes config
	router := gin.Default()
	router.Use(gin.Logger())

	auth := router.Group("/auth")
	auth.POST("/signup", handler.SignUp)
	auth.POST("/signin", handler.SignIn)
	///////////////////////////////////

	protected := router.Group("/api")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/profile", Profile)
	}

	/////////////////////////////////
	router.Run(":8080")

	fmt.Println("Running auth-service")
}

func Profile(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"Message":"profile page"})
}