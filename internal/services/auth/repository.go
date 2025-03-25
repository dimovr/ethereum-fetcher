package auth

import (
	"ethereum_fetcher/db/models"
	"ethereum_fetcher/pkg/passwords"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func (r *UserRepo) FindUser(username string, password string) (models.User, error) {
	var user models.User

	if err := r.db.Table("users").
		Select("users.*").
		Where("users.username = ?", username).
		First(&user).Error; err != nil {
		return user, UsernameNotFound
	}

	if err := passwords.ComparePasswords(user.PasswordHash, password); err != nil {
		return user, InvalidPassword
	}

	return user, nil
}
