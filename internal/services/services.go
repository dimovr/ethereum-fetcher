package services

import (
	"ethereum_fetcher/internal/config"
	"ethereum_fetcher/internal/services/auth"
	"ethereum_fetcher/internal/services/transactions"
	"fmt"

	"gorm.io/gorm"
)

type Services struct {
	Auth auth.AuthService
	Tx   transactions.TxnService
}

func Init(db *gorm.DB, cfg config.Config) (*Services, error) {
	authService, err := auth.NewAuthService(db, cfg.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth service:  %w", err)
	}

	txService, err := transactions.NewTxnService(db, cfg.EthNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create txn service:  %w", err)
	}

	return &Services{Auth: authService, Tx: txService}, nil
}
