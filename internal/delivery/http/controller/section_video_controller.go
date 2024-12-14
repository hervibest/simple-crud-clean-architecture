package http

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SectionVideoController struct {
	UseCase *usecase.SecVideoUseCase
	Log     *logrus.Logger
}

func NewSecVideoController(useCase *usecase.SecVideoUseCase, log *logrus.Logger) *SectionVideoController {
	return &SectionVideoController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *SectionVideoController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateSecVideoRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating course category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) List(ctx *fiber.Ctx) error {

	request := &model.SearchSecVideosRequest{
		Title: ctx.Query("title", ""),
		Notes: ctx.Query("notes", ""),
		Page:  ctx.QueryInt("page", 1),
		Size:  ctx.QueryInt("size", 10),
	}

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	responses, pageMetadata, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching course sections")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.SecVideoResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *SectionVideoController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("secVideoID"))

	c.Log.Infof(parsedUUID.String())

	if err != nil {
		c.Log.WithError(err).Error("error parsing uuid course controller")
		return fiber.ErrBadRequest
	}

	request := &model.GetSecVideoRequest{
		UUID: parsedUUID,
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateSecVideoRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("secVideoID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing uuid")
		return fiber.ErrBadRequest
	}

	request.VideoUUID = parsedUUID

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating course category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) Delete(ctx *fiber.Ctx) error {

	request := new(model.DeleteSecVideoRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("secVideoID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing uuid")
		return fiber.ErrBadRequest
	}

	request.VideoUUID = parsedUUID

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error deleting course category")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}
