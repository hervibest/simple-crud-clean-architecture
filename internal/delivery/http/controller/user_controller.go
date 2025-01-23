package http

import (
	"simple-crud-clean-architecture/internal/delivery/http/middleware"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	c.Log.Printf("Register accessed")
	request := new(model.RegisterUserRequest)
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
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	data := map[string]interface{}{
		"user": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *UserController) VerifyEmail(ctx *fiber.Ctx) error {
	request := new(model.VerifyEmailUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.Token = ctx.Params("token")

	err = c.UseCase.VerifyEmail(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to verify user : %+v", err)
		return err
	}

	data := map[string]interface{}{}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *UserController) RequestEmailVerification(ctx *fiber.Ctx) error {
	request := new(model.ResendEmailUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	var rateLimit bool = helper.RateLimit("requestResetPass"+request.Email, 1, 60*2)
	if !rateLimit {
		return fiber.NewError(fiber.StatusTooManyRequests, "can't resend verification email more than once in 2 minutes")
	}

	err = c.UseCase.RequestEmailVerification(ctx.UserContext(), request.Email, false)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
	})
}

func (c *UserController) RequestResetPassword(ctx *fiber.Ctx) error {
	request := new(model.SendResetPasswordRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	var rateLimit bool = helper.RateLimit("requestResetPass"+request.Email, 1, 60*2)
	if !rateLimit {
		c.Log.Warnf("Failed to register user : %+v", err)
		return fiber.NewError(fiber.StatusTooManyRequests, "can't resend reset password email more than once in 2 minutes")
	}

	err = c.UseCase.RequestResetPassword(ctx.UserContext(), request.Email, false)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
	})
}

func (c *UserController) ValidateResetToken(ctx *fiber.Ctx) error {
	request := new(model.ValidateResetTokenRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	valid, err := c.UseCase.ValidateResetToken(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to validate reset token : %+v", err)
		return err
	}

	data := map[string]interface{}{
		"valid": valid,
	}
	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *UserController) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ResetPasswordUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}
	request.Token = ctx.Params("token")

	err = c.UseCase.ResetPassword(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
	})
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
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
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	data := map[string]interface{}{
		"user": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *UserController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := &model.GetUserRequest{
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
		c.Log.WithError(err).Warnf("Failed to get current user")
		return err
	}

	data := map[string]interface{}{
		"user": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *UserController) RequestAccessToken(ctx *fiber.Ctx) error {
	request := new(model.AccessTokenRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.AccessTokenRequest(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.Warnf("Failed to login user : %+v", err)
		return err
	}

	data := map[string]interface{}{
		"user": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.LogoutUserRequest)
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
		c.Log.WithError(err).Warnf("Failed to logout user")
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

func (c *UserController) Update(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	request.Email = auth.Email

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if validationErr, ok := err.(*helper.UseCaseValError); ok {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr.GetValidationErrors(),
			Message: "validation error occurred",
		})
	} else if err != nil {
		c.Log.WithError(err).Warnf("Failed to update user")
		return err
	}

	data := map[string]interface{}{
		"user": response,
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    data,
	})
}
