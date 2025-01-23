package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CareerController struct {
	UseCase *usecase.CareerUseCase
	Log     *logrus.Logger
}

func NewCareerController(useCase *usecase.CareerUseCase, log *logrus.Logger) *CareerController {
	return &CareerController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *CareerController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateCareerRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	// @TODO : Refactor uuid parser into reusable component in helper/parser
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
		c.Log.WithError(err).Error("error when creating career")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CareerResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CareerController) List(ctx *fiber.Ctx) error {

	request := &model.SearchCareerRequest{
		Title:       ctx.Query("title", ""),
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
		c.Log.WithError(err).Error("error when searching careers")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.CareerResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *CareerController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("careerId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing career uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid career uuid")
	}

	request := &model.GetCareerRequest{
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
		c.Log.WithError(err).Error("error when getting career")
		return err
	}
	return ctx.JSON(model.DataResponse[*model.CareerResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CareerController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateCareerRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("careerId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing career uuid")
		return fiber.NewError(fiber.StatusBadRequest, "invalid career uuid")
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
		c.Log.WithError(err).Error("error when updating career")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.CareerResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *CareerController) UploadThumbnail(ctx *fiber.Ctx) error {

	parsedCareerUUID, err := uuid.Parse(ctx.Params("careerId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing career uuid")
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

	request := &model.CareerThumbnailRequest{
		CareerUUID: parsedCareerUUID,
	}

	response, err := c.UseCase.UploadThumbnail(ctx.UserContext(), file, request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error when uploading career thumbnail")
		return err
	}
	return ctx.JSON(model.DataResponse[*model.CareerResponse]{
		Success: true,
		Data:    response,
	})
}
