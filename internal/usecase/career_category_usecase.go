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

type CareerCatUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	CareerCatRepository *repository.CareerCategoryRepository
	Validator           helper.CustomValidator
}

func NewCareerCatUseCase(db *gorm.DB, logger *logrus.Logger, careerCatRepository *repository.CareerCategoryRepository, validator helper.CustomValidator) *CareerCatUseCase {
	return &CareerCatUseCase{
		DB:                  db,
		Log:                 logger,
		CareerCatRepository: careerCatRepository,
		Validator:           validator,
	}
}

func (c *CareerCatUseCase) Create(ctx context.Context, request *model.CreateCareerCatRequest) (*model.CareerCatResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	total, err := c.CareerCatRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed career category user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	careerCategory := &entity.CareerCategory{
		UUID:        uuid.New(),
		Name:        request.Name,
		Slug:        slugHelper.Make(request.Name),
		Description: request.Description,
	}

	if err := c.CareerCatRepository.Create(tx, careerCategory); err != nil {
		c.Log.WithError(err).Error("error creating career category")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating career category")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CareerCatToResponse(careerCategory), nil
}

func (c *CareerCatUseCase) Search(ctx context.Context, request *model.SearchCareerCatRequest) ([]model.CareerCatResponse, *model.PageMetadata, error) {

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, nil, validationErr
	}

	careerCats, pageMetadata, err := c.CareerCatRepository.Search(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting career category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CareerCatResponse, len(careerCats))
	for i, careerCat := range careerCats {
		responses[i] = *converter.CareerCatToResponse(&careerCat)
	}

	return responses, pageMetadata, nil
}

func (c *CareerCatUseCase) Get(ctx context.Context, request *model.GetCareerCatRequest) (*model.CareerCatResponse, error) {

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	careerCategory := new(entity.CareerCategory)
	if err := c.CareerCatRepository.FindByUUID(c.DB, careerCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting careerCategory")
		return nil, fiber.ErrNotFound
	}

	return converter.CareerCatToResponse(careerCategory), nil
}

func (c *CareerCatUseCase) Update(ctx context.Context, request *model.UpdateCareerCatRequest) (*model.CareerCatResponse, error) {

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	careerCategory := new(entity.CareerCategory)
	if err := c.CareerCatRepository.FindByUUID(tx, careerCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	total, err := c.CareerCatRepository.CountByNameAndNotID(tx, request.Name, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get career category user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	careerCategory.Name = request.Name
	careerCategory.Slug = slugHelper.Make(request.Name)
	careerCategory.Description = request.Description

	if err := c.CareerCatRepository.Update(tx, careerCategory); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CareerCatToResponse(careerCategory), nil
}
