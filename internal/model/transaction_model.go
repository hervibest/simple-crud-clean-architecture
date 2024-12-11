package model

import (
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
