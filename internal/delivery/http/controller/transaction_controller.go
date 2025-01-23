package http

import (
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	Log      *logrus.Logger
	UseCase  *usecase.TransactionUseCase
	Midtrans *helper.MidtransClient
}

func NewTransactionController(useCase *usecase.TransactionUseCase, log *logrus.Logger, midtransHelper *helper.MidtransClient) *TransactionController {
	return &TransactionController{
		UseCase:  useCase,
		Log:      log,
		Midtrans: midtransHelper,
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

	responses, pageMetadata, err := c.UseCase.EmployeeSearch(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
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
		return fiber.NewError(fiber.StatusBadRequest, "invalid course UUID")
	}

	request.CourseUUID = parsedUUID
	request.UserID = auth.Id
	request.Email = auth.Email

	response, err := c.UseCase.CreateTransaction(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error creating transaction")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SnapshotTokenResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *TransactionController) Notify(ctx *fiber.Ctx) error {

	request := new(model.UpdateTransactionWebhookRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	helper.SanitiseStruct(request)

	parsedOrderID, err := uuid.Parse(request.OrderID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid uuid")
	}

	request.ParsedOrderID = parsedOrderID
	request.Body = ctx.Body()

	err = c.UseCase.UpdateTransactionWebhook(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})

	}
	return ctx.JSON(model.DataResponse[*model.TransactionResponse]{
		Success: true,
	})
}
