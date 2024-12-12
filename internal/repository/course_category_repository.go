package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseCategoryRepository struct {
	Repository[entity.CourseCategory]
	Log *logrus.Logger
}

func NewCourseCatRepository(log *logrus.Logger) *CourseCategoryRepository {
	return &CourseCategoryRepository{
		Log: log,
	}
}

func (r *CourseCategoryRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.CourseCategory)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CourseCategoryRepository) FindManyByUUIDs(db *gorm.DB, uuids []uuid.UUID) ([]*entity.CourseCategory, error) {
	var courseCategories []*entity.CourseCategory

	if err := db.Where("uuid IN ?", uuids).Find(&courseCategories).Error; err != nil {
		return nil, err
	}

	return courseCategories, nil
}

func (r *CourseCategoryRepository) Search(db *gorm.DB, request *model.SearchCourseCatRequest) ([]entity.CourseCategory, *model.PageMetadata, error) {
	var courseCats []entity.CourseCategory

	var totalItems int64
	if err := db.Model(&entity.CourseCategory{}).Scopes(r.FilterCourseCat(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	if err := db.Scopes(r.FilterCourseCat(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&courseCats).Error; err != nil {
		return nil, nil, err
	}

	return courseCats, pageMetadata, nil
}

func (r *CourseCategoryRepository) FilterCourseCat(request *model.SearchCourseCatRequest) func(tx *gorm.DB) *gorm.DB {
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
