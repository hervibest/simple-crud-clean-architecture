package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/enum"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model/converter"
	"time"

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
	UserRepository        *repository.UserRepository
	VoucherRepository     *repository.VoucherRepository
	Midtrans              *helper.MidtransClient
}

func NewTransactionUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	transactionRepository *repository.TransactionRepository, courseRepository *repository.CourseRepository, userRepository *repository.UserRepository, voucherRepository *repository.VoucherRepository, midtransHelper *helper.MidtransClient) *TransactionUseCase {
	return &TransactionUseCase{
		DB:                    db,
		Log:                   logger,
		Validate:              validate,
		TransactionRepository: transactionRepository,
		CourseRepository:      courseRepository,
		UserRepository:        userRepository,
		VoucherRepository:     voucherRepository,

		Midtrans: midtransHelper,
	}
}

func (c *TransactionUseCase) EmployeeSearch(ctx context.Context, request *model.SearchTransactionRequest) ([]model.TransactionResponse, *model.PageMetadata, error) {

	transactions, pageMetadata, err := c.TransactionRepository.EmployeeSearch(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting transactions")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *converter.TransactionToResponse(&transaction)
	}

	return responses, pageMetadata, nil
}

func (c *TransactionUseCase) GetDetailTransaction(ctx context.Context, trxId uuid.UUID) (*model.TransactionResponse, error) {
	transaction, err := c.TransactionRepository.GetTransactionWithDetails(c.DB, trxId)
	if err != nil {
		return nil, fiber.ErrNotFound
	}

	return converter.TransactionToResponse(transaction), nil
}

func (c *TransactionUseCase) CreateTransaction(ctx context.Context, request *model.CreateTransactionRequest) (*entity.Transaction, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)
	course, err := c.CourseRepository.FindWithDetails(tx, request.CourseUUID, false, true, false)
	if err != nil {
		c.Log.WithError(err).Error("error finding course")
		return nil, fiber.ErrBadRequest
	}

	//!!! Separation of Logic Concern in this code !!!!
	courseResponse := converter.CourseToResponse(course)

	finalPrice := courseResponse.FinalPrice
	voucher := new(entity.Voucher)
	if request.VoucherCode != "" {

		voucherReponse, err := c.VoucherRepository.FindVoucherByCodeAndCourse(tx, course.ID, request.VoucherCode)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusNotFound, "Invalid voucher code")
		}
		voucher = voucherReponse

		if voucher.Type == "PERCENT" {
			finalPrice = courseResponse.FinalPrice * (1 - voucher.Value/100)
		} else {
			finalPrice = courseResponse.FinalPrice - voucher.Value
		}
	}

	transaction := &entity.Transaction{
		TrxID:     uuid.New(),
		UserID:    request.UserID,
		CourseID:  course.ID,
		Amount:    finalPrice,
		Status:    enum.TransactionStatusPending,
		VoucherID: voucher.ID,
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

	transaction.SnapToken = response.Token
	transaction.Status = enum.TransactionStatusPending
	if err := c.TransactionRepository.Update(tx, transaction); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SnapToResponse(transaction.TrxID.String(), response), nil
}

func (c *TransactionUseCase) GetByTrxID(ctx context.Context, trxID uuid.UUID) (*entity.Transaction, error) {

	transaction := &entity.Transaction{}
	err := c.TransactionRepository.FindByTrxID(c.DB, transaction, trxID)
	if err != nil {
		return nil, fiber.ErrNotFound
	}

	return transaction, nil
}

func (c *TransactionUseCase) UpdateTransactionStatus(ctx context.Context, request *model.UpdateTransactionStatus) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return fiber.ErrBadRequest
	}

	helper.SanitiseStruct(request)

	transaction := new(entity.Transaction)
	err = c.TransactionRepository.FindByTrxID(tx, transaction, request.TransactionID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "transaction not found")
	}

	if transaction.Status == enum.TransactionStatusSuccess {
		return nil
	}

	var transactionStatus enum.TransactionStatus

	switch request.Status {
	case enum.PaymentStatusSettlement:
		transactionStatus = enum.TransactionStatusSuccess
		transaction.PaidAt = time.Now()
		transaction.SnapToken = ""
	case enum.PaymentStatusPending:
		transactionStatus = enum.TransactionStatusPending
	case enum.PaymentStatusExpire:
		transactionStatus = enum.TransactionStatusExpired
	case enum.PaymentStatusFailure:
		transactionStatus = enum.TransactionStatusFailed
	}

	transaction.Status = transactionStatus
	transaction.ExternalStatus = request.Status
	transaction.ExternalCallbackResponse = request.ExternalCallbackResponse

	if err := c.TransactionRepository.Update(tx, transaction); err != nil {
		return fiber.ErrInternalServerError
	}

	if transactionStatus != enum.TransactionStatusSuccess {
		return nil
	}

	course := new(entity.Course)
	if err := c.CourseRepository.FindById(tx, course, transaction.CourseID); err != nil {
		c.Log.WithError(err).Error("error finding course")
		return fiber.ErrInternalServerError
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, transaction.UserID); err != nil {
		c.Log.WithError(err).Error("error finding user")
		return fiber.ErrInternalServerError
	}

	if err := c.UserRepository.AppendCourse(tx, user, course); err != nil {
		c.Log.WithError(err).Error("error append course")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil

}
