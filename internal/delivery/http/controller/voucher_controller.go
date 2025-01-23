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

type VoucherController struct {
	UseCase *usecase.VoucherUseCase
	Log     *logrus.Logger
}

func NewVoucherController(useCase *usecase.VoucherUseCase, log *logrus.Logger) *VoucherController {
	return &VoucherController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *VoucherController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateVoucherRequest)
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

	validateRequest := &model.ValidateVoucherRequest{
		UUID:          uuid.Nil,
		Name:          request.Name,
		Code:          request.Code,
		Type:          request.Type,
		Value:         request.Value,
		StartActiveAt: request.StartActiveAt,
		ValidUntil:    request.ValidUntil,
		CourseUUIDs:   courseUUIDs,
	}

	err := c.UseCase.Validate(ctx.UserContext(), validateRequest)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error validating voucher")
		return err
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating voucher")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.VoucherResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *VoucherController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateVoucherRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("voucherID"))
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

	validateRequest := &model.ValidateVoucherRequest{
		UUID:          request.UUID,
		Name:          request.Name,
		Code:          request.Code,
		Type:          request.Type,
		Value:         request.Value,
		StartActiveAt: request.StartActiveAt,
		ValidUntil:    request.ValidUntil,
		CourseUUIDs:   request.CourseUUIDs,
	}

	err = c.UseCase.Validate(ctx.UserContext(), validateRequest)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error validating voucher")
		return err
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating voucher")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.VoucherResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *VoucherController) List(ctx *fiber.Ctx) error {

	request := &model.SearchVoucherRequest{
		Name: ctx.Query("name", ""),
		Code: ctx.Query("code", ""),
		Type: enum.VoucherType(ctx.Query("type", "")),
		Page: ctx.QueryInt("page", 1),
		Size: ctx.QueryInt("size", 10),
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

	return ctx.JSON(model.DataResponse[[]model.VoucherResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *VoucherController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("voucherID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing course uuids")
		return fiber.NewError(fiber.StatusBadRequest, "invalid voucher uuid")
	}

	request := &model.GetVoucherRequest{
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

	return ctx.JSON(model.DataResponse[*model.VoucherResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *VoucherController) ApplyVoucher(ctx *fiber.Ctx) error {

	voucherCode := ctx.Params("voucherCode")
	request := new(model.ApplyVoucherRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.VoucherCode = voucherCode

	response, err := c.UseCase.ApplyVoucher(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Error("error getting voucher")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.VoucherResponse]{
		Success: true,
		Data:    response,
	})
}
