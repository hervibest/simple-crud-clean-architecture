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

	slugHelper "github.com/gosimple/slug"
)

type CourseUseCase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	Validate            *validator.Validate
	CourseRepository    *repository.CourseRepository
	CourseCatRepository *repository.CourseCategoryRepository
	Minio               *helper.Minio
}

func NewCourseUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	courseRepository *repository.CourseRepository, courseCatRepository *repository.CourseCategoryRepository,
	minio *helper.Minio) *CourseUseCase {
	return &CourseUseCase{
		DB:                  db,
		Log:                 logger,
		Validate:            validate,
		CourseRepository:    courseRepository,
		CourseCatRepository: courseCatRepository,
		Minio:               minio,
	}
}

func (c *CourseUseCase) Create(ctx context.Context, request *model.CreateCourseRequest) (*model.CourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.CourseRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed course user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	course := &entity.Course{
		UUID:        uuid.New(),
		Name:        request.Name,
		Slug:        slugHelper.Make(request.Name),
		Price:       request.Price,
		Description: request.Description,
	}

	if err := c.CourseRepository.Create(tx, course); err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	if request.CategoryUUIDs == nil {
		if err := c.CourseRepository.ClearCategory(tx, course); err != nil {
			c.Log.WithError(err).Error("error clearing course categories")
			return nil, fiber.ErrInternalServerError
		}
	} else if request.CategoryUUIDs != nil {
		courseCategories, err := c.CourseCatRepository.FindManyByUUIDs(tx, request.CategoryUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding course categories")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.CourseRepository.SyncCategory(tx, course, courseCategories); err != nil {
			c.Log.WithError(err).Error("error syncing course categories")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating course")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseToResponse(course), nil
}

func (c *CourseUseCase) Search(ctx context.Context, request *model.SearchCourseRequest) ([]model.CourseResponse, *model.PageMetadata, error) {

	courses, pageMetadata, err := c.CourseRepository.Search(c.DB, request, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting course category")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = *converter.CourseToResponse(&course)
	}

	return responses, pageMetadata, nil
}

func (c *CourseUseCase) Get(ctx context.Context, request *model.GetCourseRequest) (*model.CourseResponse, error) {

	course, err := c.CourseRepository.FindWithDetails(c.DB, request.UUID, true, true, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting course")
		return nil, fiber.ErrNotFound
	}

	return converter.CourseToResponse(course), nil
}

func (c *CourseUseCase) Update(ctx context.Context, request *model.UpdateCourseRequest) (*model.CourseResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(tx, course, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	total, err := c.CourseRepository.CountByNameAndNotID(tx, request.Name, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get course user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	course.Name = request.Name
	course.Slug = slugHelper.Make(request.Name)
	course.Description = request.Description
	course.Price = request.Price
	course.IsActive = request.IsActive

	if err := c.CourseRepository.Update(tx, course); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if request.CategoryUUIDs == nil {
		if err := c.CourseRepository.ClearCategory(tx, course); err != nil {
			c.Log.WithError(err).Error("error clearing course categories")
			return nil, fiber.ErrInternalServerError
		}
	} else if len(request.CategoryUUIDs) > 0 {
		courseCategories, err := c.CourseCatRepository.FindManyByUUIDs(tx, request.CategoryUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding course categories UUIDs")
			return nil, fiber.NewError(fiber.StatusBadRequest, "course categories UUIDs not valid")
		}

		if err := c.CourseRepository.ClearCategory(tx, course); err != nil {
			c.Log.WithError(err).Error("error clearing course categories")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.CourseRepository.SyncCategory(tx, course, courseCategories); err != nil {
			c.Log.WithError(err).Error("error syncing course categories")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("Course updated : %+v", course)

	return converter.CourseToResponse(course), nil
}

func (c *CourseUseCase) GetPurchasedCourseUUID(ctx context.Context, request *model.GetPurchasedCourseRequest) ([]model.CourseUUIDresponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	courses, err := c.CourseRepository.GetPurchasedCourseUUID(c.DB, request.UserID)
	if err != nil {
		c.Log.WithError(err).Error("error getting user purchased course")
		return nil, fiber.ErrNotFound
	}

	responses := make([]model.CourseUUIDresponse, len(courses))
	for i, course := range courses {
		responses[i] = *converter.CourseUUIDToResponse(&course)
	}

	return responses, nil

}

func (c *CourseUseCase) UserGetPurchasedCourse(ctx context.Context, request *model.SearchCourseRequest, userID int) ([]model.CourseResponse, *model.PageMetadata, error) {

	courses, pageMetadata, err := c.CourseRepository.UserGetPurchasedCourse(c.DB, request, userID)
	if err != nil {
		c.Log.WithError(err).Error("error getting user purchased course")
		return nil, nil, fiber.ErrNotFound
	}

	responses := make([]model.CourseResponse, len(courses))
	for i, course := range courses {
		responses[i] = *converter.CourseToResponse(&course)
	}

	return responses, pageMetadata, nil

}

func (c *CourseUseCase) UploadThumbnail(ctx context.Context, file *multipart.FileHeader, request *model.CourseThumbnailRequest) (*model.CourseResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(tx, course, request.CourseUUID); err != nil {
		c.Log.Warnf("Failed find course from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	upload, err := c.Minio.UploadFileToMinio(ctx, file, "thumbnail")
	if err != nil {
		return nil, err
	}

	newThumbnailFile := &entity.File{

		UUID:     uuid.New(),
		Filename: upload.Filename,
		Mimetype: upload.Mimetype,
		Path:     upload.URL,
		Size:     upload.Size,
	}

	if err := c.CourseRepository.AttachThumbnail(tx, course, newThumbnailFile); err != nil {
		c.Log.Warnf("Failed to attach video")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to attach thumbnail"+err.Error())
	}
	// OriginalName string    `gorm:"column:original_name"`
	// OriginalSize float64   `gorm:"column:original_size"`
	// OriginalMime string    `gorm:"column:original_mime"`
	// MediaID      string    `gorm:"column:media_id"`
	// Bucket       string    `gorm:"column:bucket"`
	// Dir          string    `gorm:"column:dir"`

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error uploading thumbnail")
		return nil, fiber.ErrInternalServerError
	}

	return converter.CourseToResponse(course), nil

}
