package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CareerCategoryRepository struct {
	Repository[entity.CareerCategory]
	Log *logrus.Logger
}

func NewCareerCatRepository(log *logrus.Logger) *CareerCategoryRepository {
	return &CareerCategoryRepository{
		Log: log,
	}
}

func (r *CareerCategoryRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.CareerCategory)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CareerCategoryRepository) FindManyByUUIDs(db *gorm.DB, uuids []uuid.UUID) ([]*entity.CareerCategory, error) {
	var careerCategories []*entity.CareerCategory

	if err := db.Where("uuid IN ?", uuids).Find(&careerCategories).Error; err != nil {
		return nil, err
	}

	return careerCategories, nil
}

func (r *CareerCategoryRepository) Search(db *gorm.DB, request *model.SearchCareerCatRequest) ([]entity.CareerCategory, *model.PageMetadata, error) {
	var careerCats []entity.CareerCategory

	var totalItems int64
	if err := db.Model(&entity.CareerCategory{}).Scopes(r.FilterCareerCat(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	if err := db.Scopes(r.FilterCareerCat(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&careerCats).Error; err != nil {
		return nil, nil, err
	}

	return careerCats, pageMetadata, nil
}

func (r *CareerCategoryRepository) FilterCareerCat(request *model.SearchCareerCatRequest) func(tx *gorm.DB) *gorm.DB {
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
