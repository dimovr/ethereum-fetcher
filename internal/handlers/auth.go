package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ethereum_fetcher/api"
	"ethereum_fetcher/internal/services/auth"
)

func Authenticate(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req api.AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, api.Error{Msg: "Invalid request"})
			return
		}

		jwt, err := authService.Authenticate(req)
		if err != nil {
			c.JSON(toStatusCode(err), mapError(err))
			return
		}

		c.JSON(http.StatusOK, api.AuthResponse{Token: jwt})
	}
}

const AuthTokenHeader = "AUTH_TOKEN"

func JwtMiddleware(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader(AuthTokenHeader)
		if tokenString == "" {
			c.Next()
			return
		}

		userId, err := authService.GetUserId(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, mapError(err))
			return
		}

		c.Set(auth.UserClaim, userId)
		c.Next()
	}
}
