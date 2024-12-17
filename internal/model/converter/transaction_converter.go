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

	var userResponse *model.UserResponse
	if transaction.User != nil {
		userResponseValue := DTOUserToResponse(*transaction.User)
		userResponse = &userResponseValue
	}

	var courseResponse *model.CourseResponse
	if transaction.Course != nil {
		courseResponseValue := DTOCourseToResponse(*transaction.Course)
		courseResponse = &courseResponseValue

	}

	var voucherResponse *model.VoucherResponse
	if transaction.Voucher != nil {
		voucherResponseValue := DTOVoucherToResponse(*transaction.Voucher)
		voucherResponse = &voucherResponseValue

	}

	return &model.TransactionResponse{
		TrxID:    transaction.TrxID,
		UserID:   transaction.User.UUID,
		CourseID: transaction.Course.UUID,
		Amount:   transaction.Amount,

		User:    userResponse,
		Course:  courseResponse,
		Voucher: voucherResponse,

		Status: transaction.Status,
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
