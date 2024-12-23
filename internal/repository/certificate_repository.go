package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CertificateRepository struct {
	Repository[entity.Certificate]
	Log *logrus.Logger
}

func NewCertificateRepository(log *logrus.Logger) *CertificateRepository {
	return &CertificateRepository{
		Log: log,
	}
}

func (r *CertificateRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.Certificate)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CertificateRepository) AddCategory(db *gorm.DB, certificate *entity.Certificate, category *entity.CertificateCategory) error {
	return db.Model(certificate).Association("Category").Append(category)
}

func (r *CertificateRepository) ReplaceCategory(db *gorm.DB, certificate *entity.Certificate, category *entity.CertificateCategory) error {
	return db.Model(certificate).Association("Category").Replace(category)
}

func (r *CertificateRepository) ClearCategory(db *gorm.DB, certificate *entity.Certificate) error {
	return db.Model(certificate).Association("Category").Clear()
}

func (r *CertificateRepository) FindWithDetails(db *gorm.DB, uuid uuid.UUID, withCategory bool, withMedia bool) (*entity.Certificate, error) {
	var certificate entity.Certificate
	query := db

	if withCategory {
		query = query.Preload("Category")
	}

	if withMedia {
		query = query.Preload("Thumbnail")
	}

	err := query.First(&certificate, "uuid = ?", uuid).Error
	return &certificate, err
}

func (r *CertificateRepository) Search(db *gorm.DB, request *model.SearchCertificateRequest, withCategories bool) ([]entity.Certificate, *model.PageMetadata, error) {
	var certificates []entity.Certificate

	var totalItems int64
	if err := db.Model(&entity.Certificate{}).Scopes(r.FilterCertificate(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	query := db.Scopes(r.FilterCertificate(request))

	if withCategories {
		query = query.Preload("Categories")
	}

	if err := query.Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&certificates).Error; err != nil {
		return nil, nil, err
	}

	return certificates, pageMetadata, nil
}

func (r *CertificateRepository) FilterCertificate(request *model.SearchCertificateRequest) func(tx *gorm.DB) *gorm.DB {
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

func (r *CertificateRepository) AttachThumbnail(db *gorm.DB, certificate *entity.Certificate, file *entity.File) error {
	return db.Model(certificate).Association("Thumbnail").Append(file)
}
