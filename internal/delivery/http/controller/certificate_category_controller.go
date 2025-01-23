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
	UseCase *usecase.CertifCatUseCase
	Log     *logrus.Logger
}

func NewCertificateCatController(useCase *usecase.CertifCatUseCase, log *logrus.Logger) *CertificateCatController {
	return &CertificateCatController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CertificateCatController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCertifCatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
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

	responses, pageMetadata, err := c.UseCase.Search(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
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
	parsedUUID, err := uuid.Parse(ctx.Params("certificateCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing Certificate category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid Certificate category uuid")
	}

	request := &model.GetCertificateCatRequest{
		UUID: parsedUUID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
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

	parsedUUID, err := uuid.Parse(ctx.Params("certificateCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing Certificate category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid Certificate category uuid")
	}

	request = &model.UpdateCertificateCatRequest{
		UUID:        parsedUUID,
		Name:        request.Name,
		Description: request.Description,
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error updating Certificate category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateCatResponse]{
		Success: true,
		Data:    response,
	})
}
