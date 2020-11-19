package payment

import (
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
)

const (
	GetPaymentActionID    = 9
	CreatePaymentActionID = 10
	EditPaymentActionID   = 11
	DeletePaymentActionID = 12
)

type paymentStorage interface {
	InsertPayments(payments []models.Payment, createBy, groupID, postID int) (err error)
	SelectPaymentsByPostsIDs(postIDs []int) (payments []models.Payment, err error)
	SelectPaymentsByPostID(postID int) (payments []models.Payment, err error)
	SelectPayment(paymentID int) (payment models.Payment, err error)
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

type Pay struct {
	PaymentID int    `json:"paymentID"`
	Message   string `json:"message"`
}
