package transactions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ethereum_fetcher/db/models"
	txns "ethereum_fetcher/internal/services/transactions"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.UserTransaction{})
	require.NoError(t, err)

	return db
}

func TestTransactionService(t *testing.T) {
	db := setupTestDB(t)
	ethNodeURL := "https://sepolia.infura.io/v3/dummy"

	// Prepare test data
	user := models.User{
		Username:     "txuser",
		PasswordHash: "hashedpassword",
		CreatedAt:    time.Now(),
	}
	err := db.Create(&user).Error
	require.NoError(t, err)

	transactions := []models.Transaction{
		{
			TransactionHash:   "0x123",
			TransactionStatus: 1,
			BlockNumber:       100,
			FromAddress:       "0xSender1",
		},
		{
			TransactionHash:   "0x456",
			TransactionStatus: 1,
			BlockNumber:       200,
			FromAddress:       "0xSender2",
		},
	}
	err = db.Create(&transactions).Error
	require.NoError(t, err)

	// Create user transactions
	userTransactions := []models.UserTransaction{
		{
			UserId:          user.ID,
			TransactionHash: "0x123",
			RequestedAt:     time.Now(),
		},
		{
			UserId:          user.ID,
			TransactionHash: "0x456",
			RequestedAt:     time.Now(),
		},
	}
	err = db.Create(&userTransactions).Error
	require.NoError(t, err)

	txService, err := txns.NewTxnService(db, ethNodeURL)
	require.NoError(t, err)

	t.Run("GetTransactionsByHashes", func(t *testing.T) {
		txns, err := txService.ByHashes([]string{"0x123", "0x456"}, user.ID)
		assert.NoError(t, err)
		assert.Len(t, txns, 2)
	})

	t.Run("GetUserTransactions", func(t *testing.T) {
		txns, err := txService.ForUser(user.ID)
		assert.NoError(t, err)
		assert.Len(t, txns, 2)
	})

	t.Run("GetAllTransactions", func(t *testing.T) {
		txns, err := txService.All()
		assert.NoError(t, err)
		assert.Len(t, txns, 2)
	})
}
