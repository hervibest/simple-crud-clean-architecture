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
)

type CertifMaterialUseCase struct {
	DB                           *gorm.DB
	Log                          *logrus.Logger
	CertificateRepository        *repository.CertificateRepository
	CertifMaterialtionRepository *repository.CertifMaterialRepository
	Validator                    helper.CustomValidator
}

func NewCertifMaterialUseCase(db *gorm.DB, logger *logrus.Logger, certifMaterialRepository *repository.CertifMaterialRepository,
	certificateRepository *repository.CertificateRepository, validator helper.CustomValidator) *CertifMaterialUseCase {
	return &CertifMaterialUseCase{
		DB:                           db,
		Log:                          logger,
		CertificateRepository:        certificateRepository,
		CertifMaterialtionRepository: certifMaterialRepository,
		Validator:                    validator,
	}
}

func (c *CertifMaterialUseCase) Search(ctx context.Context, request *model.SearchCertificateMatRequest) ([]model.CertificateMaterialResponse, *model.PageMetadata, error) {
	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, nil, validationErr
	}

	certificate := new(entity.Certificate)
	if err := c.CertificateRepository.FindByUUID(c.DB, certificate, request.CertificateUUID); err != nil {
		c.Log.Warnf("Failed find certificate from database : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid certificate uuid")
	}

	request.CertificateID = certificate.ID

	certifMaterials, pageMetadata, err := c.CertifMaterialtionRepository.Search(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting certificate category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CertificateMaterialResponse, len(certifMaterials))
	for i, certifMaterial := range certifMaterials {
		responses[i] = *converter.CertifMaterialToResponse(&certifMaterial)
	}

	return responses, pageMetadata, nil
}

func (c *CertifMaterialUseCase) Get(ctx context.Context, request *model.GetCertificateMatRequest) (*model.CertificateMaterialResponse, error) {
	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	certifMaterial := new(entity.Material)
	if err := c.CertifMaterialtionRepository.FindByUUID(c.DB, certifMaterial, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting certificate")
		return nil, fiber.ErrNotFound
	}

	return converter.CertifMaterialToResponse(certifMaterial), nil
}

func (c *CertifMaterialUseCase) Create(ctx context.Context, request *model.CreateCertificateMatRequest) (*model.CertificateMaterialResponse, error) {

	helper.SanitiseStruct(request)

	if validationErr := c.Validator.ValidateUseCase(request); validationErr != nil {
		c.Log.WithError(validationErr).Error("error validating request datas")
		return nil, validationErr
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.CertifMaterialtionRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed certificate user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	certificate := new(entity.Certificate)
	if err := c.CertificateRepository.FindByUUID(tx, certificate, request.CertificateUUID); err != nil {
		c.Log.Warnf("Failed find certificate from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	certifMaterial := new(entity.Material)
	certifMaterial.UUID = uuid.New()
	certifMaterial.Name = request.Name
	certifMaterial.Code = request.Code
	certifMaterial.Type = request.Type
	certifMaterial.CertificateID = certificate.ID

	if err := c.CertifMaterialtionRepository.Create(tx, certifMaterial); err != nil {
		c.Log.WithError(err).Error("error creating material")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating material")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertifMaterialToResponse(certifMaterial), nil

}

func (c *CertifMaterialUseCase) Update(ctx context.Context, request *model.UpdateCertificateMatRequest) (*model.CertificateMaterialResponse, error) {

	helper.SanitiseStruct(request)

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

	certifMaterial := new(entity.Material)
	if err := c.CertifMaterialtionRepository.FindByUUID(tx, certifMaterial, request.CertificateMatUUID); err != nil {
		c.Log.Warnf("Failed find certificate section from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	total, err := c.CertifMaterialtionRepository.CountByNameAndNotID(tx, request.Name, request.CertificateMatUUID)
	if err != nil {
		c.Log.Warnf("Failed certificate user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Title already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "title already exists")
	}

	certifMaterial.Name = request.Name
	certifMaterial.Code = request.Code
	certifMaterial.Type = request.Type
	certifMaterial.CertificateID = certificate.ID

	if err := c.CertifMaterialtionRepository.Update(tx, certifMaterial); err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertifMaterialToResponse(certifMaterial), nil

}

func (c *CertifMaterialUseCase) Delete(ctx context.Context, request *model.DeleteCertificateMatRequest) (*model.CertificateMaterialResponse, error) {

	helper.SanitiseStruct(request)

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
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid certificate uuid")

	}

	certifMaterial := new(entity.Material)
	if err := c.CertifMaterialtionRepository.FindByUUID(tx, certifMaterial, request.CertificateMatUUID); err != nil {
		c.Log.Warnf("Failed find certificate section from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid certificate section uuid")
	}

	if err := c.CertifMaterialtionRepository.Delete(tx, certifMaterial); err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertifMaterialToResponse(certifMaterial), nil
}
