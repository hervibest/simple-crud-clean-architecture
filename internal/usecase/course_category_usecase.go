package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	slugHelper "github.com/gosimple/slug"
)

type CourseCatUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	CourseCatRepository *repository.CourseCategoryRepository
	Validator           helper.CustomValidator
}

func NewCourseCatUseCase(db *gorm.DB, logger *logrus.Logger, courseCatRepository *repository.CourseCategoryRepository,
	validator helper.CustomValidator) *CourseCatUseCase {
	return &CourseCatUseCase{
		DB:                  db,
		Log:                 logger,
		CourseCatRepository: courseCatRepository,
		Validator:           validator,
	}
}

func (c *CourseCatUseCase) Create(ctx context.Context, request *model.CreateCourseCatRequest) (*model.CourseCatResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

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

func (c *CourseCatUseCase) Search(ctx context.Context, request *model.SearchCourseCatRequest) ([]model.CourseCatResponse, *model.PageMetadata, error) {

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, nil, validationErr
	}

	courseCats, pageMetadata, err := c.CourseCatRepository.Search(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting course category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CourseCatResponse, len(courseCats))
	for i, courseCat := range courseCats {
		responses[i] = *converter.CourseCatToResponse(&courseCat)
	}

	return responses, pageMetadata, nil
}

func (c *CourseCatUseCase) Get(ctx context.Context, request *model.GetCourseCatRequest) (*model.CourseCatResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	courseCategory := new(entity.CourseCategory)
	if err := c.CourseCatRepository.FindByUUID(c.DB, courseCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting courseCategory")
		return nil, fiber.ErrNotFound
	}

	return converter.CourseCatToResponse(courseCategory), nil
}

func (c *CourseCatUseCase) Update(ctx context.Context, request *model.UpdateCourseCatRequest) (*model.CourseCatResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	courseCategory := new(entity.CourseCategory)
	if err := c.CourseCatRepository.FindByUUID(tx, courseCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
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
