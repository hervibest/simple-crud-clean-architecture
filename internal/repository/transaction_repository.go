package repository

import (
	"simple-crud-clean-architecture/internal/entity"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Repository[entity.Transaction]
	Log *logrus.Logger
}

func NewTransactionRepository(log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		Log: log,
	}
}

func (r *TransactionRepository) FindByTrxID(db *gorm.DB, transaction *entity.Transaction, trxID uuid.UUID) error {
	return db.Where("trx_id = ?", trxID).First(transaction).Error
}

func (r *TransactionRepository) GetTransactionWithDetails(db *gorm.DB, trxID uuid.UUID) (*entity.Transaction, error) {
	var transaction entity.Transaction

	err := db.Joins("User").Joins("Course").
		Where("trx_id = ?", trxID).
		First(&transaction).Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
