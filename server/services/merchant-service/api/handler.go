package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/payflow/merchant-service/config"
	"github.com/k1ngalph0x/payflow/merchant-service/models"
	"gorm.io/gorm"

	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
)

type MerchantHandler struct {
	DB   *gorm.DB
	Config *config.Config
	WalletClient *walletclient.WalletClient
}

func NewMerchantHandler(db *gorm.DB, config *config.Config, walletclient *walletclient.WalletClient) *MerchantHandler {
	return &MerchantHandler{DB: db, Config: config, WalletClient: walletclient}
}

type OnboardRequest struct {
	BusinessName string `json:"business_name" binding:"required,min=2,max=100"`
}


func(h *MerchantHandler) Onboard(c *gin.Context){
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if role != "merchant"{
		c.JSON(http.StatusForbidden, gin.H{"error": "Only merchants can onboard"})
		return 
	}

	var req OnboardRequest
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return 
	}

	var existing models.Merchant
	result := h.DB.Where("user_id = ?", userID).First(&existing)
	if result.Error == nil{
		c.JSON(http.StatusConflict, gin.H{"error": "Merchant already onboarded"})
		return
	}else if !errors.Is(result.Error, gorm.ErrRecordNotFound){
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	merchant := models.Merchant{
		UserID:       userID,
		BusinessName: req.BusinessName,
		Status:       models.MerchantStatusActive,
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		  err := tx.Create(&merchant).Error
			if err != nil {
				return err
			}
			err = h.WalletClient.CreateWallet(merchant.ID)
			if err != nil{
				return err
			}
			return nil
	})

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete onboarding"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"merchant_id": merchant.ID,
		"status":      merchant.Status,
	})
}

func(h *MerchantHandler) GetMerchants(c *gin.Context){

	var merchants []models.Merchant
	result := h.DB.Where("status = ?", models.MerchantStatusActive).Order("business_name ASC").Find(&merchants)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch merchants"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"merchants": merchants})
}
func (h *MerchantHandler) OnboardingStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")
		if role != "merchant" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only merchants are allowed"})
		return
	}

	var merchant models.Merchant
	result := h.DB.Where("user_id = ?", userID).First(&merchant)
	if result.Error != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"onboarded": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"onboarded":     merchant.Status == models.MerchantStatusActive,
		"merchant_id":   merchant.ID,
		"status":        merchant.Status,
		"business_name": merchant.BusinessName,
	})
}