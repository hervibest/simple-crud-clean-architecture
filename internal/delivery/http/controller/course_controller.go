package http

import (
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CourseController struct {
	UseCase *usecase.CourseUseCase
	Log     *logrus.Logger
}

func NewCourseController(useCase *usecase.CourseUseCase, log *logrus.Logger) *CourseController {
	return &CourseController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CourseController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCourseRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	var categoryUUIDs []uuid.UUID
	for _, category := range request.Categories {
		id, err := uuid.Parse(category)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid category UUID")
		}
		categoryUUIDs = append(categoryUUIDs, id)
	}
	request.CategoryUUIDs = categoryUUIDs

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating course")
		return err
	}

	data := map[string]interface{}{
		"course": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *CourseController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("courseID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing uuid")
		return fiber.ErrBadRequest
	}

	request := &model.GetCourseRequest{
		UUID: parsedUUID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return err
	}

	data := map[string]interface{}{
		"course": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *CourseController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCourseRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("courseID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing uuid")
		return fiber.ErrBadRequest
	}

	request.UUID = parsedUUID

	var categoryUUIDs []uuid.UUID
	for _, category := range request.Categories {
		id, err := uuid.Parse(category)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid category UUID")
		}
		categoryUUIDs = append(categoryUUIDs, id)
	}
	request.CategoryUUIDs = categoryUUIDs

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return err
	}

	data := map[string]interface{}{
		"contact": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}
