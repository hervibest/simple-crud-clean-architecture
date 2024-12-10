package http

import (
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CourseCatController struct {
	UseCase *usecase.CourseCatUseCase
	Log     *logrus.Logger
}

func NewCourseCatController(useCase *usecase.CourseCatUseCase, log *logrus.Logger) *CourseCatController {
	return &CourseCatController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CourseCatController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCourseCatRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating course category")
		return err
	}

	data := map[string]interface{}{
		"course_category": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *CourseCatController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("courseCatID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing uuid")
		return fiber.ErrBadRequest
	}

	request := &model.GetCourseCatRequest{
		UUID: parsedUUID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return err
	}

	data := map[string]interface{}{
		"course_category": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
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
		c.Log.WithError(err).Error("error parsing uuid")
		return fiber.ErrBadRequest
	}

	request = &model.UpdateCourseCatRequest{
		UUID:        parsedUUID,
		Name:        request.Name,
		Description: request.Description,
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return err
	}

	// Populate the Data map with multiple fields
	data := map[string]interface{}{
		"contact": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}
