package models

import (
	"time"
)

type User struct {
	ID           uint64 `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null" json:"-"`
	CreatedAt    time.Time
}

type Transaction struct {
	TransactionHash   string `gorm:"primaryKey;size:66"`
	TransactionStatus int
	BlockHash         string
	BlockNumber       uint64
	FromAddress       string
	ToAddress         *string
	ContractAddress   *string
	LogsCount         int
	Input             string
	Value             string
	CreatedAt         time.Time
}

// UserTransaction stores which users requested which transactions
type UserTransaction struct {
	ID              uint64 `gorm:"primaryKey"`
	UserId          uint64 `gorm:"size:255;not null;"`
	TransactionHash string `gorm:"size:66;not null;uniqueIndex:idx_user_transaction"`
	RequestedAt     time.Time
}
