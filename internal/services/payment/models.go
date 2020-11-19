package payment

import (
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	payment "github.com/Solar-2020/Payment-Backend/internal/storages/paymentStorage"
)

const (
	GetPaymentActionID    = 9
	CreatePaymentActionID = 10
	EditPaymentActionID   = 11
	DeletePaymentActionID = 12
)

type paymentStorage interface {
	InsertPayments(payments []payment.Payment, createBy, groupID, postID int) (err error)
	SelectPaymentsByPostsIDs(postIDs []int) (payments []payment.Payment, err error)
	SelectPaymentsByPostID(postID int) (payments []payment.Payment, err error)
	SelectPayment(paymentID int) (payment payment.Payment, err error)
}

type moneyClient interface {
	GetInstanceID() (instanceID string, err error)
	CreatePayment(yandexPayment money.Payment) (requestID string, err error)
	CreatePaymentURL(requestID string) (paymentPage money.PaymentPage, err error)
}

type groupClient interface {
	CheckPermission(userID, groupId, actionID int) (err error)
}

type errorWorker interface {
	NewError(httpCode int, responseError error, fullError error) (err error)
}

type CreateRequest struct {
	CreateBy int               `json:"createBy"`
	GroupID  int               `json:"groupID"`
	PostID   int               `json:"postID"`
	Payments []payment.Payment `json:"payments"`
}

type Pay struct {
	PaymentID int    `json:"paymentID"`
	Message   string `json:"message"`
}
