package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
)

type HTTPHandler struct {
	WalletHandler *WalletHandler
}

func NewHTTPHandler(walletHandler *WalletHandler) *HTTPHandler {
	return &HTTPHandler{WalletHandler: walletHandler}
}

func (h *HTTPHandler) GetBalance(c *gin.Context) {
	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.WalletHandler.GetBalance(ctx, &walletpb.GetBalanceRequest{
		UserId: userId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": resp.Balance})
}

func (h *HTTPHandler) GetTransactions(c *gin.Context) {
	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	
	limit := int32(20)
	offset := int32(0)

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = int32(val)
		}
	}

	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = int32(val)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.WalletHandler.GetTransactions(ctx, &walletpb.GetTransactionsRequest{
		UserId: userId,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": resp.Transactions,
		"limit":        limit,
		"offset":       offset,
	})
}