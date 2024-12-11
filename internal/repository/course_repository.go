package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseRepository struct {
	Repository[entity.Course]
	Log *logrus.Logger
}

func NewCourseRepository(log *logrus.Logger) *CourseRepository {
	return &CourseRepository{
		Log: log,
	}
}

func (r *CourseRepository) CountByName(db *gorm.DB, email string) (int64, error) {
	var total int64
	err := db.Model(new(entity.Course)).Where("name = ?", email).Count(&total).Error
	return total, err
}

func (r *CourseRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.Course)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *CourseRepository) ClearCategory(db *gorm.DB, course *entity.Course) error {
	return db.Model(course).Association("Categories").Clear()
}

func (r *CourseRepository) SyncCategory(db *gorm.DB, course *entity.Course, categories []*entity.CourseCategory) error {
	return db.Model(course).Association("Categories").Replace(categories)
}

func (r *CourseRepository) FindWithDetails(db *gorm.DB, uuid uuid.UUID, withCategories bool) (*entity.Course, error) {
	var course entity.Course
	query := db
	if withCategories {
		query = query.Preload("Categories")
	}

	err := query.First(&course, "uuid = ?", uuid).Error
	return &course, err
}

func (r *CourseRepository) Search(db *gorm.DB, request *model.SearchCourseRequest, withCategories bool) ([]entity.Course, *model.PageMetadata, error) {
	var courses []entity.Course

	var totalItems int64
	if err := db.Model(&entity.Course{}).Scopes(r.FilterCourse(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	query := db.Scopes(r.FilterCourse(request))

	if withCategories {
		query = query.Preload("Categories")
	}

	if err := query.Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&courses).Error; err != nil {
		return nil, nil, err
	}

	return courses, pageMetadata, nil
}

func (r *CourseRepository) FilterCourse(request *model.SearchCourseRequest) func(tx *gorm.DB) *gorm.DB {
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
