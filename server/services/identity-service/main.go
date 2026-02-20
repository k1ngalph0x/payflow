package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/identity-service/api"
	"github.com/k1ngalph0x/payflow/identity-service/config"
	"github.com/k1ngalph0x/payflow/identity-service/db"
	"github.com/k1ngalph0x/payflow/identity-service/middleware"
	"github.com/k1ngalph0x/payflow/identity-service/models"
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

	err = conn.AutoMigrate(
		&models.User{}, 
		&models.RefreshToken{}); 
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	walletClient, err := walletclient.NewWalletClient(config.WALLET.WalletClient)
	if err != nil{
		log.Fatalf("Error creating wallet client: %v", err)
	}

	handler := api.NewHandler(conn, config, walletClient)
	authMiddleware := middleware.NewAuthMiddleware(config.TOKEN.JwtKey)

	router := gin.Default()
	router.Use(gin.Logger())

	auth := router.Group("/auth")
	auth.POST("/signup", handler.SignUp)	
	auth.POST("/signin", handler.SignIn)
	auth.POST("/refresh", handler.Refresh)

	protected := router.Group("/api")
	protected.Use(authMiddleware.RequireAuth())
	{
		protected.GET("/profile", Profile)
	}

	err = router.Run(":8080")

	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func Profile(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"Message":"profile page"})
}