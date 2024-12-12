package repository

import (
	"simple-crud-clean-architecture/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseSectionRepository struct {
	Repository[entity.CourseSection]
	Log *logrus.Logger
}

func NewCourseSectionRepository(log *logrus.Logger) *CourseSectionRepository {
	return &CourseSectionRepository{
		Log: log,
	}
}

func (r *CourseSectionRepository) CountByTitle(db *gorm.DB, title string) (int64, error) {
	var total int64
	err := db.Model(new(entity.CourseSection)).Where("title = ?", title).Count(&total).Error
	return total, err
}

func (r *CourseSectionRepository) GetMaxSequence(db *gorm.DB, courseId int) (int, error) {
	var maxSequence int
	err := db.Model(new(entity.CourseSection)).
		Where("course_id = ?", courseId).
		Select("COALESCE(MAX(sequence), 0)").
		Scan(&maxSequence).Error
	return maxSequence, err
}
func (r *CourseSectionRepository) UpdateIncrementSequence(db *gorm.DB,
	courseId int, sequence int) error {
	return db.Model(new(entity.CourseSection)).
		Where("course_id = ? AND sequence >= ?", courseId, sequence).
		Update("sequence", gorm.Expr("sequence + 1")).Error
}
