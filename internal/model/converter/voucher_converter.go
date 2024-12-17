package converter

import (
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
)

func VoucherToResponse(voucher *entity.Voucher) *model.VoucherResponse {

	courseResponse := Map(voucher.Courses, DTOCourseToResponse)
	transactionResponse := Map(voucher.Transactions, DTOTransactionToResponse)

	return &model.VoucherResponse{
		UUID:          voucher.UUID,
		Name:          voucher.Name,
		Code:          voucher.Code,
		Type:          voucher.Type,
		Value:         voucher.Value,
		IsActive:      voucher.IsActive,
		StartActiveAt: voucher.StartActiveAt,
		ValidUntil:    voucher.ValidUntil,
		Courses:       courseResponse,
		Transaction:   transactionResponse,
		CreatedAt:     voucher.CreatedAt,
		UpdatedAt:     voucher.UpdatedAt,
	}
}

func DTOVoucherToResponse(voucher entity.Voucher) model.VoucherResponse {
	return model.VoucherResponse{
		UUID:          voucher.UUID,
		Name:          voucher.Name,
		Code:          voucher.Code,
		Type:          voucher.Type,
		Value:         voucher.Value,
		IsActive:      voucher.IsActive,
		StartActiveAt: voucher.StartActiveAt,
		ValidUntil:    voucher.ValidUntil,
	}
}
