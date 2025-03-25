package routes

import (
	"github.com/gin-gonic/gin"

	"ethereum_fetcher/internal/handlers"
	"ethereum_fetcher/internal/services"
)

func SetupRoutes(services *services.Services, r *gin.Engine) {
	r.Use(gin.Recovery())

	r.POST("/lime/authenticate", handlers.Authenticate(services.Auth))

	authMiddleware := handlers.JwtMiddleware(services.Auth)
	txHandler := handlers.NewTxnHandler(services.Tx)

	r.GET("/lime/eth", authMiddleware, txHandler.FetchTransactions)
	r.GET("/lime/eth/:rlphex", authMiddleware, txHandler.FetchTransactionsByRLP)
	r.GET("/lime/all", txHandler.AllTransactions)
	r.GET("/lime/my", authMiddleware, txHandler.ForUser)

}
