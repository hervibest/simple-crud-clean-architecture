package usecase

import (
	"context"
	"mime/multipart"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SecVideoUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validate                *validator.Validate
	CourseSectionRepository *repository.CourseSectionRepository
	SectionVideoRepository  *repository.SectionVideoRepository
	Minio                   *helper.Minio
}

func NewSecVideoUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	courseSecRepository *repository.CourseSectionRepository, sectionVideoRepository *repository.SectionVideoRepository,
	minio *helper.Minio) *SecVideoUseCase {
	return &SecVideoUseCase{
		DB:                      db,
		Log:                     logger,
		Validate:                validate,
		CourseSectionRepository: courseSecRepository,
		SectionVideoRepository:  sectionVideoRepository,
		Minio:                   minio,
	}
}

func (c *SecVideoUseCase) Search(ctx context.Context, request *model.SearchSecVideosRequest) ([]model.SecVideoResponse, *model.PageMetadata, error) {

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, nil, fiber.ErrBadRequest
	}

	section := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(c.DB, section, request.SectionUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")
	}

	request.SectionID = section.ID

	courseSections, pageMetadata, err := c.SectionVideoRepository.Search(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting course category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.SecVideoResponse, len(courseSections))
	for i, courseSec := range courseSections {
		responses[i] = *converter.SecVideoToResponse(&courseSec)
	}

	return responses, pageMetadata, nil
}

func (c *SecVideoUseCase) Get(ctx context.Context, request *model.GetSecVideoRequest) (*model.SecVideoResponse, error) {

	sectionVideo := new(entity.SectionVideo)
	if err := c.SectionVideoRepository.FindByUUID(c.DB, sectionVideo, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting video")
		return nil, fiber.ErrNotFound
	}

	return converter.SecVideoToResponse(sectionVideo), nil
}

