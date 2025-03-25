package api

import "math/big"

type AuthToken = string

type AuthRequest struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type AuthResponse struct {
	Token *AuthToken `json:"token"`
}

type Transaction struct {
	TransactionHash   string   `json:"transactionHash"`
	TransactionStatus int      `json:"transactionStatus"`
	BlockHash         string   `json:"blockHash"`
	BlockNumber       *big.Int `json:"blockNumber"`
	From              string   `json:"from"`
	To                *string  `json:"to,omitempty"`
	ContractAddress   *string  `json:"contractAddress,omitempty"`
	LogsCount         int      `json:"logsCount"`
	Input             string   `json:"input"`
	Value             string   `json:"value"`
}

type TransactionResponse struct {
	Transactions *[]Transaction `json:"transactions"`
}

type Error struct {
	Msg string `json:"error"`
}
