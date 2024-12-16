package entity

import (
	"simple-crud-clean-architecture/internal/enum"
	"time"

	"github.com/google/uuid"
)

type Discount struct {
	ID            int               `gorm:"column:id;primaryKey;autoIncrement"`
	UUID          uuid.UUID         `gorm:"column:uuid"`
	Name          string            `gorm:"column:name"`
	Value         float64           `gorm:"column:value"`
	IsActive      bool              `gorm:"column:is_active"`
	Type          enum.DiscountType `gorm:"column:type"`
	StartActiveAt time.Time         `gorm:"column:start_active_at"`
	ValidUntil    time.Time         `gorm:"column:valid_until"`
	CreatedAt     *time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     *time.Time        `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Courses       []Course          `gorm:"many2many:course_discount;foreignKey:id;joinForeignKey:discount_id;references:id;joinReferences:course_id"`
}

func (d *Discount) TableName() string {
	return "discounts"
}
