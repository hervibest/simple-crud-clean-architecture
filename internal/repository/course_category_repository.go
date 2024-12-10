package repository

import (
	"simple-crud-clean-architecture/internal/entity"

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

func (r *CourseCategoryRepository) CountByName(db *gorm.DB, email string) (int64, error) {
	var total int64
	err := db.Model(new(entity.CourseCategory)).Where("name = ?", email).Count(&total).Error
	return total, err
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
