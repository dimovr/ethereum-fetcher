package transactions

import "ethereum_fetcher/internal/services/errors"

type TxnError struct {
	errors.ServiceError
}

func NewTxnError(msg string) TxnError {
	return TxnError{errors.NewServiceError(msg)}
}

type EthError struct {
	TxnError
}

func NewEthError(msg string) EthError {
	return EthError{NewTxnError(msg)}
}

var (
	FailedToFetchTransaction = NewEthError("failed to fetch transaction")
)

type RlpError struct {
	EthError
}

func rlpError(msg string) RlpError {
	return RlpError{NewEthError(msg)}
}

var (
	InvalidHexEncoding = rlpError("invalid hex encoding")
	InvalidRlpEncoding = rlpError("invalid RLP encoding")
	InvalidHashLength  = rlpError("invalid transaction hash length")
)
