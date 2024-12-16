package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DiscountRepository struct {
	Repository[entity.Discount]
	Log *logrus.Logger
}

func NewDiscountRepository(log *logrus.Logger) *DiscountRepository {
	return &DiscountRepository{
		Log: log,
	}
}

func (r *DiscountRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.Discount)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *DiscountRepository) ClearCourse(db *gorm.DB, course *entity.Discount) error {
	return db.Model(course).Association("Courses").Clear()
}

func (r *DiscountRepository) SyncCourse(db *gorm.DB, course *entity.Discount, courses []*entity.Course) error {
	return db.Model(course).Association("Courses").Append(courses)
}

func (r *DiscountRepository) FindWithDetails(db *gorm.DB, uuid uuid.UUID, withCourses bool, withActive bool) (*entity.Discount, error) {
	var course entity.Discount
	query := db

	if withCourses {
		query = query.Preload("Courses")
	}

	if withActive {
		query = query.Where("is_active = true")
	}

	err := query.First(&course, "uuid = ?", uuid).Error
	return &course, err
}

func (r *DiscountRepository) Search(db *gorm.DB, request *model.SearchDiscountRequest, withCourses bool) ([]entity.Discount, *model.PageMetadata, error) {
	var courses []entity.Discount

	var totalItems int64
	if err := db.Model(&entity.Discount{}).Scopes(r.FilterDiscount(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	query := db.Scopes(r.FilterDiscount(request))

	if withCourses {
		query = query.Preload("Courses")
	}

	if err := query.Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&courses).Error; err != nil {
		return nil, nil, err
	}

	return courses, pageMetadata, nil
}

func (r *DiscountRepository) FilterDiscount(request *model.SearchDiscountRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}

		if slug := request.Type; slug != "" {
			slug = "%" + slug + "%"
			tx = tx.Where("slug LIKE ?", slug)
		}

		return tx
	}
}

func (r *DiscountRepository) ActivateDiscount(db *gorm.DB) error {
	return db.Model(new(entity.Discount)).
		Where("is_active = ?", false).
		Where("start_active_at <= ?", time.Now()).
		Where("valid_until >= ?", time.Now()).
		Update("is_active", true).Error
}

func (r *DiscountRepository) DeactivateDiscount(db *gorm.DB) error {
	return db.Model(new(entity.Discount)).
		Where("is_active = ?", true).
		Where("valid_until < ?", time.Now()).
		Or("start_active_at > ?", time.Now()).
		Update("is_active", false).Error
}
