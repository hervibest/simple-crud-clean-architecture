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
	UseCase   *usecase.VoucherUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
}

func NewVoucherController(useCase *usecase.VoucherUseCase, log *logrus.Logger, validator helper.CustomValidator) *VoucherController {
	return &VoucherController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
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

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

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

	if err := c.UseCase.Validate(ctx.UserContext(), validateRequest); err != nil {
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

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

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

	if err := c.UseCase.Validate(ctx.UserContext(), validateRequest); err != nil {
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

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.ApplyVoucher(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting voucher")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.VoucherResponse]{
		Success: true,
		Data:    response,
	})
}
