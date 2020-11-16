package payment

import (
	"errors"
	"fmt"
	"github.com/Solar-2020/Payment-Backend/internal/clients/group"
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	payment "github.com/Solar-2020/Payment-Backend/internal/storages/paymentStorage"
	"github.com/shopspring/decimal"
)

var (
	maxTotalCost             int64 = 10000
	minTotalCost             int64 = 10
	moneyAccountNumberLength       = 15
)

type service struct {
	paymentStorage paymentStorage
	moneyClient    moneyClient
	groupClient    group.Client
}

func NewService(paymentStorage paymentStorage, moneyClient moneyClient, groupClient group.Client) *service {
	return &service{
		paymentStorage: paymentStorage,
		moneyClient:    moneyClient,
		groupClient:    groupClient,
	}
}

func (s *service) Create(createRequest CreateRequest) (createdPayments []payment.Payment, err error) {
	if err = s.validateCreate(createRequest.Payments); err != nil {
		return
	}

	roleID, err := s.groupClient.GetUserRole(createRequest.CreateBy, createRequest.GroupID)
	if err != nil {
		err = fmt.Errorf("restricted")
		return
	}

	if roleID > 2 {
		return createdPayments, errors.New("permission denied")
	}

	err = s.paymentStorage.InsertPayments(createRequest.Payments, createRequest.GroupID, createRequest.PostID)
	if err != nil {
		return
	}

	return s.paymentStorage.SelectPaymentsByPostID(createRequest.PostID)
}

func (s *service) GetByPostIDs(postIDs []int) (payments []payment.Payment, err error) {
	payments, err = s.paymentStorage.SelectPaymentsByPostsIDs(postIDs)
	if err != nil {
		return
	}

	return
}

func (s *service) Pay(pay Pay) (paymentPage money.PaymentPage, err error) {
	payment, err := s.paymentStorage.SelectPayment(pay.PaymentID)
	if err != nil {
		return paymentPage, err
	}

	yandexPayment := money.Payment{
		To:        payment.PaymentAccount,
		AmountDue: payment.TotalCost,
		Message:   pay.Message,
	}

	requestID, err := s.moneyClient.CreatePayment(yandexPayment)
	if err != nil {
		return paymentPage, err
	}

	paymentPage, err = s.moneyClient.CreatePaymentURL(requestID)
	if err != nil {
		return paymentPage, err
	}

	return
}

func (s *service) validateCreate(payments []payment.Payment) (err error) {
	for _, payment := range payments {
		if len(payment.PaymentAccount) != moneyAccountNumberLength {
			return errors.New("Неверный номер кошелька")
		}

		if payment.TotalCost.GreaterThan(decimal.NewFromInt(maxTotalCost)) || payment.TotalCost.LessThan(decimal.NewFromInt(minTotalCost)) {
			return errors.New("Недопустимое значение суммы оплаты")
		}
	}

	return
}
