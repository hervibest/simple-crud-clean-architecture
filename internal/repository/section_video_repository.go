package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SectionVideoRepository struct {
	Repository[entity.SectionVideo]
	Log *logrus.Logger
}

func NewSecVideoRepository(log *logrus.Logger) *SectionVideoRepository {
	return &SectionVideoRepository{
		Log: log,
	}
}

func (r *SectionVideoRepository) CountByTitle(db *gorm.DB, title string) (int64, error) {
	var total int64
	err := db.Model(new(entity.SectionVideo)).Where("title = ?", title).Count(&total).Error
	return total, err
}

func (r *SectionVideoRepository) CountByTitleAndNotID(db *gorm.DB, title string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.SectionVideo)).
		Where("title = ? AND uuid != ?", title, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *SectionVideoRepository) GetMaxSequence(db *gorm.DB, sectionId int) (int, error) {
	var maxSequence int
	err := db.Model(new(entity.SectionVideo)).
		Where("section_id = ?", sectionId).
		Select("COALESCE(MAX(sequence), 0)").
		Scan(&maxSequence).Error
	return maxSequence, err
}

func (r *SectionVideoRepository) UpdateIncrementSequence(db *gorm.DB,
	sectionId int, sequence int) error {
	return db.Model(new(entity.SectionVideo)).
		Where("section_id = ? AND sequence >= ?", sectionId, sequence).
		Update("sequence", gorm.Expr("sequence + 1")).Error
}

func (r *SectionVideoRepository) UpdateBetweenSequence(db *gorm.DB,
	sectionId int, sequenceOld int, sequenceNew int, moreThanOld bool) error {
	if moreThanOld {
		return db.Model(new(entity.SectionVideo)).
			Where("section_id = ? AND sequence >= ? AND sequence <= ?", sectionId, sequenceOld, sequenceNew).
			Update("sequence", gorm.Expr("sequence - 1")).Error

	} else {
		return db.Model(new(entity.SectionVideo)).
			Where("section_id = ? AND sequence >= ? AND sequence <= ?", sectionId, sequenceNew, sequenceOld).
			Update("sequence", gorm.Expr("sequence + 1")).Error
	}
}

func (r *SectionVideoRepository) UpdateDecrementSequence(db *gorm.DB,
	sectionId int, sequence int) error {
	return db.Model(new(entity.SectionVideo)).
		Where("section_id = ? AND sequence >= ?", sectionId, sequence).
		Update("sequence", gorm.Expr("sequence - 1")).Error
}

func (r *SectionVideoRepository) CountBySequence(db *gorm.DB, sectionVideo *entity.SectionVideo, sequence int) (int64, error) {
	var total int64
	err := db.Model(sectionVideo).Where("sequence = ?", sequence).Count(&total).Error
	return total, err
}

func (r *SectionVideoRepository) Search(db *gorm.DB, request *model.SearchSecVideosRequest) ([]entity.SectionVideo, *model.PageMetadata, error) {
	var sectionVideos []entity.SectionVideo

	var totalItems int64
	if err := db.Model(&entity.SectionVideo{}).Scopes(r.FilterSecVideo(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	if err := db.Scopes(r.FilterSecVideo(request)).Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&sectionVideos).Error; err != nil {
		return nil, nil, err
	}

	return sectionVideos, pageMetadata, nil
}

func (r *SectionVideoRepository) FilterSecVideo(request *model.SearchSecVideosRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tx = tx.Where("section_id = ?", request.SectionID)

		if title := request.Title; title != "" {
			title = "%" + title + "%"
			tx = tx.Where("title LIKE ?", title)
		}

		if notes := request.Notes; notes != "" {
			notes = "%" + notes + "%"
			tx = tx.Where("notes LIKE ?", notes)
		}
		return tx
	}
}
