package transactions

import (
	"errors"
	"ethereum_fetcher/db/models"
	"time"

	"gorm.io/gorm"
)

type TxnRepo interface {
	Save(txns []models.Transaction) error
	AddUserTransactions(txnHashes []string, userId uint64) error
	GetForHashes(txnHashes []string) ([]models.Transaction, error)
	GetUserTransactions(userId uint64) ([]models.Transaction, error)
	GetAll() ([]models.Transaction, error)
}

func NewTxnRepo(db *gorm.DB) TxnRepo {
	return &repoImpl{db: db}
}

type repoImpl struct {
	db *gorm.DB
}

func (r *repoImpl) Save(txns []models.Transaction) error {
	return r.db.Create(&txns).Error
}

func (r *repoImpl) AddUserTransactions(txnHashes []string, userId uint64) error {
	newUserTxns := make([]models.UserTransaction, 0, len(txnHashes))

	for _, txnHash := range txnHashes {
		var existingUserTxn models.UserTransaction
		result := r.db.Where(&models.UserTransaction{
			UserId:          userId,
			TransactionHash: txnHash,
		}).First(&existingUserTxn)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newUserTxns = append(newUserTxns, models.UserTransaction{
				UserId:          userId,
				TransactionHash: txnHash,
				RequestedAt:     time.Now(),
			})
		}
	}

	if len(newUserTxns) > 0 {
		return r.db.Create(&newUserTxns).Error
	}

	return nil
}

func (r *repoImpl) GetForHashes(txnHashes []string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("transaction_hash IN ?", txnHashes).Find(&transactions).Error
	return transactions, err
}

func (r *repoImpl) GetUserTransactions(userId uint64) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Table("transactions").
		Select("transactions.*").
		Joins("JOIN user_transactions ON transactions.transaction_hash = user_transactions.transaction_hash").
		Where("user_transactions.user_id = ?", userId).
		Find(&transactions).Error

	return transactions, err
}

func (r *repoImpl) GetAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Find(&transactions).Error
	return transactions, err
}
