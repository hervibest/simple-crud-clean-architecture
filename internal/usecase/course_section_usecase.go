package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseSecUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	CourseRepository        *repository.CourseRepository
	CourseSectionRepository *repository.CourseSectionRepository
}

func NewCourseSecUseCase(db *gorm.DB, logger *logrus.Logger,
	courseSecRepository *repository.CourseSectionRepository, courseRepository *repository.CourseRepository) *CourseSecUseCase {
	return &CourseSecUseCase{
		DB:                      db,
		Log:                     logger,
		CourseRepository:        courseRepository,
		CourseSectionRepository: courseSecRepository}
}

func (c *CourseSecUseCase) Search(ctx context.Context, request *model.SearchCourseSecRequest) ([]model.CourseSectionResponse, *model.PageMetadata, error) {

	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(c.DB, course, request.CourseUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")
	}

	request.CourseID = course.ID

	courseSections, pageMetadata, err := c.CourseSectionRepository.Search(c.DB, request)
	if err != nil {
		c.Log.WithError(err).Error("error getting course category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CourseSectionResponse, len(courseSections))
	for i, courseSec := range courseSections {
		responses[i] = *converter.CourseSecToResponse(&courseSec)
	}

	return responses, pageMetadata, nil
}

func (c *CourseSecUseCase) Get(ctx context.Context, request *model.GetCourseSecRequest) (*model.CourseSectionResponse, error) {
	courseSection := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(c.DB, courseSection, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting course")
		return nil, fiber.ErrNotFound
	}

	return converter.CourseSecToResponse(courseSection), nil
}

func (c *CourseSecUseCase) Create(ctx context.Context, request *model.CreateCourseSecRequest) (*model.CourseSectionResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.CourseSectionRepository.CountByTitle(tx, request.Title)
	if err != nil {
		c.Log.Warnf("Failed course user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(tx, course, request.CourseUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	courseMaxSequence, err := c.CourseSectionRepository.GetMaxSequence(tx, course.ID)
	if err != nil {
		c.Log.Warnf("Error getting max sequence of course")
		return nil, fiber.ErrInternalServerError
	}

	courseSection := new(entity.CourseSection)
	courseSection.UUID = uuid.New()
	courseSection.Title = request.Title
	courseSection.Description = request.Description
	courseSection.CourseID = course.ID

	if courseMaxSequence == 0 {
		courseSection.Sequence = 1
	} else {
		courseSection.Sequence = request.Sequence
	}

	if courseSection.Sequence < courseMaxSequence+1 {
		c.CourseSectionRepository.UpdateIncrementSequence(tx, course.ID, courseSection.Sequence)
	}

	if err := c.CourseSectionRepository.Create(tx, courseSection); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseSecToResponse(courseSection), nil

}

func (c *CourseSecUseCase) Update(ctx context.Context, request *model.UpdateCourseSecRequest) (*model.CourseSectionResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(tx, course, request.CourseUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	courseSection := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(tx, courseSection, request.CourseSecUUID); err != nil {
		c.Log.Warnf("Failed find course section from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	total, err := c.CourseSectionRepository.CountByTitleAndNotID(tx, request.Title, request.CourseSecUUID)
	if err != nil {
		c.Log.Warnf("Failed course user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Title already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "title already exists")
	}

	courseSection.Title = request.Title
	courseSection.Description = request.Description
	courseSection.CourseID = course.ID

	sequenceIsUsed := false

	total, err = c.CourseSectionRepository.CountBySequence(tx, courseSection, request.Sequence)
	if err != nil {
		c.Log.Warnf("Error getting sequence of course section")
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		sequenceIsUsed = true
	}

	sequenceOld := courseSection.Sequence

	//implement stack logic

	if sequenceIsUsed && request.Sequence != courseSection.Sequence {
		var moreThanOld bool
		if sequenceOld < request.Sequence {
			moreThanOld = true
			if err := c.CourseSectionRepository.UpdateBetweenSequence(tx, course.ID, sequenceOld, request.Sequence, moreThanOld); err != nil {
				return nil, fiber.ErrInternalServerError
			}
		} else {
			moreThanOld = false
			if err := c.CourseSectionRepository.UpdateBetweenSequence(tx, course.ID, sequenceOld, request.Sequence, moreThanOld); err != nil {
				return nil, fiber.ErrInternalServerError
			}
		}
	}

	courseSection.Sequence = request.Sequence

	if err := c.CourseSectionRepository.Update(tx, courseSection); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseSecToResponse(courseSection), nil

}

func (c *CourseSecUseCase) Delete(ctx context.Context, request *model.DeleteCourseSecRequest) (*model.CourseSectionResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(tx, course, request.CourseUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid course uuid")

	}

	courseSection := new(entity.CourseSection)
	if err := c.CourseSectionRepository.FindByUUID(tx, courseSection, request.CourseSecUUID); err != nil {
		c.Log.Warnf("Failed find course section from database : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid course section uuid")
	}

	courseMaxSequence, err := c.CourseSectionRepository.GetMaxSequence(tx, course.ID)
	if err != nil {
		c.Log.Warnf("Error getting max sequence of course")
		return nil, fiber.ErrInternalServerError
	}

	if courseSection.Sequence < courseMaxSequence {
		c.CourseSectionRepository.UpdateDecrementSequence(tx, course.ID, courseSection.Sequence)
	}

	if err := c.CourseSectionRepository.Delete(tx, courseSection); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseSecToResponse(courseSection), nil
}
