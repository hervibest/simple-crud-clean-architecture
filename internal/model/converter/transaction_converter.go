package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"

	"github.com/midtrans/midtrans-go/snap"
)

func SnapToResponse(transaction_id string, snap *snap.Response) *model.SnapshotTokenResponse {
	return &model.SnapshotTokenResponse{
		TransactionID: transaction_id,
		Token:         snap.Token,
	}
}

func TransactionToResponse(transaction *entity.Transaction) *model.TransactionResponse {
	return &model.TransactionResponse{
		TrxID:    transaction.TrxID,
		UserID:   transaction.User.UUID,
		CourseID: transaction.Course.UUID,
		Amount:   transaction.Amount,
		Status:   transaction.Status,
	}
}

func DTOTransactionToResponse(transaction entity.Transaction) model.TransactionResponse {
	return model.TransactionResponse{
		TrxID:    transaction.TrxID,
		UserID:   transaction.User.UUID,
		CourseID: transaction.Course.UUID,
		Amount:   transaction.Amount,
		Status:   transaction.Status,
	}
}
