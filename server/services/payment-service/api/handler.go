package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k1ngalph0x/payflow/payment-service/config"
)

type PaymentHandler struct {
	DB *sql.DB
	Config *config.Config
}

type CreatePaymentRequest struct{
	MerchantID string `json:"merchant_id"`
	Amount float64 `json:"amount"`
}

func NewPaymentHandler(db *sql.DB, cfg *config.Config)*PaymentHandler{
	return &PaymentHandler{DB: db, Config: cfg}
}


func (h *PaymentHandler) CreatePayment(c *gin.Context){
	var req CreatePaymentRequest
	userId := c.GetString("user_id")
	ref := uuid.NewString()
	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request"})
		return 
	}

	insertQuery := `INSERT INTO payflow_payments (user_id, merchant_id, amount, status, reference)
	VALUES ($1, $2, $3, 'CREATED', $4)
	`	
	_, err = h.DB.Exec(insertQuery, userId, req.MerchantID, req.Amount, ref)

	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to create payment"})
	}

	c.JSON(http.StatusCreated, gin.H{
		"reference":ref,
		"status":"PENDING",
	})
}