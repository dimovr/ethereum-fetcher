package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ethereum_fetcher/api"
	"ethereum_fetcher/db/models"
	"ethereum_fetcher/internal/services/auth"
	"ethereum_fetcher/pkg/passwords"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func TestAuthService(t *testing.T) {
	db := setupTestDB(t)
	secretKey := "test_secret_key"

	hashedPassword, _ := passwords.HashPassword("password")

	// Prepare a test user
	user := models.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}
	err := db.Create(&user).Error
	require.NoError(t, err)

	authService, err := auth.NewAuthService(db, secretKey)
	require.NoError(t, err)

	t.Run("Authenticate", func(t *testing.T) {
		token, err := authService.Authenticate(api.AuthRequest{
			Username: user.Username,
			Password: "password",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("GetUserId", func(t *testing.T) {
		// Generate a token for the user
		tokenString, err := authService.Authenticate(api.AuthRequest{
			Username: user.Username,
			Password: "password",
		})
		require.NoError(t, err)

		// Extract user ID from the token
		userId, err := authService.GetUserId(*tokenString)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, userId)
	})

	t.Run("AuthenticateInvalidCredentials", func(t *testing.T) {
		_, err := authService.Authenticate(api.AuthRequest{
			Username: user.Username,
			Password: "wrongpassword",
		})
		assert.Error(t, err)
	})
}
