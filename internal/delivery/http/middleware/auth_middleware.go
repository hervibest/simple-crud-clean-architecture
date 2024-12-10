package middleware

import (
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUserCase *usecase.UserUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		request := &model.VerifyUserRequest{Token: token}
		auth, err := userUserCase.Verify(ctx.UserContext(), request)
		if err != nil {
			return err
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}
