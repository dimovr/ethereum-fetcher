package transactions

import (
	"ethereum_fetcher/internal/services/transactions/ethereum"
	types "ethereum_fetcher/internal/services/transactions/types"
	"ethereum_fetcher/pkg/logging"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TxnService interface {
	ByHashes(hashes []string, userId uint64) ([]types.ApiTxn, error)
	FromRLPHex(rlpHex string, userId uint64) ([]types.ApiTxn, error)
	ForUser(userId uint64) ([]types.ApiTxn, error)
	All() ([]types.ApiTxn, error)
}

type impl struct {
	repo   TxnRepo
	eth    ethereum.EthService
	cache  TxnCache
	logger *logrus.Logger
}

func NewTxnService(db *gorm.DB, ethNodeURL string) (TxnService, error) {
	logger := logging.New()
	TxnRepo := NewTxnRepo(db)

	ethService, err := ethereum.NewEthereumService(ethNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ethereum service:  %w", err)
	}

	cache := NewTxnCache()

	return &impl{repo: TxnRepo, eth: ethService, cache: cache, logger: logger}, nil
}

func (s *impl) FromRLPHex(rlpHex string, userId uint64) ([]types.ApiTxn, error) {
	hashes, err := s.eth.DecodeHashes(rlpHex)
	if err != nil {
		return nil, types.InvalidRlpEncoding
	}
	return s.ByHashes(hashes, userId)
}

func (s *impl) ByHashes(hashes []string, userId uint64) ([]types.ApiTxn, error) {
	// todo: implement ishex check
	s.recordUserTransactions(hashes, userId)

	cacheResult := s.loadFromCache(hashes)
	if len(cacheResult.MissingHashes) == 0 {
		s.logger.Infof("Fetched all transactions from the cache: '%s'", hashes)
		return toApiTxns(cacheResult.ExistingTxns), nil
	} else {
		s.logger.Infof("Transactions for hashes: '%s' found in the cache", cacheResult.ExistingHashes)
	}

	dbResult, err := s.loadFromDb(cacheResult.MissingHashes)
	if err != nil {
		s.logger.Infof("failed to load existing transactions for hashes: '%s'", cacheResult.MissingHashes)
		return nil, types.NewTxnError("failed to load existing transactions")
	}
	s.cacheTxns(dbResult.ExistingTxns)

	if len(dbResult.MissingHashes) == 0 {
		s.logger.Infof("Fetched all transactions from the database: '%s'", cacheResult.MissingHashes)
		return toApiTxns(dbResult.ExistingTxns), nil
	} else {
		s.logger.Infof("Transactions for hashes: '%s' fetched from the database", dbResult.ExistingHashes)
	}

	newTxns, err := s.getFromEth(dbResult.MissingHashes)
	if err != nil || len(newTxns) == 0 {
		return nil, types.FailedToFetchTransaction
	}

	storeErr := s.storeTxns(newTxns)
	if storeErr != nil {
		return nil, storeErr
	}
	s.cacheTxns(newTxns)

	all := append(append(cacheResult.ExistingTxns, dbResult.ExistingTxns...), newTxns...)

	return toApiTxns(all), nil
}

func (s *impl) loadFromCache(hashes []string) types.TxnsResult {
	return s.cache.GetMany(hashes)
}

func (s *impl) loadFromDb(hashes []string) (types.TxnsResult, error) {
	existingTxns, err := s.repo.GetForHashes(hashes)
	if err != nil {
		return types.TxnsResult{}, err
	}

	existingHashes := make([]string, 0, len(existingTxns))
	existingHashSet := make(map[string]bool, len(existingTxns))
	missingHashes := make([]string, 0, len(hashes))

	for _, tx := range existingTxns {
		existingHashes = append(existingHashes, tx.TransactionHash)
		existingHashSet[tx.TransactionHash] = true
	}

	for _, hash := range hashes {
		if !existingHashSet[hash] {
			missingHashes = append(missingHashes, hash)
		}
	}

	return types.TxnsResult{
		ExistingTxns:   existingTxns,
		ExistingHashes: existingHashes,
		MissingHashes:  missingHashes,
	}, nil
}

func (s *impl) getFromEth(hashes []string) ([]types.DbTxn, error) {
	s.logger.Infof("Fetching transactions for hashes: '%s' from the Ethereum node", hashes)
	ethTxnsResult, err := s.eth.ByHashes(hashes)
	if err != nil {
		s.logger.Errorf("failed to fetch missing transactions for hashes: '%s':  %w", hashes, err)
		return nil, types.NewEthError("failed to fetch missing transactions")
	}

	newTxns, err := toDbTxns(ethTxnsResult)
	if err != nil {
		s.logger.Errorf("failed to convert transactions to DB models for hashes: '%s':  %w", hashes, err)
		return nil, types.NewTxnError("failed to convert transactions to DB models")

	}

	return newTxns, nil
}

func (s *impl) cacheTxns(txns []types.DbTxn) {
	s.cache.SetMany(txns)
}

func (s *impl) storeTxns(txns []types.DbTxn) error {
	s.logger.Infof("Saving new transactions for hashes: '%s' to the database", txns)
	if err := s.repo.Save(txns); err != nil {
		s.logger.Errorf("failed to store new transactions for hashes: '%s':  %w", txns, err)
		return types.NewTxnError("failed to store new transactions")
	}
	return nil
}

func (s *impl) recordUserTransactions(hashes []string, userId uint64) error {
	if userId == 0 {
		return nil
	}
	s.logger.Infof("Storing user transactions for user: '%s' and hashes: '%s'", userId, hashes)
	return s.repo.AddUserTransactions(hashes, userId)
}

func (s *impl) ForUser(userId uint64) ([]types.ApiTxn, error) {
	txns, err := s.repo.GetUserTransactions(userId)
	if err != nil {
		s.logger.Errorf("failed to fetch user transactions for user '%s':  %v", userId, err)
		return nil, types.NewTxnError("failed to fetch user transactions")
	}

	return toApiTxns(txns), nil
}

func (s *impl) All() ([]types.ApiTxn, error) {
	txns, err := s.repo.GetAll()
	if err != nil {
		s.logger.Errorf("failed to fetch all transactions:  %w", err)
		return nil, types.NewTxnError("failed to fetch all transactions")
	}
	return toApiTxns(txns), nil
}
