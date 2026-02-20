package models

import "time"

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeCredit TransactionType = "CREDIT"
	TransactionTypeDebit  TransactionType = "DEBIT"
	TransactionStatusSuccess            TransactionStatus = "SUCCESS"
	TransactionStatusInsufficientFunds  TransactionStatus = "INSUFFICIENT_FUNDS"
)

type Wallet struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex"                 json:"user_id"`
	Balance   float64   `gorm:"not null;default:0"                             json:"balance"`
	CreatedAt time.Time `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"                                 json:"updated_at"`
}

type Transaction struct {
	ID        string            `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	WalletID  string            `gorm:"type:uuid;not null;index"                       json:"wallet_id"`
	UserID    string            `gorm:"type:uuid;not null;index"                       json:"user_id"`
	Type      TransactionType   `gorm:"type:varchar(10);not null"                      json:"type"`
	Amount    float64           `gorm:"not null"                                       json:"amount"`
	Reference string            `gorm:"not null"                                       json:"reference"`
	Status    TransactionStatus `gorm:"type:varchar(20);not null"                      json:"status"`
	CreatedAt time.Time         `gorm:"autoCreateTime"                                 json:"created_at"`
}