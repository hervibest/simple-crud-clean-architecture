package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CertifSkkniController struct {
	UseCase   *usecase.CertifSkkniUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewCertifSkkniController(useCase *usecase.CertifSkkniUseCase, log *logrus.Logger, validator helper.CustomValidator) *CertifSkkniController {
	return &CertifSkkniController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
	}
}

func (c *CertifSkkniController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCertificateSkkniRequest)
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
		c.Log.WithError(err).Error("error creating skkni category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateSkkniResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertifSkkniController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCertificateSkkniRequest{
		Name: ctx.Query("name", ""),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
	}

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

	responses, pageMetadata, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching skkni sections")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CertificateSkkniResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CertifSkkniController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("skkniId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section skkni uuid")
		return fiber.ErrBadRequest
	}

	request := &model.GetCertificateSkkniRequest{
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
		c.Log.WithError(err).Error("error getting contact")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateSkkniResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertifSkkniController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCertificateSkkniRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("skkniId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing CertifSkkni Repository  uuid")
		return fiber.ErrBadRequest
	}

	request.CertificateSkkniUUID = parsedUUID

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating CertifSkkniRepository category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateSkkniResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertifSkkniController) Delete(ctx *fiber.Ctx) error {

	request := new(model.DeleteCertificateSkkniRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("skkniID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section skkni uuid")
		return fiber.ErrBadRequest
	}

	request.CertificateSkkniUUID = parsedUUID

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error deleting skkni category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateSkkniResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertifSkkniController) UploadThumbnail(ctx *fiber.Ctx) error {

	parsedSkkniUUID, err := uuid.Parse(ctx.Params("skkniId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing skkni uuid")
		return fiber.ErrBadRequest
	}

	file, err := ctx.FormFile("thumbnail")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 2 * 1024 * 1024 // 10MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	request := &model.SkkniThumbnailRequest{
		SkkniUUID: parsedSkkniUUID,
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.UploadFile(ctx.UserContext(), file, request)
	if err != nil {
		c.Log.WithError(err).Error("error updating and uploading skkni file")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateSkkniResponse]{
		Success: true,
		Data:    response,
	})
}
