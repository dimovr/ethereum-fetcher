package ethereum

import (
	"context"
	"ethereum_fetcher/pkg/logging"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	custom "ethereum_fetcher/internal/services/transactions/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthService interface {
	DecodeHashes(rlpHex string) ([]string, error)
	ByHashes(hashes []string) ([]custom.EthTxnWithReceipt, error)
}

type impl struct {
	client *ethclient.Client
	logger *logrus.Logger
}

func NewEthereumService(ethNodeURL string) (EthService, error) {
	logger := logging.New()

	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create ethereum client:  %w", err)
	}

	return &impl{client: client, logger: logger}, nil
}

func (s *impl) DecodeHashes(rlpHex string) ([]string, error) {
	return DecodeHashes(rlpHex)
}

func (s *impl) ByHashes(hashes []string) ([]custom.EthTxnWithReceipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sem := make(chan struct{}, 10)

	var (
		txns    = make([]custom.EthTxnWithReceipt, 0, len(hashes))
		errChan = make(chan error, len(hashes))
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	for _, hash := range hashes {
		wg.Add(1)

		go func(hash string) {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}

			result, err := s.fetchSingle(ctx, hash)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch tx %s: %w", hash, err)
				return
			}

			mu.Lock()
			txns = append(txns, result)
			mu.Unlock()
		}(hash)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	select {
	case err := <-errChan:
		return nil, err
	default:
		return txns, nil
	}
}

func (s *impl) fetchSingle(ctx context.Context, hash string) (custom.EthTxnWithReceipt, error) {
	txHash := common.HexToHash(hash)

	tx, isPending, err := s.client.TransactionByHash(ctx, txHash)
	if err != nil {
		s.logger.Errorf("Error fetching transaction '%s': %v", txHash.Hex(), err)
		return custom.EthTxnWithReceipt{}, err
	}

	if isPending {
		return custom.EthTxnWithReceipt{}, nil
	}

	receipt, err := s.client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		s.logger.Errorf("Error fetching receipt for '%s': %v", tx.Hash().Hex(), err)
		return custom.EthTxnWithReceipt{}, err
	}

	return custom.EthTxnWithReceipt{Txn: tx, Receipt: receipt}, nil
}
