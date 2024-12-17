package entity

import (
	"encoding/json"
	"simple-crud-clean-architecture/internal/enum"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID        int                    `gorm:"column:id;primaryKey;autoIncrement"`
	TrxID     uuid.UUID              `gorm:"column:trx_id"`
	UserID    int                    `gorm:"column:user_id"`
	CourseID  int                    `gorm:"column:course_id"`
	Amount    float64                `gorm:"column:amount"`
	Status    enum.TransactionStatus `gorm:"column:status"`
	VoucherID int                    `gorm:"coloumn:voucher_id"`

	User    *User    `gorm:"foreignKey:user_id"`
	Course  *Course  `gorm:"foreignKey:course_id"`
	Voucher *Voucher `gorm:"foreignKey:voucher_id"`

	SnapToken                string                     `gorm:"column:snap_token"`
	ExternalStatus           enum.MidtransPaymentStatus `gorm:"column:external_status"`
	ExternalCallbackResponse json.RawMessage            `gorm:"column:external_callback_response"`
	PaidAt                   time.Time                  `gorm:"column:paid_at"`
	CreatedAt                time.Time                  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt                time.Time                  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *Transaction) TableName() string {
	return "transactions"
}
