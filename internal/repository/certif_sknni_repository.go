package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CertifSkkniRepository struct {
	Repository[entity.Skkni]
	Log *logrus.Logger
}

func NewCertifSkkniRepository(log *logrus.Logger) *CertifSkkniRepository {
	return &CertifSkkniRepository{
		Log: log,
	}
}

func (r *CertifSkkniRepository) CountByName(db *gorm.DB, name string) (int64, error) {
	var total int64
	err := db.Model(new(entity.Skkni)).Where("LOWER(name) = LOWER(?)", name).Count(&total).Error
	return total, err
}

func (r *CertifSkkniRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.Skkni)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CertifSkkniRepository) Search(db *gorm.DB, request *model.SearchCertificateSkkniRequest) ([]entity.Skkni, *model.PageMetadata, error) {
	var certificateSkkni []entity.Skkni

	var totalItems int64
	if err := db.Model(&entity.Skkni{}).Scopes(r.FilterSkkniSec(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	db = db.Preload("Certificate").Preload("File")

	if err := db.Scopes(r.FilterSkkniSec(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&certificateSkkni).Error; err != nil {
		return nil, nil, err
	}

	return certificateSkkni, pageMetadata, nil
}

func (r *CertifSkkniRepository) FilterSkkniSec(request *model.SearchCertificateSkkniRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("certificate_id = ?", request.CertificateID)

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}
		return tx
	}
}

func (r *CertifSkkniRepository) AttachFile(db *gorm.DB, skkni *entity.Skkni, file *entity.File) error {
	return db.Model(skkni).Association("File").Append(file)
}

func (r *CertifSkkniRepository) GetDetailByUUID(db *gorm.DB, skkni *entity.Skkni, uuid uuid.UUID) error {
	return db.Joins("Certificate").Joins("File").Where("skkni.uuid = ?", uuid).Take(skkni).Error
}
