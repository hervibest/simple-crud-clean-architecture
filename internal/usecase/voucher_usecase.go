package usecase

import (
	"context"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/enum"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/model/converter"
	"simple-crud-clean-architecture/internal/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VoucherUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	VoucherRepository *repository.VoucherRepository
	CourseRepository  *repository.CourseRepository
}

func NewVoucherUseCase(db *gorm.DB, logger *logrus.Logger, voucherRepository *repository.VoucherRepository,
	courseRepository *repository.CourseRepository) *VoucherUseCase {
	return &VoucherUseCase{
		DB:                db,
		Log:               logger,
		VoucherRepository: voucherRepository,
		CourseRepository:  courseRepository,
	}
}

func (c *VoucherUseCase) Create(ctx context.Context, request *model.CreateVoucherRequest) (*model.VoucherResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	helper.SanitiseStruct(request)

	total, err := c.VoucherRepository.CountByName(tx, request.Name)
	if err != nil {
		c.Log.Warnf("Failed voucher user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	voucher := &entity.Voucher{
		UUID:          uuid.New(),
		Name:          request.Name,
		Code:          request.Code,
		Value:         request.Value,
		Type:          request.Type,
		StartActiveAt: request.StartActiveAt,
		ValidUntil:    request.ValidUntil,
	}

	if err := c.VoucherRepository.Create(tx, voucher); err != nil {
		c.Log.WithError(err).Error("error creating voucher")
		return nil, fiber.ErrInternalServerError
	}

	if request.CourseUUIDs == nil {
		if err := c.VoucherRepository.ClearCourse(tx, voucher); err != nil {
			c.Log.WithError(err).Error("error clearing course")
			return nil, fiber.ErrInternalServerError
		}
	} else if request.CourseUUIDs != nil {
		courses, err := c.CourseRepository.FindManyByUUIDs(tx, request.CourseUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding course")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.VoucherRepository.SyncCourse(tx, voucher, courses); err != nil {
			c.Log.WithError(err).Error("error syncing course")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error creating voucher")
		return nil, fiber.ErrInternalServerError
	}

	return converter.VoucherToResponse(voucher), nil
}

func (c *VoucherUseCase) Validate(ctx context.Context, request *model.ValidateVoucherRequest) error {

	helper.SanitiseStruct(request)

	if request.UUID != uuid.Nil {
		voucher := new(entity.Voucher)
		if err := c.VoucherRepository.FindByUUID(c.DB, voucher, request.UUID); err != nil {
			c.Log.Warnf("Failed voucher user from database : %+v", err)
			return fiber.ErrInternalServerError
		}

		if voucher.IsActive {
			return fiber.NewError(fiber.StatusBadRequest, "current valid and active voucher cannot be update")
		}
	}

	courses, err := c.CourseRepository.FindManyByUUIDs(c.DB, request.CourseUUIDs)
	if err != nil {
		c.Log.WithError(err).Error("error finding course")
		return fiber.ErrInternalServerError
	}
	c.Log.Infof("test debug courses" + courses[0].Name)

	for _, course := range courses {
		coursePriceToCompare := course.Price
		if request.Type == enum.VoucherTypeRebate && coursePriceToCompare < request.Value {
			return fiber.NewError(fiber.StatusUnprocessableEntity, "requested Voucher will cause a minus price on course")
		}
		c.Log.Infof("ini adalah perhitungan validasi" + string(request.Type) + string(enum.VoucherTypePercent))

		if request.Type == enum.VoucherTypePercent {
			coursePriceToCompare = course.Price - request.Value/100*course.Price
			c.Log.Infof("ini adalah perhitungan validasi" + strconv.FormatFloat(coursePriceToCompare, 'f', 2, 64))
			if coursePriceToCompare <= 0 {
				return fiber.NewError(fiber.StatusUnprocessableEntity, "requested Voucher will cause a minus price on course")

			}
		}
	}

	return nil
}

func (c *VoucherUseCase) Update(ctx context.Context, request *model.UpdateVoucherRequest) (*model.VoucherResponse, error) {

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	voucher := new(entity.Voucher)
	if err := c.VoucherRepository.FindByUUID(tx, voucher, request.UUID); err != nil {
		c.Log.WithError(err).Error("error getting contact")
		return nil, fiber.ErrNotFound
	}

	total, err := c.VoucherRepository.CountByNameAndNotID(tx, request.Name, request.UUID)
	if err != nil {
		c.Log.Warnf("Failed to get voucher user from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	c.Log.Infof("Total : %d", total)

	if total > 0 {
		c.Log.Warnf("Name already exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, "name already exists")
	}

	voucher.Name = request.Name
	voucher.Code = request.Code
	voucher.Type = request.Type
	voucher.StartActiveAt = request.StartActiveAt
	voucher.ValidUntil = request.ValidUntil
	voucher.Value = request.Value

	if err := c.VoucherRepository.Update(tx, voucher); err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	if request.CourseUUIDs == nil {
		if err := c.VoucherRepository.ClearCourse(tx, voucher); err != nil {
			c.Log.WithError(err).Error("error clearing voucher categories")
			return nil, fiber.ErrInternalServerError
		}
	} else if len(request.CourseUUIDs) > 0 {
		courses, err := c.CourseRepository.FindManyByUUIDs(tx, request.CourseUUIDs)
		if err != nil {
			c.Log.WithError(err).Error("error finding voucher categories UUIDs")
			return nil, fiber.NewError(fiber.StatusBadRequest, "voucher categories UUIDs not valid")
		}

		if err := c.VoucherRepository.ClearCourse(tx, voucher); err != nil {
			c.Log.WithError(err).Error("error clearing voucher categories")
			return nil, fiber.ErrInternalServerError
		}

		if err := c.VoucherRepository.SyncCourse(tx, voucher, courses); err != nil {
			c.Log.WithError(err).Error("error syncing voucher categories")
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("error updating contact")
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("Voucher updated : %+v", voucher)

	return converter.VoucherToResponse(voucher), nil
}

func (c *VoucherUseCase) Search(ctx context.Context, request *model.SearchVoucherRequest) ([]model.VoucherResponse, *model.PageMetadata, error) {

	courses, pageMetadata, err := c.VoucherRepository.Search(c.DB, request, true)
	if err != nil {
		c.Log.WithError(err).Error("error getting voucher")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.VoucherResponse, len(courses))
	for i, course := range courses {
		responses[i] = *converter.VoucherToResponse(&course)
	}

	return responses, pageMetadata, nil
}

func (c *VoucherUseCase) Get(ctx context.Context, request *model.GetVoucherRequest) (*model.VoucherResponse, error) {

	course, err := c.VoucherRepository.FindWithDetails(c.DB, request.UUID, true, false)
	if err != nil {
		c.Log.WithError(err).Error("error getting course")
		return nil, fiber.ErrNotFound
	}

	return converter.VoucherToResponse(course), nil
}

func (c *VoucherUseCase) ApplyVoucher(ctx context.Context, request *model.ApplyVoucherRequest) (*model.VoucherResponse, error) {

	course := new(entity.Course)
	if err := c.CourseRepository.FindByUUID(c.DB, course, request.CourseUUID); err != nil {
		c.Log.WithError(err).Error("invald course")
		return nil, fiber.ErrNotFound
	}

	voucher, err := c.VoucherRepository.FindVoucherByCodeAndCourse(c.DB, course.ID, request.VoucherCode)
	if err != nil {
		c.Log.WithError(err).Error("invald voucher code or course")
		return nil, fiber.ErrNotFound
	}

	if voucher.IsActive != true {
		c.Log.WithError(err).Error("voucher has expired")
		return nil, fiber.NewError(fiber.StatusBadRequest, "voucher is expired")
	}

	return converter.VoucherToResponse(voucher), nil
}

func (c *VoucherUseCase) ActivateVoucher(ctx context.Context) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.VoucherRepository.ActivateVoucher(tx); err != nil {
		c.Log.Warnf("Error activating voucher using cron" + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Error comitting tx activate voucher using cron" + err.Error())
		return err
	}

	c.Log.Warnf("Activate voucher using Job")

	return nil
}

func (c *VoucherUseCase) DeactivateVoucher(ctx context.Context) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.VoucherRepository.DeactivateVoucher(tx); err != nil {
		c.Log.Warnf("Error deactivating voucher using cron" + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Error commiting tx deactivate voucher using cron" + err.Error())
		return err
	}

	c.Log.Warnf("Deactivate voucher using Job")

	return nil
}
