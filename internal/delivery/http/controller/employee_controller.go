package http

import (
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeController struct {
	Log     *logrus.Logger
	UseCase *usecase.EmployeeUseCase
}

func NewEmployeeController(useCase *usecase.EmployeeUseCase, logger *logrus.Logger) *EmployeeController {
	return &EmployeeController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *EmployeeController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterEmployeeRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to register employee : %+v", err)
		return err
	}

	data := map[string]interface{}{
		"employee": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *EmployeeController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginEmployeeRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to login employee : %+v", err)
		return err
	}

	data := map[string]interface{}{
		"employee": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *EmployeeController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetEmployee(ctx)

	request := &model.GetEmployeeRequest{
		Email: auth.Email,
	}

	response, err := c.UseCase.Current(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Warnf("Failed to get current employee")
		return err
	}

	data := map[string]interface{}{
		"employee": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *EmployeeController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetEmployee(ctx)

	request := new(model.LogoutEmployeeRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.Email = auth.Email
	request.AccessToken = auth.Token

	response, err := c.UseCase.Logout(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Warnf("Failed to logout employee")
		return err
	}

	data := map[string]interface{}{
		"logout": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}
