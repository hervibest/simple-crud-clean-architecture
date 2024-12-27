package middleware

import (
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewUserAuth(userUseCase *usecase.UserUseCase, validator helper.CustomValidator) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		request := &model.VerifyUserRequest{Token: token}
		if validationErr := validator.Validate(request); validationErr != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
				Success: false,
				Errors:  validationErr,
				Message: "validation error",
			})
		}

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

func NewBuyableCourse(courseUseCase *usecase.CourseUseCase, validator helper.CustomValidator) fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		request := new(model.CreateTransactionRequest)
		if err := ctx.BodyParser(request); err != nil {
			return fiber.NewError(fiber.StatusNotFound, "invalid uuid")
		}

		courseRequest := new(model.GetCourseRequest)
		parsedUUID, err := uuid.Parse(request.CourseUUIDStr)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "invalid parsing uuid")
		}

		courseRequest.UUID = parsedUUID
		_, err = courseUseCase.Get(ctx.UserContext(), courseRequest)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "invalid course uuid")
		}

		auth := GetUser(ctx)

		purchasedRequest := new(model.GetPurchasedCourseRequest)
		purchasedRequest.UserID = auth.Id

		if validationErr := validator.Validate(request); validationErr != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
				Success: false,
				Errors:  validationErr,
				Message: "validation error",
			})
		}

		courses, _ := courseUseCase.GetPurchasedCourseUUID(ctx.UserContext(), purchasedRequest)
		for _, course := range courses {
			if course.UUID == parsedUUID {
				return fiber.NewError(fiber.StatusBadRequest, "course already purchased")
			}

		}
		return ctx.Next()
	}
}
