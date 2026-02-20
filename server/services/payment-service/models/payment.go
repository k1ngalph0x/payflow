package models

import "time"

type PaymentStatus string

const (
	PaymentStatusCreated        PaymentStatus = "CREATED"
	PaymentStatusProcessing     PaymentStatus = "PROCESSING"
	PaymentStatusFundsCaptured  PaymentStatus = "FUNDS_CAPTURED"
	PaymentStatusSettled        PaymentStatus = "SETTLED"
	PaymentStatusFailed         PaymentStatus = "FAILED"
)

type Payment struct {
	ID             string        `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         string        `gorm:"type:uuid;not null;index"                       json:"user_id"`
	MerchantID     string        `gorm:"type:uuid;not null;index"                       json:"merchant_id"`
	MerchantUserID string        `gorm:"type:uuid;not null"                             json:"merchant_user_id"`
	Amount         float64       `gorm:"not null"                                       json:"amount"`
	Status         PaymentStatus `gorm:"type:varchar(20);not null;default:'CREATED'"    json:"status"`
	Reference      string        `gorm:"uniqueIndex;not null"                           json:"reference"`
	CreatedAt      time.Time     `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime"                                 json:"updated_at"`
}

type IdempotencyKey struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	IdempotencyKey   string    `gorm:"uniqueIndex;not null"`
	UserID           string    `gorm:"type:uuid;not null"`
	PaymentReference string    `gorm:"not null"`
	RequestHash      string    `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}