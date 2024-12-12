package enum

type TransactionStatus string

const (
	TransactionStatusSuccess  TransactionStatus = "SUCCESS"
	TransactionStatusFailed   TransactionStatus = "FAILED"
	TransactionStatusExpired  TransactionStatus = "EXPIRED"
	TransactionStatusCanceled TransactionStatus = "CANCELLED"
	TransactionStatusRefund   TransactionStatus = "REFUND"
	TransactionStatusPending  TransactionStatus = "PENDING"
)
