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

type CertifCatUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	CertifCatRepository *repository.CertifCategoryRepository
}

func NewCertifCatUseCase(db *gorm.DB, logger *logrus.Logger, certifCatRepository *repository.CertifCategoryRepository) *CertifCatUseCase {
	return &CertifCatUseCase{
		DB:                  db,
		Log:                 logger,
		CertifCatRepository: certifCatRepository,
	}
}

func (c *CertifCatUseCase) Create(ctx context.Context, request *model.CreateCertifCatRequest) (*model.CertificateCatResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.CertifCatRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed certif category user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	certifCategory := &entity.CertificateCategory{
		UUID:        uuid.New(),
		Name:        request.Name,
		Slug:        slugHelper.Make(request.Name),
		Description: request.Description,
	}

	if err := c.CertifCatRepository.Create(tx, certifCategory); err != nil {
		c.Log.WithError(err).Error("error creating certif category")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating certif category")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertificateCatToResponse(certifCategory), nil
}

func (c *CertifCatUseCase) Search(ctx context.Context, request *model.SearchCertificateCatRequest) ([]model.CertificateCatResponse, *model.PageMetadata, error) {

	certifCats, pageMetadata, err := c.CertifCatRepository.Search(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting certif category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CertificateCatResponse, len(certifCats))
	for i, certifCat := range certifCats {
		responses[i] = *converter.CertificateCatToResponse(&certifCat)
	}

	return responses, pageMetadata, nil
}

func (c *CertifCatUseCase) Get(ctx context.Context, request *model.GetCertificateCatRequest) (*model.CertificateCatResponse, error) {

	certifCategory := new(entity.CertificateCategory)
	if err := c.CertifCatRepository.FindByUUID(c.DB, certifCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting certifCategory")
		return nil, fiber.ErrNotFound
	}

	return converter.CertificateCatToResponse(certifCategory), nil
}

func (c *CertifCatUseCase) Update(ctx context.Context, request *model.UpdateCertificateCatRequest) (*model.CertificateCatResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	certifCategory := new(entity.CertificateCategory)
	if err := c.CertifCatRepository.FindByUUID(tx, certifCategory, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	total, err := c.CertifCatRepository.CountByNameAndNotID(tx, request.Name, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get certif category user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	certifCategory.Name = request.Name
	certifCategory.Slug = slugHelper.Make(request.Name)
	certifCategory.Description = request.Description

	if err := c.CertifCatRepository.Update(tx, certifCategory); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertificateCatToResponse(certifCategory), nil
}
