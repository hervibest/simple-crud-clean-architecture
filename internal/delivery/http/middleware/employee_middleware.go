package middleware

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewEmployeeAuth(employeeUseCase *usecase.EmployeeUseCase, validator helper.CustomValidator) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		request := &model.VerifyEmployeeRequest{Token: token}
		if validationErr := validator.ValidateUseCase(request); validationErr != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
				Success: false,
				Errors:  validationErr,
				Message: "validation error",
			})
		}

		auth, err := employeeUseCase.Verify(ctx.UserContext(), request)
		if err != nil {
			return err
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetEmployee(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
