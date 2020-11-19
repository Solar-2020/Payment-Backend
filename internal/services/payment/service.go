package payment

import (
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	payment "github.com/Solar-2020/Payment-Backend/internal/storages/paymentStorage"
	"github.com/valyala/fasthttp"
)

var (
	maxTotalCost             int64 = 10000
	minTotalCost             int64 = 10
	moneyAccountNumberLength       = 15
)

type service struct {
	paymentStorage paymentStorage
	moneyClient    moneyClient
	groupClient    groupClient
	errorWorker    errorWorker
}

func NewService(paymentStorage paymentStorage, moneyClient moneyClient, groupClient groupClient, errorWorker errorWorker) *service {
	return &service{
		paymentStorage: paymentStorage,
		moneyClient:    moneyClient,
		groupClient:    groupClient,
		errorWorker:    errorWorker,
	}
}

func (s *service) Create(createRequest CreateRequest) (createdPayments []payment.Payment, err error) {
	if err = s.validateCreate(createRequest.Payments); err != nil {
		return
	}

	err = s.groupClient.CheckPermission(createRequest.CreateBy, createRequest.GroupID, CreatePaymentActionID)
	if err != nil {
		return
	}

	err = s.paymentStorage.InsertPayments(createRequest.Payments, createRequest.CreateBy, createRequest.GroupID, createRequest.PostID)
	if err != nil {
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}

	createdPayments, err = s.paymentStorage.SelectPaymentsByPostID(createRequest.PostID)
	if err != nil {
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	return
}

func (s *service) GetByPostIDs(postIDs []int) (payments []payment.Payment, err error) {
	payments, err = s.paymentStorage.SelectPaymentsByPostsIDs(postIDs)
	if err != nil {
		err = s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	return
}

func (s *service) Pay(pay Pay) (paymentPage money.PaymentPage, err error) {
	payment, err := s.paymentStorage.SelectPayment(pay.PaymentID)
	if err != nil {
		return paymentPage, s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	yandexPayment := money.Payment{
		To:        payment.PaymentAccount,
		AmountDue: payment.TotalCost,
		Message:   pay.Message,
	}

	requestID, err := s.moneyClient.CreatePayment(yandexPayment)
	if err != nil {
		return paymentPage, s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	paymentPage, err = s.moneyClient.CreatePaymentURL(requestID)
	if err != nil {
		return paymentPage, s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	return
}

func (s *service) validateCreate(payments []payment.Payment) (err error) {
	//for _, payment := range payments {
		//if len(payment.PaymentAccount) != moneyAccountNumberLength {
		//	return errors.New("Неверный номер кошелька")
		//}

		//if payment.TotalCost.GreaterThan(decimal.NewFromInt(maxTotalCost)) || payment.TotalCost.LessThan(decimal.NewFromInt(minTotalCost)) {
		//	return errors.New("Недопустимое значение суммы оплаты")
		//}
	//}

	return
}
