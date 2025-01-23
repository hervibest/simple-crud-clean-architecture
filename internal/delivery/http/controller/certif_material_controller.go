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
	UseCase *usecase.CertifMaterialUseCase
	Log     *logrus.Logger
}

func NewCertifMaterialController(useCase *usecase.CertifMaterialUseCase, log *logrus.Logger) *CertifMaterialController {
	return &CertifMaterialController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CertifMaterialController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCertificateMatRequest)
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

	responses, pageMetadata, err := c.UseCase.Search(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
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

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
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

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
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

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error deleting material category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CertificateMaterialResponse]{
		Success: true,
		Data:    response,
	})
}
