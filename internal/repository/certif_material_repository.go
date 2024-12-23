package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CertifMaterialRepository struct {
	Repository[entity.Material]
	Log *logrus.Logger
}

func NewCertifMaterialRepository(log *logrus.Logger) *CertifMaterialRepository {
	return &CertifMaterialRepository{
		Log: log,
	}
}

func (r *CertifMaterialRepository) CountByName(db *gorm.DB, name string) (int64, error) {
	var total int64
	err := db.Model(new(entity.Material)).Where("LOWER(name) = LOWER(?)", name).Count(&total).Error
	return total, err
}

func (r *CertifMaterialRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.Material)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CertifMaterialRepository) Search(db *gorm.DB, request *model.SearchCertificateMatRequest) ([]entity.Material, *model.PageMetadata, error) {
	var certifMaterials []entity.Material

	var totalItems int64
	if err := db.Model(&entity.Material{}).Scopes(r.FilterMaterialSec(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	if err := db.Scopes(r.FilterMaterialSec(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&certifMaterials).Error; err != nil {
		return nil, nil, err
	}

	return certifMaterials, pageMetadata, nil
}

func (r *CertifMaterialRepository) FilterMaterialSec(request *model.SearchCertificateMatRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("certificate_id = ?", request.CertificateID)

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}
		if code := request.Code; code != "" {
			code = "%" + code + "%"
			tx = tx.Where("code LIKE ?", code)
		}

		return tx
	}
}
