package models

import "time"

type Payment struct {
	ID         string    `json:"id"`
	UserId     string    `json:"user_id"`
	MerchantID string    `json:"merchant_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	Reference  string    `json:"reference"`
	CreatedAt  time.Time `json:"created_at"`
}