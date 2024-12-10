package http

import (
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
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	err = c.UseCase.RequestEmailVerification(ctx.UserContext(), response.Email, true)
	if err != nil {
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
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
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
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse{
		Success: true,
	})

}
