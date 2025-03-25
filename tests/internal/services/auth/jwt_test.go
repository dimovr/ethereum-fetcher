package auth

import (
	"strconv"
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ethereum_fetcher/internal/services/auth"
)

func TestJwtManager(t *testing.T) {
	secretKey := "test_secret_key"
	jwtManager := auth.NewJwtManager(secretKey)
	userId := uint64(1)

	t.Run("GenerateJwt", func(t *testing.T) {
		tokenString, err := jwtManager.GenerateJwt(userId)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)
	})

	t.Run("ParseJwt", func(t *testing.T) {
		tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIn0.jtMjdef2GATPLXke5Qj7TME1Yauo2e2ybP8Cern1eDg"

		// Parse and verify user ID
		token, err := jwtManager.ParseJwt(tokenString)
		sub, err := token.Claims.GetSubject()
		parsedUserId, err := strconv.ParseUint(sub, 10, 64)
		require.NoError(t, err)
		assert.Equal(t, userId, parsedUserId)
	})

	t.Run("InvalidSecretKey", func(t *testing.T) {
		tokenString, err := jwtManager.GenerateJwt(userId)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("invalid_key"), nil
		})
		assert.Error(t, err)
		assert.False(t, token.Valid)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		invalidToken := "invalid.token.string"
		_, err := jwtManager.ParseJwt(invalidToken)
		assert.Error(t, err)
	})

}
