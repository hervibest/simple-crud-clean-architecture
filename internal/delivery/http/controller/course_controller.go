package http

import (
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/helper"
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
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error creating course")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CourseController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCourseRequest{
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
		c.Log.WithError(err).Error("error searching course")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CourseResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CourseController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("courseId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")
	}

	request := &model.GetCourseRequest{
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
		c.Log.WithError(err).Error("error getting course")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CourseController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCourseRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("courseId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")
	}

	request.UUID = parsedUUID

	var categoryUUIDs []uuid.UUID
	for _, category := range request.Categories {
		id, err := uuid.Parse(category)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid category UUID")
		}
		categoryUUIDs = append(categoryUUIDs, id)
	}
	request.CategoryUUIDs = categoryUUIDs

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error updating course")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CourseController) ListUserPurchased(ctx *fiber.Ctx) error {

	request := &model.SearchCourseRequest{
		Name:        ctx.Query("name", ""),
		Slug:        ctx.Query("slug", ""),
		Description: ctx.Query("description", ""),
		Page:        ctx.QueryInt("page", 1),
		Size:        ctx.QueryInt("size", 10),
	}

	auth := middleware.GetUser(ctx)
	userId := auth.Id

	responses, pageMetadata, err := c.UseCase.UserGetPurchasedCourse(ctx.UserContext(), request, userId)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error searching course")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CourseResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CourseController) UploadThumbnail(ctx *fiber.Ctx) error {

	parsedCourseUUID, err := uuid.Parse(ctx.Params("courseId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course uuid")
		return fiber.ErrBadRequest
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 2 * 1024 * 1024 // 10MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	request := &model.CourseThumbnailRequest{
		CourseUUID: parsedCourseUUID,
	}

	response, err := c.UseCase.UploadThumbnail(ctx.UserContext(), file, request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error updating section video")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CourseResponse]{
		Success: true,
		Data:    response,
	})
}
