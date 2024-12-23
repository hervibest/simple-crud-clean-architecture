package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CertifMaterialController struct {
	UseCase   *usecase.CertifMaterialUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewCertifMaterialController(useCase *usecase.CertifMaterialUseCase, log *logrus.Logger, validator helper.CustomValidator) *CertifMaterialController {
	return &CertifMaterialController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
	}
}

func (c *CertifMaterialController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCertificateMatRequest)
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
		c.Log.WithError(err).Error("error creating material category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateMaterialResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertifMaterialController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCertificateMatRequest{
		Name: ctx.Query("name", ""),
		Code: ctx.Query("code", ""),
		Type: ctx.Query("type", ""),
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
		c.Log.WithError(err).Error("error searching material sections")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CertificateMaterialResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CertifMaterialController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("materialID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section skkni uuid")
		return fiber.ErrBadRequest
	}

	request := &model.GetCertificateMatRequest{
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

	return ctx.JSON(model.DataResponse[*model.CertificateMaterialResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertifMaterialController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCertificateMatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("materialID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing material section uuid")
		return fiber.ErrBadRequest
	}

	request.CertificateMatUUID = parsedUUID

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating material category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateMaterialResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CertifMaterialController) Delete(ctx *fiber.Ctx) error {

	request := new(model.DeleteCertificateMatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("materialID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section skkni uuid")
		return fiber.ErrBadRequest
	}

	request.CertificateMatUUID = parsedUUID

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error deleting material category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateMaterialResponse]{
		Success: true,
		Data:    response,
	})
}
