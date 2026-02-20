package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/k1ngalph0x/payflow/wallet-service/models"
	"gorm.io/gorm"
)

type WalletHTTPHandler struct {
	DB *gorm.DB
}


func NewWalletHTTPHandler(db *gorm.DB) *WalletHTTPHandler {
	return &WalletHTTPHandler{DB: db}
}

func (h *WalletHTTPHandler) GetBalance(c *gin.Context) {
	userID := c.GetString("user_id")

	var wallet models.Wallet
	result :=  h.DB.Where("user_id = ?", userID).First(&wallet)
	if result.Error != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
			return 
	}else if errors.Is(result.Error, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return 
	}

	c.JSON(http.StatusOK, gin.H{"balance": wallet.Balance})
}

func (h *WalletHTTPHandler) GetTransactions(c *gin.Context) {
	userID := c.GetString("user_id")
	limit, offset := parsePagination(c)

	var transactions []models.Transaction
	result := h.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	})
}

func parsePagination(c *gin.Context) (limit, offset int) {
	limit = 20
	offset = 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}
	return
}