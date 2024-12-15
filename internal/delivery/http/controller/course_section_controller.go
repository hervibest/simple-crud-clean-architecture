package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CourseSectionController struct {
	UseCase   *usecase.CourseSecUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewCourseSecController(useCase *usecase.CourseSecUseCase, log *logrus.Logger, validator helper.CustomValidator) *CourseSectionController {
	return &CourseSectionController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
	}
}

func (c *CourseSectionController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCourseSecRequest)
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

	return ctx.JSON(model.DataResponse[*model.CourseSectionResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CourseSectionController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCourseSecRequest{
		Title:       ctx.Query("title", ""),
		Description: ctx.Query("description", ""),
		Page:        ctx.QueryInt("page", 1),
		Size:        ctx.QueryInt("size", 10),
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
		c.Log.WithError(err).Error("error searching course sections")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CourseSectionResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CourseSectionController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("courseSecID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section video uuid")
		return fiber.ErrBadRequest
	}

	request := &model.GetCourseSecRequest{
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

	return ctx.JSON(model.DataResponse[*model.CourseSectionResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CourseSectionController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCourseSecRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("courseSecID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course section uuid")
		return fiber.ErrBadRequest
	}

	request.CourseSecUUID = parsedUUID

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

	return ctx.JSON(model.DataResponse[*model.CourseSectionResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CourseSectionController) Delete(ctx *fiber.Ctx) error {

	request := new(model.DeleteCourseSecRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("courseSecID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section video uuid")
		return fiber.ErrBadRequest
	}

	request.CourseSecUUID = parsedUUID

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error deleting course category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseSectionResponse]{
		Success: true,
		Data:    response,
	})
}
