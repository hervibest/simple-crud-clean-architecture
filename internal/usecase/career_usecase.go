package usecase

import (
	"context"
	"mime/multipart"
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

type CareerUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	CareerRepository    *repository.CareerRepository
	CareerCatRepository *repository.CareerCategoryRepository
	Minio               *helper.Minio
	Validator           helper.CustomValidator
}

func NewCareerUseCase(db *gorm.DB, logger *logrus.Logger, careerRepository *repository.CareerRepository,
	careerCatRepository *repository.CareerCategoryRepository, minio *helper.Minio, validator helper.CustomValidator) *CareerUseCase {
	return &CareerUseCase{
		DB:                  db,
		Log:                 logger,
		CareerRepository:    careerRepository,
		CareerCatRepository: careerCatRepository,
		Minio:               minio,
		Validator:           validator,
	}
}

func (c *CareerUseCase) Create(ctx context.Context, request *model.CreateCareerRequest) (*model.CareerResponse, error) {

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	total, err := c.CareerRepository.CountByName(tx, request.Title)
	if err != nil {
		c.Log.Warnf("Failed career user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	career := &entity.Career{
		UUID:        uuid.New(),
		Title:       request.Title,
		Slug:        slugHelper.Make(request.Title),
		Price:       request.Price,
		Description: request.Description,
	}

	if err := c.CareerRepository.Create(tx, career); err != nil {
		c.Log.WithError(err).Error("error creating career")
		return nil, fiber.ErrInternalServerError
	}

	if request.CategoryUUIDs == nil {
		if err := c.CareerRepository.ClearCategory(tx, career); err != nil {
			c.Log.WithError(err).Error("error clearing career categories")
			return nil, fiber.ErrInternalServerError
		}
	} else if request.CategoryUUIDs != nil {
		careerCategories, err := c.CareerCatRepository.FindManyByUUIDs(tx, request.CategoryUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding career categories")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.CareerRepository.SyncCategory(tx, career, careerCategories); err != nil {
			c.Log.WithError(err).Error("error syncing career categories")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating career")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CareerToResponse(career), nil
}

func (c *CareerUseCase) Search(ctx context.Context, request *model.SearchCareerRequest) ([]model.CareerResponse, *model.PageMetadata, error) {

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, nil, validationErr
	}

	careers, pageMetadata, err := c.CareerRepository.Search(c.DB, request, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting career category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CareerResponse, len(careers))
	for i, career := range careers {
		responses[i] = *converter.CareerToResponse(&career)
	}

	return responses, pageMetadata, nil
}

func (c *CareerUseCase) Get(ctx context.Context, request *model.GetCareerRequest) (*model.CareerResponse, error) {
	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	career, err := c.CareerRepository.FindWithDetails(c.DB, request.UUID, true, true, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting career")
		return nil, fiber.ErrNotFound
	}

	return converter.CareerToResponse(career), nil
}

func (c *CareerUseCase) Update(ctx context.Context, request *model.UpdateCareerRequest) (*model.CareerResponse, error) {

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	career := new(entity.Career)
	if err := c.CareerRepository.FindByUUID(tx, career, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	total, err := c.CareerRepository.CountByNameAndNotID(tx, request.Title, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get career user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	career.Title = request.Title
	career.Slug = slugHelper.Make(request.Title)
	career.Description = request.Description
	career.Price = request.Price
	career.IsActive = request.IsActive

	if err := c.CareerRepository.Update(tx, career); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if request.CategoryUUIDs == nil {
		if err := c.CareerRepository.ClearCategory(tx, career); err != nil {
			c.Log.WithError(err).Error("error clearing career categories")
			return nil, fiber.ErrInternalServerError
		}
	} else if len(request.CategoryUUIDs) > 0 {
		careerCategories, err := c.CareerCatRepository.FindManyByUUIDs(tx, request.CategoryUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding career categories UUIDs")
			return nil, fiber.NewError(fiber.StatusBadRequest, "career categories UUIDs not valid")
		}

		if err := c.CareerRepository.ClearCategory(tx, career); err != nil {
			c.Log.WithError(err).Error("error clearing career categories")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.CareerRepository.SyncCategory(tx, career, careerCategories); err != nil {
			c.Log.WithError(err).Error("error syncing career categories")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("Career updated : %+v", career)

	return converter.CareerToResponse(career), nil
}

func (c *CareerUseCase) UploadThumbnail(ctx context.Context, file *multipart.FileHeader, request *model.CareerThumbnailRequest) (*model.CareerResponse, error) {
	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	career := new(entity.Career)
	if err := c.CareerRepository.FindByUUID(tx, career, request.CareerUUID); err != nil {
		c.Log.Warnf("Failed find career from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	upload, err := c.Minio.UploadFileToMinio(ctx, file, "thumbnail")
	if err != nil {
		return nil, err
	}

	newThumbnailFile := &entity.File{

		UUID:     uuid.New(),
		Filename: upload.Filename,
		Mimetype: upload.Mimetype,
		Path:     upload.URL,
		Size:     upload.Size,
	}

	if err := c.CareerRepository.AttachThumbnail(tx, career, newThumbnailFile); err != nil {
		c.Log.Warnf("Failed to attach video")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to attach thumbnail"+err.Error())
	}
	// OriginalName string    `gorm:"column:original_name"`
	// OriginalSize float64   `gorm:"column:original_size"`
	// OriginalMime string    `gorm:"column:original_mime"`
	// MediaID      string    `gorm:"column:media_id"`
	// Bucket       string    `gorm:"column:bucket"`
	// Dir          string    `gorm:"column:dir"`

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error uploading thumbnail")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CareerToResponse(career), nil

}
