package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CareerRepository struct {
	Repository[entity.Career]
	Log *logrus.Logger
}

func NewCareerRepository(log *logrus.Logger) *CareerRepository {
	return &CareerRepository{
		Log: log,
	}
}

func (r *CareerRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.Career)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CareerRepository) ClearCategory(db *gorm.DB, career *entity.Career) error {
	return db.Model(career).Association("Categories").Clear()
}

func (r *CareerRepository) SyncCategory(db *gorm.DB, career *entity.Career, categories []*entity.CareerCategory) error {
	return db.Model(career).Association("Categories").Append(categories)
}

func (r *CareerRepository) FindWithDetails(db *gorm.DB, uuid uuid.UUID, withCategories bool, withDiscount bool, withMedia bool) (*entity.Career, error) {
	var career entity.Career
	query := db

	if withCategories {
		query = query.Preload("Categories")
	}

	if withDiscount {
		query = query.Preload("Discounts", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("created_at asc").Limit(1)
		})
	}

	if withMedia {
		query = query.Preload("Thumbnail")
	}

	err := query.First(&career, "uuid = ?", uuid).Error
	return &career, err
}

func (r *CareerRepository) Search(db *gorm.DB, request *model.SearchCareerRequest, withCategories bool) ([]entity.Career, *model.PageMetadata, error) {
	var careers []entity.Career

	var totalItems int64
	if err := db.Model(&entity.Career{}).Scopes(r.FilterCareer(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	query := db.Scopes(r.FilterCareer(request))

	if withCategories {
		query = query.Preload("Categories")
	}

	if err := query.Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&careers).Error; err != nil {
		return nil, nil, err
	}

	return careers, pageMetadata, nil
}

func (r *CareerRepository) FilterCareer(request *model.SearchCareerRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if title := request.Title; title != "" {
			title = "%" + title + "%"
			tx = tx.Where("name LIKE ?", title)
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

func (r *CareerRepository) AttachThumbnail(db *gorm.DB, career *entity.Career, file *entity.File) error {
	return db.Model(career).Association("Thumbnail").Append(file)
}
