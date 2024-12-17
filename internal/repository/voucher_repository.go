package repository

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VoucherRepository struct {
	Repository[entity.Voucher]
	Log *logrus.Logger
}

func NewVoucherRepository(log *logrus.Logger) *VoucherRepository {
	return &VoucherRepository{
		Log: log,
	}
}

func (r *VoucherRepository) CountByNameAndNotID(db *gorm.DB, name string, excludeUUID uuid.UUID) (int64, error) {
	var total int64
	err := db.Model(new(entity.Voucher)).
		Where("name = ? AND uuid != ?", name, excludeUUID).
		Count(&total).Error
	return total, err
}

func (r *VoucherRepository) ClearCourse(db *gorm.DB, voucher *entity.Voucher) error {
	return db.Model(voucher).Association("Courses").Clear()
}

func (r *VoucherRepository) SyncCourse(db *gorm.DB, voucher *entity.Voucher, courses []*entity.Course) error {
	return db.Model(voucher).Association("Courses").Append(courses)
}

func (r *VoucherRepository) FindWithDetails(db *gorm.DB, uuid uuid.UUID, withCourses bool, withActive bool) (*entity.Voucher, error) {
	var voucher entity.Voucher
	query := db

	if withCourses {
		query = query.Preload("Courses")
	}

	if withActive {
		query = query.Where("is_active = true")
	}

	err := query.First(&voucher, "uuid = ?", uuid).Error
	return &voucher, err
}

func (r *VoucherRepository) Search(db *gorm.DB, request *model.SearchVoucherRequest, withCourses bool) ([]entity.Voucher, *model.PageMetadata, error) {
	var vouchers []entity.Voucher

	var totalItems int64
	if err := db.Model(&entity.Voucher{}).Scopes(r.FilterVoucher(request)).Count(&totalItems).Error; err != nil {
		return nil, nil, err
	}

	pageMetadata := helper.CalculatePagination(totalItems, request.Page, request.Size)

	query := db.Scopes(r.FilterVoucher(request))

	if withCourses {
		query = query.Preload("Courses")
	}

	if err := query.Offset(pageMetadata.Offset).Limit(pageMetadata.Size).Find(&vouchers).Error; err != nil {
		return nil, nil, err
	}

	return vouchers, pageMetadata, nil
}

func (r *VoucherRepository) FilterVoucher(request *model.SearchVoucherRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}

		if code := request.Code; code != "" {
			code = "%" + code + "%"
			tx = tx.Where("code LIKE ?", code)
		}

		return tx
	}
}

func (r *VoucherRepository) ActivateVoucher(db *gorm.DB) error {
	return db.Model(new(entity.Voucher)).
		Where("is_active = ?", false).
		Where("start_active_at <= ?", time.Now()).
		Where("valid_until >= ?", time.Now()).
		Update("is_active", true).Error
}

func (r *VoucherRepository) DeactivateVoucher(db *gorm.DB) error {
	return db.Model(new(entity.Voucher)).
		Where("is_active = ?", true).
		Where("valid_until < ?", time.Now()).
		Or("start_active_at > ?", time.Now()).
		Update("is_active", false).Error
}

func (r *VoucherRepository) FindVoucherByCodeAndCourse(db *gorm.DB, courseId int, code string) (*entity.Voucher, error) {
	var voucher entity.Voucher

	err := db.Joins("JOIN course_voucher ON course_voucher.voucher_id = vouchers.id").
		Where("course_voucher.course_id = ? AND vouchers.code = ?", courseId, code).
		First(&voucher).Error

	if err != nil {
		return nil, err
	}

	return &voucher, nil
}
