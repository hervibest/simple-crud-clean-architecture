package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CertificateCatController struct {
	UseCase   *usecase.CertifCatUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewCertificateCatController(useCase *usecase.CertifCatUseCase, log *logrus.Logger, validator helper.CustomValidator) *CertificateCatController {
	return &CertificateCatController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
	}
}

func (c *CertificateCatController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCertifCatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating Certificate category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateCatResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertificateCatController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCertificateCatRequest{
		Name:        ctx.Query("name", ""),
		Slug:        ctx.Query("slug", ""),
		Description: ctx.Query("description", ""),
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

	responses, pageMetadata, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching Certificate category")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CertificateCatResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CertificateCatController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("CertificateCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing Certificate category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid Certificate category uuid")
	}

	request := &model.GetCertificateCatRequest{
		UUID: parsedUUID,
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting Certificate category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateCatResponse]{
		Success: true,
		Data:    response,
	})

}

func (c *CertificateCatController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCertificateCatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("CertificateCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing Certificate category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid Certificate category uuid")
	}

	request = &model.UpdateCertificateCatRequest{
		UUID:        parsedUUID,
		Name:        request.Name,
		Description: request.Description,
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating Certificate category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateCatResponse]{
		Success: true,
		Data:    response,
	})
}
