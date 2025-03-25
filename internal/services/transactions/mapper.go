package transactions

import (
	"math/big"
	"time"

	custom "ethereum_fetcher/internal/services/transactions/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func toDbTxns(txns []custom.EthTxnWithReceipt) ([]custom.DbTxn, error) {
	dbTxs := make([]custom.DbTxn, 0, len(txns))

	for _, pair := range txns {
		dbTx := toDbTxn(pair.Txn, pair.Receipt)
		dbTxs = append(dbTxs, dbTx)
	}

	return dbTxs, nil
}

func toDbTxn(tx *custom.EthTxn, receipt *custom.EthReceipt) custom.DbTxn {
	from, _ := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)

	var to *string
	if tx.To() != nil {
		toAddress := tx.To().Hex()
		to = &toAddress
	}

	var contractAddress *string
	if receipt.ContractAddress != (common.Address{}) {
		addr := receipt.ContractAddress.Hex()
		contractAddress = &addr
	}

	return custom.DbTxn{
		TransactionHash:   tx.Hash().Hex(),
		TransactionStatus: int(receipt.Status),
		BlockHash:         receipt.BlockHash.Hex(),
		BlockNumber:       receipt.BlockNumber.Uint64(),
		FromAddress:       from.Hex(),
		ToAddress:         to,
		ContractAddress:   contractAddress,
		LogsCount:         len(receipt.Logs),
		Input:             common.Bytes2Hex(tx.Data()),
		Value:             tx.Value().String(),
		CreatedAt:         time.Now(),
	}
}

func toApiTxns(dbTxs []custom.DbTxn) []custom.ApiTxn {
	apiTxs := make([]custom.ApiTxn, 0, len(dbTxs))

	for _, dbTx := range dbTxs {
		apiTx := toApiTxn(&dbTx)
		apiTxs = append(apiTxs, apiTx)
	}

	return apiTxs
}

func toApiTxn(txn *custom.DbTxn) custom.ApiTxn {
	return custom.ApiTxn{
		TransactionHash:   txn.TransactionHash,
		TransactionStatus: txn.TransactionStatus,
		BlockHash:         txn.BlockHash,
		BlockNumber:       big.NewInt(int64(txn.BlockNumber)),
		From:              txn.FromAddress,
		To:                txn.ToAddress,
		ContractAddress:   txn.ContractAddress,
		LogsCount:         txn.LogsCount,
		Input:             txn.Input,
		Value:             txn.Value,
	}
}
