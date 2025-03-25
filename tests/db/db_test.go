package db_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ethereum_fetcher/db/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.UserTransaction{})
	require.NoError(t, err)

	return db
}

func TestUserOperations(t *testing.T) {
	db := setupTestDB(t)

	t.Run("CreateAndRetrieveUser", func(t *testing.T) {
		// Create a new user
		user := models.User{
			Username:     "testuser",
			PasswordHash: "hashedpassword",
			CreatedAt:    time.Now(),
		}

		// Save the user
		result := db.Create(&user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// Retrieve the user
		var retrievedUser models.User
		result = db.First(&retrievedUser, user.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, user.Username, retrievedUser.Username)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Create a user
		user := models.User{
			Username:     "updateuser",
			PasswordHash: "initialpassword",
			CreatedAt:    time.Now(),
		}
		db.Create(&user)

		// Update the user
		user.PasswordHash = "newhashedpassword"
		result := db.Save(&user)
		assert.NoError(t, result.Error)

		// Verify the update
		var updatedUser models.User
		db.First(&updatedUser, user.ID)
		assert.Equal(t, "newhashedpassword", updatedUser.PasswordHash)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		// Create a user
		user := models.User{
			Username:     "deleteuser",
			PasswordHash: "password",
			CreatedAt:    time.Now(),
		}
		db.Create(&user)

		// Delete the user
		result := db.Delete(&user)
		assert.NoError(t, result.Error)

		// Verify deletion
		var deletedUser models.User
		err := db.First(&deletedUser, user.ID).Error
		assert.Error(t, err, "User should not be found after deletion")
	})
}

func TestTransactionOperations(t *testing.T) {
	db := setupTestDB(t)

	t.Run("CreateAndRetrieveTransaction", func(t *testing.T) {
		// Create a new transaction
		transaction := models.Transaction{
			TransactionHash:   "0x1234567890abcdef",
			TransactionStatus: 1,
			BlockHash:         "0xabcdef1234567890",
			BlockNumber:       12345,
			FromAddress:       "0xSenderAddress",
			ToAddress:         ptr("0xReceiverAddress"),
			CreatedAt:         time.Now(),
		}

		// Save the transaction
		result := db.Create(&transaction)
		assert.NoError(t, result.Error)

		// Retrieve the transaction
		var retrievedTransaction models.Transaction
		result = db.Where("transaction_hash = ?", transaction.TransactionHash).First(&retrievedTransaction)
		assert.NoError(t, result.Error)
		assert.Equal(t, transaction.BlockNumber, retrievedTransaction.BlockNumber)
	})

	t.Run("UpdateTransaction", func(t *testing.T) {
		// Create a transaction
		transaction := models.Transaction{
			TransactionHash:   "0x9876543210fedcba",
			TransactionStatus: 0,
			CreatedAt:         time.Now(),
		}
		db.Create(&transaction)

		// Update the transaction
		transaction.TransactionStatus = 2
		result := db.Save(&transaction)
		assert.NoError(t, result.Error)

		// Verify the update
		var updatedTransaction models.Transaction
		db.First(&updatedTransaction, "transaction_hash = ?", transaction.TransactionHash)
		assert.Equal(t, 2, updatedTransaction.TransactionStatus)
	})
}

func TestUserTransactionOperations(t *testing.T) {
	db := setupTestDB(t)

	t.Run("CreateUserTransaction", func(t *testing.T) {
		// Create a user and a transaction first
		user := models.User{
			Username:     "txuser",
			PasswordHash: "password",
			CreatedAt:    time.Now(),
		}
		db.Create(&user)

		transaction := models.Transaction{
			TransactionHash:   "0xuniquehash123",
			TransactionStatus: 1,
			CreatedAt:         time.Now(),
		}
		db.Create(&transaction)

		// Create a user transaction
		userTransaction := models.UserTransaction{
			UserId:          user.ID,
			TransactionHash: transaction.TransactionHash,
			RequestedAt:     time.Now(),
		}
		result := db.Create(&userTransaction)
		assert.NoError(t, result.Error)
		assert.NotZero(t, userTransaction.ID)

		// Retrieve and verify
		var retrievedUserTransaction models.UserTransaction
		result = db.Where("user_id = ? AND transaction_hash = ?", user.ID, transaction.TransactionHash).First(&retrievedUserTransaction)
		assert.NoError(t, result.Error)
		assert.Equal(t, user.ID, retrievedUserTransaction.UserId)
	})
}

func ptr[T any](v T) *T {
	return &v
}
