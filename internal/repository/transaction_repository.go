package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/enum"

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

func (r *TransactionRepository) UpdateTransactionStatus(db *gorm.DB, transaction *entity.Transaction, id int, status enum.TransactionStatus) error {
	return db.Model(transaction).Where("id = ?", id).Update("status", status).Error
}
