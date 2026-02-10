package api

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"

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
	var existingRef string
	var existingHash string
	var merchantUserId string
	userId := c.GetString("user_id")
	ref := uuid.NewString()

	idemKey := c.GetHeader("Idempotency-Key")
	if idemKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header required"})
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request"})
		return 
	}

	reqHash, err := hashRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	defer tx.Rollback()


	selectQuery := `
		SELECT payment_reference, request_hash
		FROM payflow_idempotency_keys
		WHERE idempotency_key = $1 AND user_id = $2
		`
	err = tx.QueryRow(selectQuery, idemKey, userId).Scan(&existingRef, &existingHash)

	if err == nil{
		if existingHash != reqHash {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Idempotency key reuse with different payload",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"reference": existingRef,
			"idempotent": true,
		})
		return
	}


	if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	merchantQuery := `SELECT user_id FROM payflow_merchants WHERE id = $1`
	err = tx.QueryRow(merchantQuery, req.MerchantID).Scan(&merchantUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid merchant ID"})
		return
	}

	// insertQuery := `INSERT INTO payflow_payments (user_id, merchant_id, amount, status, reference)
	// VALUES ($1, $2, $3, 'CREATED', $4)
	// `	

	insertQuery := `INSERT INTO payflow_payments (user_id, merchant_id, merchant_user_id, amount, status, reference)
	VALUES ($1, $2, $3, $4, 'CREATED', $5)
	`	
	_,  err = tx.Exec(insertQuery, userId, req.MerchantID, merchantUserId, req.Amount, ref)

	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to create payment"})
		return 
	}

	insertQuery = `
		INSERT INTO payflow_idempotency_keys (idempotency_key, user_id, payment_reference, request_hash)
		VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(insertQuery, idemKey, userId, ref, reqHash)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store idempotency key"})
		return
	}

	err = tx.Commit()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
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

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	reference := c.Query("reference")
	
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reference is required"})
		return
	}

	var payment struct {
		ID         string  `json:"id"`
		UserID     string  `json:"user_id"`
		MerchantID string  `json:"merchant_id"`
		Amount     float64 `json:"amount"`
		Status     string  `json:"status"`
		Reference  string  `json:"reference"`
		CreatedAt  string  `json:"created_at"`
	}

	query := `
		SELECT id, user_id, merchant_id, amount, status, reference, created_at
		FROM payflow_payments
		WHERE reference = $1
	`

	err := h.DB.QueryRow(query, reference).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.MerchantID,
		&payment.Amount,
		&payment.Status,
		&payment.Reference,
		&payment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment": payment})
}



func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	userId := c.GetString("user_id")
	
	
	limit := 20
	offset := 0
	
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}
	
	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	query := `
		SELECT id, user_id, merchant_id, amount, status, reference, created_at
		FROM payflow_payments
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := h.DB.Query(query, userId, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}
	defer rows.Close()

	var payments []map[string]interface{}

	for rows.Next() {
		var id, userID, merchantID, status, reference, createdAt string
		var amount float64

		err := rows.Scan(&id, &userID, &merchantID, &amount, &status, &reference, &createdAt)
		if err != nil {
			continue
		}

		payments = append(payments, map[string]interface{}{
			"id":          id,
			"user_id":     userID,
			"merchant_id": merchantID,
			"amount":      amount,
			"status":      status,
			"reference":   reference,
			"created_at":  createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
		"limit":    limit,
		"offset":   offset,
	})
}

func hashRequest(v interface{})(string, error){
	b, err := json.Marshal(v)
	if err != nil{
		return "", err
	}


	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}