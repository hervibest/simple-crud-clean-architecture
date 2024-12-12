package http

import (
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CourseSectionController struct {
	UseCase *usecase.CourseSecUseCase
	Log     *logrus.Logger
}

func NewCourseSecController(useCase *usecase.CourseSecUseCase, log *logrus.Logger) *CourseSectionController {
	return &CourseSectionController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CourseSectionController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCourseSecRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
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
