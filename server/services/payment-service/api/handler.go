package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
)

type PaymentHandler struct {
	DB *sql.DB
	Config *config.Config
	WalletClient *walletclient.WalletClient
}

type CreatePaymentRequest struct{
	MerchantID string `json:"merchant_id"`
	Amount float64 `json:"amount"`
}

func NewPaymentHandler(db *sql.DB, cfg *config.Config, walletClient *walletclient.WalletClient)*PaymentHandler{
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
		return 
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	// defer cancel()

	// resp, err := h.WalletClient.Client.Debit(ctx, &walletpb.DebitRequest{
	// 	UserId: userId,
	// 	Amount: req.Amount,
	// 	Reference: ref,
	// })

	// if err != nil{
	// 	query := `UPDATE payflow_payments SET status = 'FAILED' WHERE reference = $1`
	// 	h.DB.Exec(query, ref)

	// 	c.JSON(http.StatusInternalServerError, gin.H{"error":"wallet debit failed"})
	// 	return 
	// }

	// if resp.Status != "SUCCESS"{
	// 	updateQuery := `UPDATE payflow_payments SET status = 'FAILED' WHERE reference=$1`
	// 	h.DB.Exec(updateQuery, ref)

	// 	c.JSON(http.StatusPaymentRequired, gin.H{"status":resp.Status})
	// 	return 
	// }

	// updateQuery := `UPDATE payflow_payments SET status = 'SUCCESS' WHERE reference=$1`
	// _, err = h.DB.Exec(
	// 	updateQuery, ref,
	// )

	// if err != nil{
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error":"payment update failed"})
	// 	return 
	// }

	c.JSON(http.StatusCreated, gin.H{
		"reference": ref,
		"status":"CREATED",
	})
}

func (h *PaymentHandler) SettlePayment (c *gin.Context){
	
	ref := c.Param("reference")

	var payment struct{
		UserId string
		MerchantId string
		Amount float64
		Status string
	}

	query := `
		SELECT user_id, merchant_id, amount, status
		FROM payflow_payments
		WHERE reference = $1
		FOR UPDATE
	`

	tx, err := h.DB.Begin()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal server error"})
		return 
	}
	defer tx.Rollback()

	err = tx.QueryRow(query, ref).Scan(
		&payment.UserId, &payment.MerchantId, &payment.Amount, &payment.Status)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Internal server error"})
		return 
	}

	if payment.Status != "CREATED"{
		c.JSON(http.StatusConflict, gin.H{"error":"payment already processed"})
		return 
	}

	updateProcessingQuery := `UPDATE payflow_payments SET status = 'PROCESSING' WHERE reference = $1`
	_, err = tx.Exec(updateProcessingQuery, ref)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal server error"})
		return 
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_, err = h.WalletClient.Client.Debit(ctx, &walletpb.DebitRequest{
		UserId: payment.UserId,
		Amount: payment.Amount,
		Reference: ref,
	})

	if err != nil{
		h.failPayment(tx, ref)	
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Insufficient funds"})
		return 
	}

	_, err = h.WalletClient.Client.Credit(ctx, &walletpb.CreditRequest{
		UserId: payment.MerchantId,
		Amount: payment.Amount,
		Reference: ref,
	})
	if err != nil{
		h.failPayment(tx, ref)	
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Merchant credit failed"})
		return 
	}


	updateQuery := `UPDATE payflow_payments SET status = 'SUCCESS' WHERE reference = $1`
	_, err = tx.Exec(updateQuery, ref)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Internal server error"})
		return 
	}

	tx.Commit()

	c.JSON(http.StatusAccepted, gin.H{"status":"SUCCESS"})
}

func (h *PaymentHandler) failPayment(tx *sql.Tx, ref string)error {
	query := `UPDATE payflow_payments SET status = 'FAILED' WHERE reference = $1`
	_, err := tx.Exec(query, ref)
	return err
}