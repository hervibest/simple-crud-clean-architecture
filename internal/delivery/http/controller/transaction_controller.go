package http

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/enum"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	Log       *logrus.Logger
	UseCase   *usecase.TransactionUseCase
	Midtrans  *helper.MidtransClient
	Validator helper.CustomValidator
}

func NewTransactionController(useCase *usecase.TransactionUseCase, log *logrus.Logger, midtransHelper *helper.MidtransClient, validator helper.CustomValidator) *TransactionController {
	return &TransactionController{
		UseCase:   useCase,
		Log:       log,
		Midtrans:  midtransHelper,
		Validator: validator,
	}
}

func (c *TransactionController) List(ctx *fiber.Ctx) error {

	orderBy := strings.ToLower(ctx.Query("order_by", "asc"))

	if orderBy != "asc" && orderBy != "desc" {
		orderBy = "asc"
	}

	request := &model.SearchTransactionRequest{
		// TransactionID: ctx.Query("transaction_uuid", ""),
		CourseName:  ctx.Query("courseName", ""),
		UserEmail:   ctx.Query("userEmail", ""),
		VoucherName: ctx.Query("voucherName", ""),
		OrderBy:     orderBy,
		Page:        ctx.QueryInt("page", 1),
		Size:        ctx.QueryInt("size", 10),
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	responses, pageMetadata, err := c.UseCase.EmployeeSearch(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching course")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.TransactionResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *TransactionController) GetTransaction(ctx *fiber.Ctx) error {
	transactionID := ctx.Params("trxId")
	parsedUUID, err := uuid.Parse(transactionID)
	if err != nil {
		c.Log.WithError(err).Error("error parsing transaction UUID")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid transaction UUID")
	}

	transaction, err := c.UseCase.GetDetailTransaction(ctx.UserContext(), parsedUUID)
	if err != nil {
		c.Log.WithError(err).Error("error getting transaction detail")
		return fiber.ErrNotFound
	}

	return ctx.JSON(model.DataResponse[*model.TransactionResponse]{
		Success: true,
		Data:    transaction,
	})
}

func (c *TransactionController) Buy(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateTransactionRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(request.CourseUUIDStr)
	if err != nil {
		c.Log.WithError(err).Error("error parsing course UUID")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid course UUID")
	}

	request.CourseUUID = parsedUUID
	request.UserID = auth.Id

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	transaction, err := c.UseCase.CreateTransaction(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating transaction")
		return err
	}

	email := auth.Email
	response, err := c.UseCase.GetPaymentTransactionToken(ctx.UserContext(), transaction, email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error getting payment transaction token")
	}

	return ctx.JSON(model.DataResponse[*model.SnapshotTokenResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *TransactionController) Notify(ctx *fiber.Ctx) error {

	request := new(model.MidtransNotifyRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	helper.SanitiseStruct(request)

	parsedOrderID, err := uuid.Parse(request.OrderID)

	transaction, err := c.UseCase.GetByTrxID(ctx.UserContext(), parsedOrderID)
	if err != nil {
		c.Log.WithError(err).Error("error transaction not found")
		return fiber.ErrNotFound
	}

	signatureToCompare := transaction.TrxID.String() + request.StatusCode + request.GrossAmount + c.Midtrans.GetMidtransKey()

	hash := sha512.New()
	hash.Write([]byte(signatureToCompare))
	hashedSignature := hex.EncodeToString(hash.Sum(nil))

	requestIsValid := hashedSignature == request.SignatureKey
	if !requestIsValid {
		c.Log.Error("error invalid signature key")
		return fiber.ErrForbidden
	}
	transactionStatus := enum.MidtransPaymentStatus(request.TransactionStatus)

	if enum.PaymentStatusPending == transactionStatus {
		return ctx.JSON(model.WebResponse{
			Success: true,
		})
	}

	body := ctx.Body()

	updateRequest := &model.UpdateTransactionStatus{
		TransactionID:            parsedOrderID,
		Status:                   transactionStatus,
		ExternalCallbackResponse: json.RawMessage(body),
	}

	err = c.UseCase.UpdateTransactionStatus(ctx.UserContext(), updateRequest)
	if err != nil {
		c.Log.WithError(err).Error("error updating transaction status")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.TransactionResponse]{
		Success: true,
	})
}
