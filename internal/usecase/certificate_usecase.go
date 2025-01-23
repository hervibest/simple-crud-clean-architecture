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

type CertificateUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	CertificateRepository    *repository.CertificateRepository
	CertifCategoryRepository *repository.CertifCategoryRepository
	Minio                    *helper.Minio
	Validator                helper.CustomValidator
}

func NewCertificateUseCase(db *gorm.DB, logger *logrus.Logger, certificateRepository *repository.CertificateRepository,
	certificateCatRepository *repository.CertifCategoryRepository, minio *helper.Minio, validator helper.CustomValidator) *CertificateUseCase {
	return &CertificateUseCase{
		DB:                       db,
		Log:                      logger,
		CertificateRepository:    certificateRepository,
		CertifCategoryRepository: certificateCatRepository,
		Minio:                    minio,
		Validator:                validator,
	}
}

func (c *CertificateUseCase) Create(ctx context.Context, request *model.CreateCertificateRequest) (*model.CertificateResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.CertificateRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed certificate user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	certificate := &entity.Certificate{
		UUID:        uuid.New(),
		Name:        request.Name,
		Slug:        slugHelper.Make(request.Name),
		Price:       request.Price,
		Description: request.Description,
	}

	if request.CategoryUUID != uuid.Nil {
		certifCategory := new(entity.CertificateCategory)
		if err := c.CertifCategoryRepository.FindByUUID(tx, certifCategory, request.CategoryUUID); err != nil {
			c.Log.WithError(err).Error("error finding certificate categories UUIDs")
			return nil, fiber.NewError(fiber.StatusBadRequest, "certificate categories UUIDs not valid")
		}
		certificate.Category = certifCategory
	}

	if err := c.CertificateRepository.Create(tx, certificate); err != nil {
		c.Log.WithError(err).Error("error creating certificates")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating certificates")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertificateToResponse(certificate), nil
}

func (c *CertificateUseCase) Search(ctx context.Context, request *model.SearchCertificateRequest) ([]model.CertificateResponse, *model.PageMetadata, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, nil, validationErr
	}

	certificates, pageMetadata, err := c.CertificateRepository.Search(c.DB, request, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting certificate category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CertificateResponse, len(certificates))
	for i, certificate := range certificates {
		responses[i] = *converter.CertificateToResponse(&certificate)
	}

	return responses, pageMetadata, nil
}

func (c *CertificateUseCase) Get(ctx context.Context, request *model.GetCertificateRequest) (*model.CertificateResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	certificate, err := c.CertificateRepository.FindWithDetails(c.DB, request.UUID, true, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting certificate")
		return nil, fiber.ErrNotFound
	}

	return converter.CertificateToResponse(certificate), nil
}

func (c *CertificateUseCase) Update(ctx context.Context, request *model.UpdateCertificateRequest) (*model.CertificateResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	certificate := new(entity.Certificate)
	if err := c.CertificateRepository.FindByUUID(tx, certificate, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	total, err := c.CertificateRepository.CountByNameAndNotID(tx, request.Name, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get certificate user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	certificate.Name = request.Name
	certificate.Slug = slugHelper.Make(request.Name)
	certificate.Description = request.Description
	certificate.Price = request.Price
	certificate.IsActive = request.IsActive

	if err := c.CertificateRepository.Update(tx, certificate); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if request.CategoryUUID == uuid.Nil {
		if err := c.CertificateRepository.ClearCategory(tx, certificate); err != nil {
			c.Log.WithError(err).Error("error clearing certificate categories")
			return nil, fiber.ErrInternalServerError
		}
	} else {
		certifCategory := new(entity.CertificateCategory)
		if err := c.CertifCategoryRepository.FindByUUID(tx, certifCategory, request.CategoryUUID); err != nil {
			c.Log.WithError(err).Error("error finding certificate categories UUIDs")
			return nil, fiber.NewError(fiber.StatusBadRequest, "certificate categories UUIDs not valid")
		}

		if err := c.CertificateRepository.ReplaceCategory(tx, certificate, certifCategory); err != nil {
			c.Log.WithError(err).Error("error clearing certificate categories")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("Certificate updated : %+v", certificate)

	return converter.CertificateToResponse(certificate), nil
}

func (c *CertificateUseCase) UploadThumbnail(ctx context.Context, file *multipart.FileHeader, request *model.CertificateThumbnailRequest) (*model.CertificateResponse, error) {
	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	certificate := new(entity.Certificate)
	if err := c.CertificateRepository.FindByUUID(tx, certificate, request.CertificateUUID); err != nil {
		c.Log.Warnf("Failed find certificate from database : %+v", err)
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

	if err := c.CertificateRepository.AttachThumbnail(tx, certificate, newThumbnailFile); err != nil {
		c.Log.Warnf("Failed to attach thumbnail")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to attach thumbnail"+err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error uploading thumbnail")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertificateToResponse(certificate), nil

}