func (c *SecVideoUseCase) Create(ctx context.Context, request *model.CreateSecVideoRequest) (*model.SecVideoResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.SectionVideoRepository.CountByTitle(tx, request.Title)
	if err != nil {
		c.Log.Warnf("Failed course user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Title already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "Title already exists")
	}

	section := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(tx, section, request.SectionUUID); err != nil {
		c.Log.Warnf("Failed find section from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	videoMaxSequence, err := c.SectionVideoRepository.GetMaxSequence(tx, section.ID)
	if err != nil {
		c.Log.Warnf("Error getting max sequence of section")
		return nil, fiber.ErrInternalServerError
	}

	sectionVideo := new(entity.SectionVideo)
	sectionVideo.UUID = uuid.New()
	sectionVideo.Title = request.Title
	sectionVideo.Notes = request.Notes
	sectionVideo.SectionID = section.ID

	if videoMaxSequence == 0 {
		sectionVideo.Sequence = 1
	} else {
		sectionVideo.Sequence = request.Sequence
	}

	if sectionVideo.Sequence < videoMaxSequence+1 {
		c.CourseSectionRepository.UpdateIncrementSequence(tx, section.ID, sectionVideo.Sequence)
	}

	if err := c.SectionVideoRepository.Create(tx, sectionVideo); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.SecVideoToResponse(sectionVideo), nil

}

func (c *SecVideoUseCase) Update(ctx context.Context, request *model.UpdateSecVideoRequest) (*model.SecVideoResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	section := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(tx, section, request.SectionUUID); err != nil {
		c.Log.Warnf("Failed find section from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	sectionVideo := new(entity.SectionVideo)
	if err := c.SectionVideoRepository.FindByUUID(tx, sectionVideo, request.VideoUUID); err != nil {
		c.Log.Warnf("Failed find video from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	total, err := c.SectionVideoRepository.CountByTitleAndNotID(tx, request.Title, request.VideoUUID)
	if err != nil {
		c.Log.Warnf("Failed course user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Title already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "title already exists")
	}

	sectionVideo.Title = request.Title
	sectionVideo.Notes = request.Description
	sectionVideo.SectionID = section.ID

	sequenceIsUsed := false

	total, err = c.SectionVideoRepository.CountBySequence(tx, sectionVideo, request.Sequence)
	if err != nil {
		c.Log.Warnf("Error getting sequence of section video")
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		sequenceIsUsed = true
	}

	sequenceOld := sectionVideo.Sequence

	if sequenceIsUsed && request.Sequence != sectionVideo.Sequence {
		var moreThanOld bool
		if sequenceOld < request.Sequence {
			moreThanOld = true
			if err := c.SectionVideoRepository.UpdateBetweenSequence(tx, section.ID, sequenceOld, request.Sequence, moreThanOld); err != nil {
				return nil, fiber.ErrInternalServerError
			}
		} else {
			moreThanOld = false
			if err := c.SectionVideoRepository.UpdateBetweenSequence(tx, section.ID, sequenceOld, request.Sequence, moreThanOld); err != nil {
				return nil, fiber.ErrInternalServerError
			}
		}
	}

	sectionVideo.Sequence = request.Sequence

	if err := c.SectionVideoRepository.Update(tx, sectionVideo); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.SecVideoToResponse(sectionVideo), nil

}

func (c *SecVideoUseCase) Delete(ctx context.Context, request *model.DeleteSecVideoRequest) (*model.SecVideoResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	section := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(tx, section, request.SectionUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")

	}

	sectionVideo := new(entity.SectionVideo)
	if err := c.SectionVideoRepository.FindByUUID(tx, sectionVideo, request.VideoUUID); err != nil {
		c.Log.Warnf("Failed find course section from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid course section uuid")
	}

	videoMaxSequence, err := c.SectionVideoRepository.GetMaxSequence(tx, section.ID)
	if err != nil {
		c.Log.Warnf("Error getting max sequence of course")
		return nil, fiber.ErrInternalServerError
	}

	if sectionVideo.Sequence < videoMaxSequence {
		c.SectionVideoRepository.UpdateDecrementSequence(tx, section.ID, sectionVideo.Sequence)
	}

	if err := c.SectionVideoRepository.Delete(tx, sectionVideo); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.SecVideoToResponse(sectionVideo), nil
}

func (c *SecVideoUseCase) UploadVideo(ctx context.Context, file *multipart.FileHeader, request *model.UploadVideoRequest) (*model.SecVideoResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	section := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(tx, section, request.SectionUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	sectionVideo := new(entity.SectionVideo)
	if err := c.SectionVideoRepository.FindByUUID(tx, sectionVideo, request.VideoUUID); err != nil {
		c.Log.Warnf("Failed find course section from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	upload, err := c.Minio.UploadFileToMinio(ctx, file, "video")
	if err != nil {
		return nil, err
	}

	newVideoFile := &entity.File{

		UUID:     uuid.New(),
		Filename: upload.Filename,
		Mimetype: upload.Mimetype,
		Path:     upload.URL,
		Size:     upload.Size,
	}

	if err := c.SectionVideoRepository.AttachVideo(tx, sectionVideo, newVideoFile); err != nil {
		c.Log.Warnf("Failed to attach video")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to attach video"+err.Error())
	}
	// OriginalName string    `gorm:"column:original_name"`
	// OriginalSize float64   `gorm:"column:original_size"`
	// OriginalMime string    `gorm:"column:original_mime"`
	// MediaID      string    `gorm:"column:media_id"`
	// Bucket       string    `gorm:"column:bucket"`
	// Dir          string    `gorm:"column:dir"`

	sectionVideo.OriginalName = newVideoFile.Filename
	sectionVideo.OriginalSize = float64(newVideoFile.Size)
	sectionVideo.OriginalMime = newVideoFile.Mimetype
	sectionVideo.Dir = newVideoFile.Path

	if err := c.SectionVideoRepository.Update(tx, sectionVideo); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.SecVideoToResponse(sectionVideo), nil

}
