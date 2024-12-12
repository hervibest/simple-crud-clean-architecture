package usecase

import (
	"context"
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

type CourseSecUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validate                *validator.Validate
	CourseRepository        *repository.CourseRepository
	CourseSectionRepository *repository.CourseSectionRepository
}

func NewCourseSecUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	courseSecRepository *repository.CourseSectionRepository, courseRepository *repository.CourseRepository) *CourseSecUseCase {
	return &CourseSecUseCase{
		DB:                      db,
		Log:                     logger,
		Validate:                validate,
		CourseRepository:        courseRepository,
		CourseSectionRepository: courseSecRepository}
}

func (c *CourseSecUseCase) Create(ctx context.Context, request *model.CreateCourseSecRequest) (*model.CourseSectionResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

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
