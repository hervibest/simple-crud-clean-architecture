package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CertificateController struct {
	UseCase *usecase.CertificateUseCase
	Log     *logrus.Logger
}

func NewCertificateController(useCase *usecase.CertificateUseCase, log *logrus.Logger) *CertificateController {
	return &CertificateController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CertificateController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateCertificateRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	id, err := uuid.Parse(request.Category)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid category UUID")
	}
	request.CategoryUUID = id

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertificateController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCertificateRequest{
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
		c.Log.WithError(err).Error("error searching certificate")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CertificateResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CertificateController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("certificateId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing certificate uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid certificate uuid")
	}

	request := &model.GetCertificateRequest{
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
		c.Log.WithError(err).Error("error getting certificate")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertificateController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCertificateRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("certificateId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing certificate uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid certificate uuid")
	}

	request.UUID = parsedUUID

	id, err := uuid.Parse(request.Category)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid category UUID")
	}

	request.CategoryUUID = id

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error updating certificate")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertificateController) UploadThumbnail(ctx *fiber.Ctx) error {

	parsedCertificateUUID, err := uuid.Parse(ctx.Params("certificateId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing certificate uuid")
		return fiber.ErrBadRequest
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 2 * 1024 * 1024 // 10MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	request := &model.CertificateThumbnailRequest{
		CertificateUUID: parsedCertificateUUID,
	}

	response, err := c.UseCase.UploadThumbnail(ctx.UserContext(), file, request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error updating section video")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateResponse]{
		Success: true,
		Data:    response,
	})
}
