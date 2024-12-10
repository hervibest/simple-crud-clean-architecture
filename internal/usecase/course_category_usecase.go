package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	slugHelper "github.com/gosimple/slug"
)

type CourseCatUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	Validate            *validator.Validate
	CourseCatRepository *repository.CourseCategoryRepository
}

func NewCourseCatUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	courseCatRepository *repository.CourseCategoryRepository) *CourseCatUseCase {
	return &CourseCatUseCase{
		DB:                  db,
		Log:                 logger,
		Validate:            validate,
		CourseCatRepository: courseCatRepository,
	}
}

func (c *CourseCatUseCase) Create(ctx context.Context, request *model.CreateCourseCatRequest) (*model.CourseCatResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	helper.SanitiseStruct(request)

	total, err := c.CourseCatRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed course category user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	courseCategory := &entity.CourseCategory{
		UUID:        uuid.New(),
		Name:        request.Name,
		Slug:        slugHelper.Make(request.Name),
		Description: request.Description,
	}

	if err := c.CourseCatRepository.Create(tx, courseCategory); err != nil {
		c.Log.WithError(err).Error("error creating course category")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course category")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseCatToResponse(courseCategory), nil
}

func (c *CourseCatUseCase) Get(ctx context.Context, request *model.GetCourseCatRequest) (*model.CourseCatResponse, error) {

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	courseCategory := new(entity.CourseCategory)
	if err := c.CourseCatRepository.FindByUUID(c.DB, courseCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting courseCategory")
		return nil, fiber.ErrNotFound
	}

	return converter.CourseCatToResponse(courseCategory), nil
}

func (c *CourseCatUseCase) Update(ctx context.Context, request *model.UpdateCourseCatRequest) (*model.CourseCatResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	courseCategory := new(entity.CourseCategory)
	if err := c.CourseCatRepository.FindByUUID(tx, courseCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	total, err := c.CourseCatRepository.CountByNameAndNotID(tx, request.Name, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get course category user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	courseCategory.Name = request.Name
	courseCategory.Slug = slugHelper.Make(request.Name)
	courseCategory.Description = request.Description

	if err := c.CourseCatRepository.Update(tx, courseCategory); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseCatToResponse(courseCategory), nil
}
