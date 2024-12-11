package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/enum"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model/converter"

	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	TransactionRepository *repository.TransactionRepository
	CourseRepository      *repository.CourseRepository
	Midtrans              *helper.MidtransClient
}

func NewTransactionUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	transactionRepository *repository.TransactionRepository, courseRepository *repository.CourseRepository, midtransHelper *helper.MidtransClient) *TransactionUseCase {
	return &TransactionUseCase{
		DB:                    db,
		Log:                   logger,
		Validate:              validate,
		TransactionRepository: transactionRepository,
		CourseRepository:      courseRepository,
		Midtrans:              midtransHelper,
	}
}

func (c *TransactionUseCase) CreateTransaction(ctx context.Context, request *model.CreateTransactionRequest) (*entity.Transaction, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	helper.SanitiseStruct(request)
	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(tx, course, request.CourseUUID); err != nil {
		c.Log.WithError(err).Error("error finding course")
		return nil, fiber.ErrBadRequest
	}

	transaction := &entity.Transaction{
		TrxID:    uuid.New(),
		UserID:   request.UserID,
		CourseID: course.ID,
		Amount:   course.Price,
		Status:   enum.TransactionStatusPending,
	}

	if err := c.TransactionRepository.Create(tx, transaction); err != nil {
		c.Log.WithError(err).Error("error creating transaction")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating transaction")
		return nil, fiber.ErrInternalServerError
	}

	return transaction, nil
}

func (c *TransactionUseCase) GetPaymentTransactionToken(ctx context.Context, transaction *entity.Transaction, email string) (*model.SnapshotTokenResponse, error) {
	request := &model.MidtransSnapshotRequest{
		OrderID:  transaction.TrxID.String(),
		GrossAmt: int64(transaction.Amount),
		Email:    email,
	}

	response, err := c.Midtrans.CreateSnapshot(request)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.TransactionRepository.UpdateTransactionStatus(tx, transaction, transaction.ID, enum.TransactionStatusPending); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SnapToResponse(transaction.TrxID.String(), response), nil
}
