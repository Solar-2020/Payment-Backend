package payment

import (
	payment "github.com/Solar-2020/Payment-Backend/internal/storages/paymentStorage"
)

type paymentStorage interface {
	InsertPayments(payments []payment.Payment, groupID, postID int) (err error)
	SelectPayments(postIDs []int) (payments []payment.Payment, err error)
	SelectPayment(paymentID int) (payment payment.Payment, err error)
}

type CreateRequest struct {
	CreateBy int               `json:"id"`
	GroupID  int               `json:"groupID"`
	PostID   int               `json:"postID"`
	Payments []payment.Payment `json:"payments"`
}

type Pay struct {
	PaymentID int    `json:"paymentID"`
	Message   string `json:"message"`
}
