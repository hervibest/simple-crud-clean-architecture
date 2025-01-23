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
	UseCase *usecase.CareerCatUseCase
	Log     *logrus.Logger
}

func NewCareerCatController(useCase *usecase.CareerCatUseCase, log *logrus.Logger) *CareerCatController {
	return &CareerCatController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CareerCatController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCareerCatRequest)
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
		c.Log.WithError(err).Error("error when creating career category")
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

	responses, pageMetadata, err := c.UseCase.Search(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error when searching career categories")
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

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error when getting career category")
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

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error when updating career category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CareerCatResponse]{
		Success: true,
		Data:    response,
	})
}
