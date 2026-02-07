package models

type Merchant struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	BusinessName string `json:"business_name"`
	Status       string `json:"status"`
}