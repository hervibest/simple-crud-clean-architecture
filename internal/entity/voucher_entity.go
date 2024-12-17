package entity

import (
	"simple-crud-clean-architecture/internal/enum"
	"time"

	"github.com/google/uuid"
)

type Voucher struct {
	ID            int              `gorm:"column:id;primaryKey;autoIncrement"`
	UUID          uuid.UUID        `gorm:"column:uuid"`
	Name          string           `gorm:"column:name"`
	Code          string           `gorm:"column:code"`
	Value         float64          `gorm:"column:value"`
	IsActive      bool             `gorm:"column:is_active"`
	Type          enum.VoucherType `gorm:"column:type"`
	StartActiveAt time.Time        `gorm:"column:start_active_at"`
	ValidUntil    time.Time        `gorm:"column:valid_until"`
	CreatedAt     *time.Time       `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time       `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Courses       []Course         `gorm:"many2many:course_voucher;foreignKey:id;joinForeignKey:voucher_id;references:id;joinReferences:course_id"`
	Transactions  []Transaction    `gorm:"foreignKey:voucher_id;references:id"`
}

func (d *Voucher) TableName() string {
	return "vouchers"
}
