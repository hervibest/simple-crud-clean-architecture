package repository

import (
	"simple-crud-clean-architecture/internal/entity"

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
