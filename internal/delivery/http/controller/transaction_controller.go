package http

import (
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	Log     *logrus.Logger
	UseCase *usecase.TransactionUseCase
}

func NewTransactionController(useCase *usecase.TransactionUseCase, log *logrus.Logger) *TransactionController {
	return &TransactionController{
		UseCase: useCase,
		Log:     log,
	}
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
