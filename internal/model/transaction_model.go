package model

import (
	"encoding/json"
	"simple-crud-clean-architecture/internal/enum"
	"time"

	"github.com/google/uuid"
)

type CreateTransactionRequest struct {
	UserID        int
	CourseUUIDStr string `json:"course_uuid,omitempty"`
	CourseUUID    uuid.UUID
}

type MidtransSnapshotRequest struct {
	OrderID  string `json:"order_id,omitempty"`
	GrossAmt int64  `json:"gross_amount,omitempty"`
	Email    string `json:"email,omitempty"`
}

type SnapshotTokenResponse struct {
	TransactionID string `json:"id,omitempty"`
	Token         string `json:"token"`
}

type MidtransNotifyRequest struct {
	SignatureKey      string `json:"signature_key,omitempty"`
	OrderID           string `json:"order_id,omitempty"`
	StatusCode        string `json:"status_code,omitempty"`
	GrossAmount       string `json:"gross_amount,omitempty"`
	TransactionStatus string `json:"transaction_status,omitempty"`
}

type UpdateTransactionStatus struct {
	TransactionID            uuid.UUID
	Status                   enum.MidtransPaymentStatus
	ExternalCallbackResponse json.RawMessage
}

type TransactionResponse struct {
	TrxID    uuid.UUID              `json:"uuid,omitempty"`
	UserID   uuid.UUID              `json:"user_uuid,omitempty"`
	CourseID uuid.UUID              `json:"course_uuid,omitempty"`
	Amount   float64                `json:"amount,omitempty"`
	Status   enum.TransactionStatus `json:"status,omitempty"`

	SnapToken                string          `json:"snap_token,omitempty"`
	ExternalStatus           string          `json:"external_status,omitempty"`
	ExternalCallbackResponse json.RawMessage `json:"external_callback_response,omitempty"`
	CreatedAt                *time.Time      `json:"created_at,omitempty"`
	UpdatedAt                *time.Time      `json:"updated_at,omitempty"`
}
