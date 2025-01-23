package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/enum"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DiscountUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	DiscountRepository *repository.DiscountRepository
	CourseRepository   *repository.CourseRepository
	Validator          helper.CustomValidator
}

func NewDiscountUseCase(db *gorm.DB, logger *logrus.Logger, discountRepository *repository.DiscountRepository,
	courseRepository *repository.CourseRepository, validator helper.CustomValidator) *DiscountUseCase {
	return &DiscountUseCase{
		DB:                 db,
		Log:                logger,
		DiscountRepository: discountRepository,
		CourseRepository:   courseRepository,
		Validator:          validator,
	}
}

func (c *DiscountUseCase) Create(ctx context.Context, request *model.CreateDiscountRequest) (*model.DiscountResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.DiscountRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed discount user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	discount := &entity.Discount{
		UUID:          uuid.New(),
		Name:          request.Name,
		Value:         request.Value,
		Type:          request.Type,
		StartActiveAt: request.StartActiveAt,
		ValidUntil:    request.ValidUntil,
	}

	if err := c.DiscountRepository.Create(tx, discount); err != nil {
		c.Log.WithError(err).Error("error creating discount")
		return nil, fiber.ErrInternalServerError
	}

	if request.CourseUUIDs == nil {
		if err := c.DiscountRepository.ClearCourse(tx, discount); err != nil {
			c.Log.WithError(err).Error("error clearing course")
			return nil, fiber.ErrInternalServerError
		}
	} else if request.CourseUUIDs != nil {
		courses, err := c.CourseRepository.FindManyByUUIDs(tx, request.CourseUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding course")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.DiscountRepository.SyncCourse(tx, discount, courses); err != nil {
			c.Log.WithError(err).Error("error syncing course")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating discount")
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountToResponse(discount), nil
}

func (c *DiscountUseCase) Validate(ctx context.Context, request *model.ValidateDiscountRequest) error {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return validationErr
	}

	helper.SanitiseStruct(request)

	if request.UUID == uuid.Nil {
		return fiber.NewError(fiber.StatusBadGateway, "uuid should not be empty")
	}

	discount := new(entity.Discount)
	if err := c.DiscountRepository.FindByUUID(c.DB, discount, request.UUID); err != nil {
		c.Log.Warnf("Failed discount user from database : %+v", err)
		return fiber.ErrInternalServerError
	}

	if discount.IsActive {
		return fiber.NewError(fiber.StatusBadRequest, "current valid and active discount cannot be update")
	}

	courses, err := c.CourseRepository.FindManyByUUIDs(c.DB, request.CourseUUIDs)
	if err != nil {
		c.Log.WithError(err).Error("error finding course")
		return fiber.ErrInternalServerError
	}
	c.Log.Infof("test debug courses" + courses[0].Name)

	for _, course := range courses {
		coursePriceToCompare := course.Price
		if request.Type == enum.DiscountTypeRebate && coursePriceToCompare < request.Value {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "requested Discount will cause a minus price on course")
		}
		c.Log.Infof("ini adalah perhitungan validasi" + string(request.Type) + string(enum.DiscountTypePercent))

		if request.Type == enum.DiscountTypePercent {
			coursePriceToCompare = course.Price - request.Value/100*course.Price
			c.Log.Infof("ini adalah perhitungan validasi" + strconv.FormatFloat(coursePriceToCompare, 'f', 2, 64))
			if coursePriceToCompare <= 0 {
				return fiber.NewError(fiber.StatusUnprocessableEntity, "requested Discount will cause a minus price on course")

			}
		}
	}

	return nil
}

func (c *DiscountUseCase) Update(ctx context.Context, request *model.UpdateDiscountRequest) (*model.DiscountResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	discount := new(entity.Discount)
	if err := c.DiscountRepository.FindByUUID(tx, discount, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	total, err := c.DiscountRepository.CountByNameAndNotID(tx, request.Name, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get discount user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	discount.Name = request.Name
	discount.Type = request.Type
	discount.StartActiveAt = request.StartActiveAt
	discount.ValidUntil = request.ValidUntil
	discount.Value = request.Value

	if err := c.DiscountRepository.Update(tx, discount); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if request.CourseUUIDs == nil {
		if err := c.DiscountRepository.ClearCourse(tx, discount); err != nil {
			c.Log.WithError(err).Error("error clearing discount categories")
			return nil, fiber.ErrInternalServerError
		}
	} else if len(request.CourseUUIDs) > 0 {
		courses, err := c.CourseRepository.FindManyByUUIDs(tx, request.CourseUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding discount categories UUIDs")
			return nil, fiber.NewError(fiber.StatusBadRequest, "discount categories UUIDs not valid")
		}

		if err := c.DiscountRepository.ClearCourse(tx, discount); err != nil {
			c.Log.WithError(err).Error("error clearing discount categories")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.DiscountRepository.SyncCourse(tx, discount, courses); err != nil {
			c.Log.WithError(err).Error("error syncing discount categories")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("Discount updated : %+v", discount)

	return converter.DiscountToResponse(discount), nil
}

func (c *DiscountUseCase) Search(ctx context.Context, request *model.SearchDiscountRequest) ([]model.DiscountResponse, *model.PageMetadata, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, nil, validationErr
	}

	courses, pageMetadata, err := c.DiscountRepository.Search(c.DB, request, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting discount")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.DiscountResponse, len(courses))
	for i, course := range courses {
		responses[i] = *converter.DiscountToResponse(&course)
	}

	return responses, pageMetadata, nil
}

func (c *DiscountUseCase) Get(ctx context.Context, request *model.GetDiscountRequest) (*model.DiscountResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	course, err := c.DiscountRepository.FindWithDetails(c.DB, request.UUID, true, false)
	if err != nil {
		c.Log.WithError(err).Error("error getting course")
		return nil, fiber.ErrNotFound
	}

	return converter.DiscountToResponse(course), nil
}

func (c *DiscountUseCase) ActivateDiscount(ctx context.Context) error {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.DiscountRepository.ActivateDiscount(tx); err != nil {
		c.Log.Warnf("Error activating discount using cron" + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Error comitting tx activate discount using cron" + err.Error())
		return err
	}

	c.Log.Warnf("Activate discount using Job")

	return nil
}

func (c *DiscountUseCase) DeactivateDiscount(ctx context.Context) error {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.DiscountRepository.DeactivateDiscount(tx); err != nil {
		c.Log.Warnf("Error deactivating discount using cron" + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Error commiting tx deactivate discount using cron" + err.Error())
		return err
	}

	c.Log.Warnf("Deactivate discount using Job")

	return nil
}
