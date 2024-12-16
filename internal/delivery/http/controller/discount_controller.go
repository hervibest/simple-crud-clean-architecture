package http

import (
	"simple-crud-clean-architecture/internal/enum"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type DiscountController struct {
	UseCase   *usecase.DiscountUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewDiscountController(useCase *usecase.DiscountUseCase, log *logrus.Logger, validator helper.CustomValidator) *DiscountController {
	return &DiscountController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
	}
}

func (c *DiscountController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateDiscountRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	if len(request.Courses) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "course_uuid should not empty")
	}

	var courseUUIDs []uuid.UUID

	for _, course := range request.Courses {
		id, err := uuid.Parse(course)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid course UUID")
		}
		courseUUIDs = append(courseUUIDs, id)
	}
	request.CourseUUIDs = courseUUIDs

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	validateRequest := &model.ValidateDiscountRequest{
		UUID:          uuid.Nil,
		Name:          request.Name,
		Type:          request.Type,
		Value:         request.Value,
		StartActiveAt: request.StartActiveAt,
		ValidUntil:    request.ValidUntil,
		CourseUUIDs:   courseUUIDs,
	}

	if err := c.UseCase.Validate(ctx.UserContext(), validateRequest); err != nil {
		c.Log.WithError(err).Error("error validating discount")
		return err
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating discount")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.DiscountResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *DiscountController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateDiscountRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("discountID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")
	}

	request.UUID = parsedUUID

	var courseUUIDs []uuid.UUID
	for _, course := range request.Courses {
		id, err := uuid.Parse(course)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid course UUID")
		}
		courseUUIDs = append(courseUUIDs, id)
	}

	request.CourseUUIDs = courseUUIDs

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	validateRequest := &model.ValidateDiscountRequest{
		UUID:          request.UUID,
		Name:          request.Name,
		Type:          request.Type,
		Value:         request.Value,
		StartActiveAt: request.StartActiveAt,
		ValidUntil:    request.ValidUntil,
		CourseUUIDs:   request.CourseUUIDs,
	}

	if err := c.UseCase.Validate(ctx.UserContext(), validateRequest); err != nil {
		c.Log.WithError(err).Error("error validating discount")
		return err
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating discount")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.DiscountResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *DiscountController) List(ctx *fiber.Ctx) error {

	request := &model.SearchDiscountRequest{
		Name: ctx.Query("name", ""),
		Type: enum.DiscountType(ctx.Query("type", "")),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
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
		c.Log.WithError(err).Error("error searching course")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.DiscountResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *DiscountController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("discountID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course uuids")
		return fiber.NewError(fiber.StatusBadRequest, "invalid discount uuid")
	}

	request := &model.GetDiscountRequest{
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
		c.Log.WithError(err).Error("error getting course")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.DiscountResponse]{
		Success: true,
		Data:    response,
	})
}
