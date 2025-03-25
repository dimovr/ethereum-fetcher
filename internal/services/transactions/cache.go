package transactions

import (
	"time"

	types "ethereum_fetcher/internal/services/transactions/types"
	"ethereum_fetcher/pkg/logging"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type TxnCache interface {
	Set(txn *types.DbTxn)
	SetMany(txns []types.DbTxn)
	Get(hash string) (*types.DbTxn, bool)
	GetMany(hashes []string) types.TxnsResult
}

type cacheimpl struct {
	cache  *cache.Cache
	logger *logrus.Logger
}

func NewTxnCache() TxnCache {
	logger := logging.New()
	c := cache.New(1*time.Hour, 1*time.Minute)
	return &cacheimpl{cache: c, logger: logger}
}

func (tc *cacheimpl) Get(hash string) (*types.DbTxn, bool) {
	if value, found := tc.cache.Get(hash); found {
		return value.(*types.DbTxn), true
	}
	return nil, false
}

func (tc *cacheimpl) Set(txn *types.DbTxn) {
	tc.logger.Debugf("Caching transaction for hash '%s '", txn.TransactionHash)
	tc.cache.Set(txn.TransactionHash, txn, cache.DefaultExpiration)
}

func (tc *cacheimpl) SetMany(txns []types.DbTxn) {
	for _, txn := range txns {
		tc.Set(&txn)
	}
}

func (tc *cacheimpl) GetMany(hashes []string) types.TxnsResult {
	var results = make([]types.DbTxn, 0, len(hashes))
	var existingHashes = make([]string, 0, len(hashes))
	var missingHashes = make([]string, 0, len(hashes))

	for _, hash := range hashes {
		if tx, found := tc.Get(hash); found {
			results = append(results, *tx)
			existingHashes = append(existingHashes, hash)
		} else {
			missingHashes = append(missingHashes, hash)
		}
	}

	return types.TxnsResult{
		ExistingTxns:   results,
		ExistingHashes: existingHashes,
		MissingHashes:  missingHashes,
	}
}
