package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	grpcclient "github.com/k1ngalph0x/payflow/wallet-service/grpc"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
)

type PaymentHandler struct {
	DB *sql.DB
	Config *config.Config
	WalletClient *grpcclient.WalletClient
}

type CreatePaymentRequest struct{
	MerchantID string `json:"merchant_id"`
	Amount float64 `json:"amount"`
}

func NewPaymentHandler(db *sql.DB, cfg *config.Config, walletClient *grpcclient.WalletClient)*PaymentHandler{
	return &PaymentHandler{DB: db, Config: cfg, WalletClient: walletClient}
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

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resp, err := h.WalletClient.Client.Debit(ctx, &walletpb.DebitRequest{
		UserId: userId,
		Amount: req.Amount,
		Reference: ref,
	})

	if err != nil{
		query := `UPDATE payflow_payments SET status = 'FAILED WHERE reference = $1`
		h.DB.Exec(query, ref)

		c.JSON(http.StatusInternalServerError, gin.H{"error":"wallet debit failed"})
		return 
	}

	if resp.Status != "SUCCESS"{
		updateQuery := `UPDATE payflow_payments SET status = "FAILED" WHERE reference=$1`
		h.DB.Exec(updateQuery, ref)

		c.JSON(http.StatusPaymentRequired, gin.H{"status":resp.Status})
		return 
	}

	updateQuery := `UPDATE payflow_payments SET status = 'SUCCESS' WHERE reference=$1`
	_, err = h.DB.Exec(
		updateQuery, ref,
	)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"payment update failed"})
		return 
	}

	c.JSON(http.StatusCreated, gin.H{
		"reference": ref,
		"status":"SUCCESS",
	})
}