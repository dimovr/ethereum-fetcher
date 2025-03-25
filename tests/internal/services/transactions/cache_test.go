package transactions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"ethereum_fetcher/db/models"
	txns "ethereum_fetcher/internal/services/transactions"
)

func TestTransactionCache(t *testing.T) {
	cache := txns.NewTxnCache()

	t.Run("SetAndGetSingleTransaction", func(t *testing.T) {
		txn := &models.Transaction{
			TransactionHash: "0x123",
			BlockNumber:     100,
			FromAddress:     "0xSender",
		}

		cache.Set(txn)

		cachedTxn, ok := cache.Get("0x123")
		assert.True(t, ok)
		assert.Equal(t, txn.TransactionHash, cachedTxn.TransactionHash)
		assert.Equal(t, txn.BlockNumber, cachedTxn.BlockNumber)
	})

	t.Run("SetAndGetMultipleTransactions", func(t *testing.T) {
		txns := []models.Transaction{
			{TransactionHash: "0x456", BlockNumber: 200},
			{TransactionHash: "0x789", BlockNumber: 300},
		}

		cache.SetMany(txns)

		result := cache.GetMany([]string{"0x456", "0x789"})
		assert.Len(t, result.ExistingTxns, 2)
		assert.True(t, result.ExistingTxns[0].TransactionHash == "0x456" && result.ExistingTxns[1].TransactionHash == "0x789")
	})

	t.Run("GetNonExistentTransaction", func(t *testing.T) {
		_, ok := cache.Get("0xNONEXISTENT")
		assert.False(t, ok)
	})
}
