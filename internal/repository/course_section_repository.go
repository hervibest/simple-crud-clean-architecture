package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
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
	err := db.Model(new(entity.CourseSection)).Where("LOWER(title) = LOWER(?)", title).Count(&total).Error
	return total, err
}

func (r *CourseSectionRepository) CountByTitleAndNotID(db *gorm.DB, title string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.CourseSection)).
		Where("title = ? AND uuid != ?", title, excludeUUID).
		Count(&total).Error
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

func (r *CourseSectionRepository) UpdateBetweenSequence(db *gorm.DB,
	courseId int, sequenceOld int, sequenceNew int, moreThanOld bool) error {
	if moreThanOld {
		return db.Model(new(entity.CourseSection)).
			Where("course_id = ? AND sequence >= ? AND sequence <= ?", courseId, sequenceOld, sequenceNew).
			Update("sequence", gorm.Expr("sequence - 1")).Error

	} else {
		return db.Model(new(entity.CourseSection)).
			Where("course_id = ? AND sequence >= ? AND sequence <= ?", courseId, sequenceNew, sequenceOld).
			Update("sequence", gorm.Expr("sequence + 1")).Error
	}
}

func (r *CourseSectionRepository) UpdateDecrementSequence(db *gorm.DB,
	courseId int, sequence int) error {
	return db.Model(new(entity.CourseSection)).
		Where("course_id = ? AND sequence >= ?", courseId, sequence).
		Update("sequence", gorm.Expr("sequence - 1")).Error
}

func (r *CourseSectionRepository) CountBySequence(db *gorm.DB, courseSection *entity.CourseSection, sequence int) (int64, error) {
	var total int64
	err := db.Model(courseSection).Where("sequence = ?", sequence).Count(&total).Error
	return total, err
}

func (r *CourseSectionRepository) Search(db *gorm.DB, request *model.SearchCourseSecRequest) ([]entity.CourseSection, *model.PageMetadata, error) {
	var courseSections []entity.CourseSection

	var totalItems int64
	if err := db.Model(&entity.CourseSection{}).Scopes(r.FilterCourseSec(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	if err := db.Scopes(r.FilterCourseSec(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&courseSections).Error; err != nil {
		return nil, nil, err
	}

	return courseSections, pageMetadata, nil
}

func (r *CourseSectionRepository) FilterCourseSec(request *model.SearchCourseSecRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("course_id = ?", request.CourseID)

		if title := request.Title; title != "" {
			title = "%" + title + "%"
			tx = tx.Where("title LIKE ?", title)
		}

		if description := request.Description; description != "" {
			description = "%" + description + "%"
			tx = tx.Where("description LIKE ?", description)
		}
		return tx
	}
}
