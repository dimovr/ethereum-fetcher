package db

import (
	"ethereum_fetcher/db/models"
	"ethereum_fetcher/pkg/passwords"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitDB(dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.UserTransaction{})
	if err != nil {
		log.Fatalf("Migration failed <- %v", err)
	}

	prepopulateUsers(db)

	return db, nil
}

func prepopulateUsers(db *gorm.DB) {
	for _, name := range []string{"alice", "bob", "carol", "dave"} {
		createUser(db, name)
	}
}

func createUser(db *gorm.DB, username string) {
	user := models.User{Username: username, PasswordHash: hashPasswordOrFail(username)}

	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "username"}}, // Define the unique column
		DoUpdates: clause.Assignments(map[string]interface{}{
			"password_hash": user.PasswordHash, // Update password hash if user exists
		}),
	}).Create(&user)

	if result.Error != nil {
		log.Fatalf("Error upserting user '%s' <- %v", username, result.Error)
	} else {
		log.Printf("Upserted user: '%s'", username)
	}
}

func hashPasswordOrFail(password string) string {
	hashed, err := passwords.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password <- %v", err)
	}
	return hashed
}
