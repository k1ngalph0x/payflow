package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	"github.com/k1ngalph0x/payflow/payment-service/internal/events"
	"github.com/k1ngalph0x/payflow/payment-service/models"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB           *gorm.DB
	Config       *config.Config
	WalletClient *walletclient.WalletClient
	Publisher    *events.Publisher
}

type CreatePaymentRequest struct {
	MerchantID string  `json:"merchant_id" binding:"required,uuid"`
	Amount     float64 `json:"amount"      binding:"required,gt=0"`
}

func hashRequest(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

func parsePagination(c *gin.Context) (limit, offset int) {
	limit, offset = 20, 0
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

func NewPaymentHandler(db *gorm.DB, config *config.Config, wc *walletclient.WalletClient, pub *events.Publisher) *PaymentHandler {
	return &PaymentHandler{DB: db, Config: config, WalletClient: wc, Publisher: pub}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID := c.GetString("user_id")

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header required"})
		return
	}

	var req CreatePaymentRequest
	err := c.ShouldBindJSON(&req); 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqHash, err := hashRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	var ikey models.IdempotencyKey
	result := h.DB.Where("idempotency_key = ? AND user_id = ?", idemKey, userID).First(&ikey)
	if result.Error == nil {
		if ikey.RequestHash != reqHash {
			c.JSON(http.StatusConflict, gin.H{"error": "Idempotency key reuse with different payload"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"reference": ikey.PaymentReference, "idempotent": true})
		return
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	var merchant struct {
		UserID string
	}

	//result = h.DB.Raw("SELECT user_id FROM payflow_merchants WHERE id = ?", req.MerchantID).Scan(&merchant)
	result = h.DB.Table("merchants").Select("user_id").Where("id = ?", req.MerchantID).Scan(&merchant)
	if result.Error != nil || merchant.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	ref := uuid.NewString()
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		payment := models.Payment{
			UserID:         userID,
			MerchantID:     req.MerchantID,
			MerchantUserID: merchant.UserID,
			Amount:         req.Amount,
			Status:         models.PaymentStatusCreated,
			Reference:      ref,
		}
		err = tx.Create(&payment).Error; 
		if err != nil {
			return err
		}

		ik := models.IdempotencyKey{
			IdempotencyKey:   idemKey,
			UserID:           userID,
			PaymentReference: ref,
			RequestHash:      reqHash,
		}
		return tx.Create(&ik).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	err = h.Publisher.Publish("payment.created", map[string]string{"reference": ref}); 
	if err != nil {
		h.DB.Model(&models.Payment{}).Where("reference = ?", ref).Update("status", models.PaymentStatusFailed)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"reference": ref, "status": models.PaymentStatusCreated})
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	ref := c.Query("reference")
	if ref == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reference is required"})
		return
	}

	var payment models.Payment
	result := h.DB.Where("reference = ?", ref).First(&payment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	limit, offset := parsePagination(c)

	var payments []models.Payment
	result := h.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payments": payments, "limit": limit, "offset": offset})
}


