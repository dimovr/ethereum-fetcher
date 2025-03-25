package handlers

import (
	"net/http"

	"ethereum_fetcher/api"
	"ethereum_fetcher/internal/services/auth"
	"ethereum_fetcher/internal/services/transactions"

	"github.com/gin-gonic/gin"
)

type TxnHandler struct {
	txService transactions.TxnService
}

func NewTxnHandler(txService transactions.TxnService) TxnHandler {
	return TxnHandler{txService: txService}
}

func (h *TxnHandler) FetchTransactions(c *gin.Context) {
	hashes := c.QueryArray("transactionHashes")

	if len(hashes) == 0 {
		c.JSON(http.StatusBadRequest, api.Error{Msg: "'transactionHashes' is a required parameter"})
		return
	}

	user := c.GetUint64(auth.UserClaim)
	txns, err := h.txService.ByHashes(hashes, user)
	response(&txns, err)(c)
}

func (h *TxnHandler) FetchTransactionsByRLP(c *gin.Context) {
	rlpHex := c.Param("rlphex")
	if rlpHex == "" {
		c.JSON(http.StatusBadRequest, api.Error{Msg: "'rlphex' is a required parameter"})
		return
	}

	user := c.GetUint64(auth.UserClaim)
	txns, err := h.txService.FromRLPHex(rlpHex, user)
	response(&txns, err)(c)
}

func (h *TxnHandler) AllTransactions(c *gin.Context) {
	txns, err := h.txService.All()
	response(&txns, err)(c)
}

func (h *TxnHandler) ForUser(c *gin.Context) {
	user := c.GetUint64(auth.UserClaim)
	txns, err := h.txService.ForUser(user)
	response(&txns, err)(c)
}

func response(txns *[]api.Transaction, err error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err != nil {
			c.JSON(toStatusCode(err), mapError(err))
			return
		}

		c.JSON(http.StatusOK, api.TransactionResponse{Transactions: txns})
	}
}
