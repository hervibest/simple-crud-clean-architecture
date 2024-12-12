package middleware

import (
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewAuth(userUseCase *usecase.UserUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		request := &model.VerifyUserRequest{Token: token}
		auth, err := userUseCase.Verify(ctx.UserContext(), request)
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

func BuyableCourse(courseUseCase *usecase.CourseUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		request := new(model.GetCourseRequest)
		parsedUUID, err := uuid.Parse(ctx.Params("courseID"))
		if err != nil {
			return fiber.ErrBadRequest
		}

		request.UUID = parsedUUID
		_, err = courseUseCase.Get(ctx.UserContext(), request)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")

		}

		return ctx.Next()
	}
}
