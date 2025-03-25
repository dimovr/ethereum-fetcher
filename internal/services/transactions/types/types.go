package transactions

import (
	eth "github.com/ethereum/go-ethereum/core/types"

	api "ethereum_fetcher/api"
	db "ethereum_fetcher/db/models"
)

type EthTxn = eth.Transaction
type EthReceipt = eth.Receipt

type DbTxn = db.Transaction
type ApiTxn = api.Transaction

type TxnsResult struct {
	ExistingTxns   []DbTxn
	ExistingHashes []string
	MissingHashes  []string
}

type EthTxnWithReceipt struct {
	Txn     *EthTxn
	Receipt *EthReceipt
}
