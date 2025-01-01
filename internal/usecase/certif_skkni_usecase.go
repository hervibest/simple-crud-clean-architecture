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
)

type CertifSkkniUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	CertificateRepository *repository.CertificateRepository
	CertifSkkniRepository *repository.CertifSkkniRepository
	Minio                 *helper.Minio
}

func NewCertifSkkniUseCase(db *gorm.DB, logger *logrus.Logger, certifSkkniRepository *repository.CertifSkkniRepository,
	certificateRepository *repository.CertificateRepository, minio *helper.Minio) *CertifSkkniUseCase {
	return &CertifSkkniUseCase{
		DB:                    db,
		Log:                   logger,
		CertificateRepository: certificateRepository,
		CertifSkkniRepository: certifSkkniRepository,
		Minio:                 minio,
	}
}

func (c *CertifSkkniUseCase) Search(ctx context.Context, request *model.SearchCertificateMatRequest) ([]model.CertificateSkkniResponse, *model.PageMetadata, error) {

	certificate := new(entity.Certificate)
	if err := c.CertificateRepository.FindByUUID(c.DB, certificate, request.CertificateUUID); err != nil {
		c.Log.Warnf("Failed find certificate from database : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid certificate uuid")
	}

	request.CertificateID = certificate.ID

	certifSkknis, pageMetadata, err := c.CertifSkkniRepository.Search(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting certificate category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CertificateSkkniResponse, len(certifSkknis))
	for i, certifSkkni := range certifSkknis {
		responses[i] = *converter.CertifSkkniToResponse(&certifSkkni)
	}

	return responses, pageMetadata, nil
}

func (c *CertifSkkniUseCase) Get(ctx context.Context, request *model.GetCertificateMatRequest) (*model.CertificateSkkniResponse, error) {
	certifSkkni := new(entity.Skkni)
	if err := c.CertifSkkniRepository.FindByUUID(c.DB, certifSkkni, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting certificate")
		return nil, fiber.ErrNotFound
	}

	return converter.CertifSkkniToResponse(certifSkkni), nil
}

func (c *CertifSkkniUseCase) Create(ctx context.Context, request *model.CreateCertificateMatRequest) (*model.CertificateSkkniResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.CertifSkkniRepository.CountByName(tx, request.Name)
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

	certifSkkni := new(entity.Skkni)
	certifSkkni.UUID = uuid.New()
	certifSkkni.Name = request.Name
	certifSkkni.CertificateID = certificate.ID

	if err := c.CertifSkkniRepository.Create(tx, certifSkkni); err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertifSkkniToResponse(certifSkkni), nil

}

func (c *CertifSkkniUseCase) Update(ctx context.Context, request *model.UpdateCertificateMatRequest) (*model.CertificateSkkniResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	certificate := new(entity.Certificate)
	if err := c.CertificateRepository.FindByUUID(tx, certificate, request.CertificateUUID); err != nil {
		c.Log.Warnf("Failed find certificate from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	certifSkkni := new(entity.Skkni)
	if err := c.CertifSkkniRepository.FindByUUID(tx, certifSkkni, request.CertificateMatUUID); err != nil {
		c.Log.Warnf("Failed find certificate section from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	total, err := c.CertifSkkniRepository.CountByNameAndNotID(tx, request.Name, request.CertificateMatUUID)
	if err != nil {
		c.Log.Warnf("Failed certificate user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Title already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "title already exists")
	}

	certifSkkni.Name = request.Name
	certifSkkni.CertificateID = certificate.ID

	if err := c.CertifSkkniRepository.Update(tx, certifSkkni); err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertifSkkniToResponse(certifSkkni), nil

}

func (c *CertifSkkniUseCase) Delete(ctx context.Context, request *model.DeleteCertificateMatRequest) (*model.CertificateSkkniResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	certificate := new(entity.Certificate)
	if err := c.CertificateRepository.FindByUUID(tx, certificate, request.CertificateUUID); err != nil {
		c.Log.Warnf("Failed find certificate from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid certificate uuid")

	}

	certifSkkni := new(entity.Skkni)
	if err := c.CertifSkkniRepository.FindByUUID(tx, certifSkkni, request.CertificateMatUUID); err != nil {
		c.Log.Warnf("Failed find certificate section from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid certificate section uuid")
	}

	if err := c.CertifSkkniRepository.Delete(tx, certifSkkni); err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating certificate")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertifSkkniToResponse(certifSkkni), nil
}

func (c *CertifSkkniUseCase) UploadFile(ctx context.Context, file *multipart.FileHeader, request *model.SkkniThumbnailRequest) (*model.CertificateSkkniResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	skkni := new(entity.Skkni)
	if err := c.CertifSkkniRepository.FindByUUID(tx, skkni, request.SkkniUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	upload, err := c.Minio.UploadFileToMinio(ctx, file, "file")
	if err != nil {
		c.Log.Warnf("Failed to upload skkni file to minio : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newThumbnailFile := &entity.File{

		UUID:     uuid.New(),
		Filename: upload.Filename,
		Mimetype: upload.Mimetype,
		Path:     upload.URL,
		Size:     upload.Size,
	}

	if err := c.CertifSkkniRepository.AttachFile(tx, skkni, newThumbnailFile); err != nil {
		c.Log.Warnf("Failed to attach video")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to attach file"+err.Error())
	}
	// OriginalName string    `gorm:"column:original_name"`
	// OriginalSize float64   `gorm:"column:original_size"`
	// OriginalMime string    `gorm:"column:original_mime"`
	// MediaID      string    `gorm:"column:media_id"`
	// Bucket       string    `gorm:"column:bucket"`
	// Dir          string    `gorm:"column:dir"`

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error uploading file")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CertifSkkniToResponse(skkni), nil

}
