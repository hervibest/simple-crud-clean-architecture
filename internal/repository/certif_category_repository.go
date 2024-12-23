package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CertifCategoryRepository struct {
	Repository[entity.CertificateCategory]
	Log *logrus.Logger
}

func NewCertifCategoryRepository(log *logrus.Logger) *CertifCategoryRepository {
	return &CertifCategoryRepository{
		Log: log,
	}
}

func (r *CertifCategoryRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.CertificateCategory)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CertifCategoryRepository) FindManyByUUIDs(db *gorm.DB, uuids []uuid.UUID) ([]*entity.CertificateCategory, error) {
	var careerCategories []*entity.CertificateCategory

	if err := db.Where("uuid IN ?", uuids).Find(&careerCategories).Error; err != nil {
		return nil, err
	}

	return careerCategories, nil
}

func (r *CertifCategoryRepository) Search(db *gorm.DB, request *model.SearchCertificateCatRequest) ([]entity.CertificateCategory, *model.PageMetadata, error) {
	var careerCats []entity.CertificateCategory

	var totalItems int64
	if err := db.Model(&entity.CertificateCategory{}).Scopes(r.FilterCertifCategory(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	if err := db.Scopes(r.FilterCertifCategory(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&careerCats).Error; err != nil {
		return nil, nil, err
	}

	return careerCats, pageMetadata, nil
}

func (r *CertifCategoryRepository) FilterCertifCategory(request *model.SearchCertificateCatRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}

		if slug := request.Slug; slug != "" {
			slug = "%" + slug + "%"
			tx = tx.Where("slug LIKE ?", slug)
		}

		if description := request.Description; description != "" {
			description = "%" + description + "%"
			tx = tx.Where("description LIKE ?", description)
		}

		return tx
	}
}
