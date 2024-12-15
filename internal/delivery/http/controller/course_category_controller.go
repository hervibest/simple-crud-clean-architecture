package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CourseCatController struct {
	UseCase   *usecase.CourseCatUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewCourseCatController(useCase *usecase.CourseCatUseCase, log *logrus.Logger, validator helper.CustomValidator) *CourseCatController {
	return &CourseCatController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
	}
}

func (c *CourseCatController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCourseCatRequest)
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
		c.Log.WithError(err).Error("error creating course category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseCatResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CourseCatController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCourseCatRequest{
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
		c.Log.WithError(err).Error("error searching course category")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CourseCatResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CourseCatController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("courseCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid course category uuid")
	}

	request := &model.GetCourseCatRequest{
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
		c.Log.WithError(err).Error("error getting course category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseCatResponse]{
		Success: true,
		Data:    response,
	})

}

func (c *CourseCatController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCourseCatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("courseCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course category uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid course category uuid")
	}

	request = &model.UpdateCourseCatRequest{
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
		c.Log.WithError(err).Error("error updating course category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseCatResponse]{
		Success: true,
		Data:    response,
	})
}
