package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	"github.com/k1ngalph0x/payflow/payment-service/internal/events"
)

type PaymentHandler struct {
	DB *sql.DB
	Config *config.Config
	//WalletClient *walletclient.WalletClient
	Publisher *events.Publisher
}

type CreatePaymentRequest struct{
	MerchantID string `json:"merchant_id"`
	Amount float64 `json:"amount"`
}

func NewPaymentHandler(db *sql.DB, cfg *config.Config, pub *events.Publisher)*PaymentHandler{
	return &PaymentHandler{DB: db, Config: cfg, Publisher: pub}
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
		return 
	}

	event := map[string]string{
		"reference":ref,
	}

	err = h.Publisher.Publish("payment.created", event)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to publish event"})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{
		"reference": ref,
		"status":"CREATED",
	})
}


func (h *PaymentHandler) failPayment(tx *sql.Tx, ref string)error {
	query := `UPDATE payflow_payments SET status = 'FAILED' WHERE reference = $1`
	_, err := tx.Exec(query, ref)
	return err
}