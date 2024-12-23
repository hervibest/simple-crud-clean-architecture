package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CareerCatController struct {
	UseCase   *usecase.CareerCatUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewCareerCatController(useCase *usecase.CareerCatUseCase, log *logrus.Logger, validator helper.CustomValidator) *CareerCatController {
	return &CareerCatController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
	}
}

func (c *CareerCatController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCareerCatRequest)
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
		c.Log.WithError(err).Error("error creating career category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CareerCatResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CareerCatController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCareerCatRequest{
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
		c.Log.WithError(err).Error("error searching career category")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CareerCatResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CareerCatController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("careerCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing career category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid career category uuid")
	}

	request := &model.GetCareerCatRequest{
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
		c.Log.WithError(err).Error("error getting career category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CareerCatResponse]{
		Success: true,
		Data:    response,
	})

}

func (c *CareerCatController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCareerCatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("careerCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing career category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid career category uuid")
	}

	request = &model.UpdateCareerCatRequest{
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
		c.Log.WithError(err).Error("error updating career category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CareerCatResponse]{
		Success: true,
		Data:    response,
	})
}
