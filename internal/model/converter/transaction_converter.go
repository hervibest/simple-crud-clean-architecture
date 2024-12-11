package converter

import (
	"simple-crud-clean-architecture/internal/model"

	"github.com/midtrans/midtrans-go/snap"
)

func SnapToResponse(transaction_id string, snap *snap.Response) *model.SnapshotTokenResponse {
	return &model.SnapshotTokenResponse{
		TransactionID: transaction_id,
		Token:         snap.Token,
	}
}
