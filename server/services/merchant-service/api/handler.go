package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/payflow/merchant-service/config"

	grpcclient "github.com/k1ngalph0x/payflow/wallet-service/grpc"
)

type MerchantHandler struct {
	DB *sql.DB
	Config *config.Config
	WalletClient *grpcclient.WalletClient
}

func NewMerchantHandler(db *sql.DB, cfg *config.Config, walletclient *grpcclient.WalletClient) *MerchantHandler {
	return &MerchantHandler{DB: db, Config: cfg, WalletClient: walletclient}
}
type OnboardRequest struct{
	BusinessName string `json:"business_name"`
}


func(h *MerchantHandler) Onboard(c *gin.Context){
	userId := c.GetString("user_id")
	role := c.GetString("role")
	var req OnboardRequest
	var merchantId string
	if role != "merchant"{
		c.JSON(http.StatusForbidden, gin.H{"error": "Only merchants can onboard"})
		return 
	}

	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return 
	}
	tx, err := h.DB.Begin()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request"})
		return 
	}

	defer tx.Rollback()

	query := `INSERT INTO payflow_merchants (user_id, business_name, status)
	VALUES ($1, $2, 'ACTIVE')
	RETURNING id
	`

	err = tx.QueryRow(query, userId, req.BusinessName).Scan(&merchantId)
	if err != nil{
		c.JSON(http.StatusConflict, gin.H{"error": "Merchant already exists"})
		return
	}

	err = tx.Commit()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete onboarding"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"merchant_id": merchantId,
		"status":"ACTIVE",
	})

}