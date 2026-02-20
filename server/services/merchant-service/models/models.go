package models

import "time"

type MerchantStatus string

const (
	MerchantStatusActive   MerchantStatus = "ACTIVE"
	MerchantStatusInactive MerchantStatus = "INACTIVE"
	MerchantStatusPending  MerchantStatus = "PENDING"
)

type Merchant struct {
	ID           string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       string         `gorm:"type:uuid;not null;uniqueIndex"                 json:"user_id"`
	BusinessName string         `gorm:"not null"                                       json:"business_name"`
	Status       MerchantStatus `gorm:"type:varchar(20);not null;default:'ACTIVE'"     json:"status"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"                                 json:"updated_at"`
}