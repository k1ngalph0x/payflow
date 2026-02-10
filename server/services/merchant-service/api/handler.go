package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/payflow/merchant-service/config"

	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
)

type MerchantHandler struct {
	DB *sql.DB
	Config *config.Config
	WalletClient *walletclient.WalletClient
}

func NewMerchantHandler(db *sql.DB, cfg *config.Config, walletclient *walletclient.WalletClient) *MerchantHandler {
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

	err = h.WalletClient.CreateWallet(merchantId)
	if err != nil{
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to create merchant wallet"})
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

func(h *MerchantHandler) GetMerchants(c *gin.Context){
	query := `
		SELECT id, user_id, business_name, status, created_at
		FROM payflow_merchants
		WHERE status = 'ACTIVE'
		ORDER BY business_name ASC
	`
	
	rows, err := h.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch merchants"})
		return
	}
	defer rows.Close()

	var merchants []map[string]interface{}

	for rows.Next() {
		var id, userId, businessName, status string
		var createdAt time.Time

		err := rows.Scan(&id, &userId, &businessName, &status, &createdAt)
		if err != nil {
			continue
		}

		merchants = append(merchants, map[string]interface{}{
			"id":            id,
			"user_id":       userId,
			"business_name": businessName,
			"status":        status,
			"created_at":    createdAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, gin.H{"merchants": merchants})
}
func (h *MerchantHandler) OnboardingStatus(c *gin.Context) {
	userId := c.GetString("user_id")
	role := c.GetString("role")
	var merchantId string
	var status string
	var businessName string 

	if role != "merchant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only merchants are allowed"})
		return
	}

	query := `
		SELECT id, status, business_name
		FROM payflow_merchants
		WHERE user_id = $1
		LIMIT 1
	`

	err := h.DB.QueryRow(query, userId).Scan(&merchantId, &status, &businessName)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"onboarded": false,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"onboarded":  status == "ACTIVE",
		"merchant_id": merchantId,
		"status":     status,
		"business_name": businessName, 
	})
}