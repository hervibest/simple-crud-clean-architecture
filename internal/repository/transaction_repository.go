package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

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

	err := db.Joins("User").Joins("Course").Joins("Voucher").
		Where("trx_id = ?", trxID).
		First(&transaction).Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) EmployeeSearch(db *gorm.DB, request *model.SearchTransactionRequest) ([]entity.Transaction, *model.PageMetadata, error) {
	var courses []entity.Transaction

	var totalItems int64
	if err := db.Model(&entity.Transaction{}).Scopes(r.FilterTransaction(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	if err := db.Debug().Scopes(r.FilterTransaction(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Order("transactions.created_at " + request.OrderBy).Find(&courses).Error; err != nil {
		return nil, nil, err
	}

	return courses, pageMetadata, nil
}

func (r *TransactionRepository) FilterTransaction(request *model.SearchTransactionRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		tx = tx.Joins("User").Joins("Course").Joins("Voucher")

		// if name := request.TransactionID; name != "" {
		// 	tx = tx.Where("trx_id = ?", name)
		// }

		if email := request.UserEmail; email != "" {
			email = "%" + email + "%"
			tx = tx.Where("\"User\".email LIKE ?", email)
		}

		if courseName := request.CourseName; courseName != "" {
			courseName = "%" + courseName + "%"
			tx = tx.Where("\"Course\".name LIKE ?", courseName)
		}

		if voucherName := request.VoucherName; voucherName != "" {
			voucherName = "%" + voucherName + "%"
			tx = tx.Where("\"Voucher\".name LIKE ?", voucherName)
		}

		return tx
	}
}
