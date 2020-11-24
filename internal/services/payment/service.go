package payment

import (
	"database/sql"
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	models2 "github.com/Solar-2020/Payment-Backend/internal/models"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
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
	accountBackend accountBackend
	errorWorker    errorWorker
}

func NewService(paymentStorage paymentStorage, moneyClient moneyClient, groupClient groupClient, accountBackend accountBackend, errorWorker errorWorker) *service {
	return &service{
		paymentStorage: paymentStorage,
		moneyClient:    moneyClient,
		groupClient:    groupClient,
		accountBackend: accountBackend,
		errorWorker:    errorWorker,
	}
}

func (s *service) Create(createRequest models.CreateRequest) (createdPayments []models.Payment, err error) {
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

func (s *service) GetByPostIDs(postIDs []int) (payments []models.Payment, err error) {
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

func (s *service) validateCreate(payments []models.Payment) (err error) {
	//for _, payment := range payments {
	//	//if len(payment.PaymentAccount) != moneyAccountNumberLength {
	//	//	return errors.New("Неверный номер кошелька")
	//	//}
	//
	//	//if payment.TotalCost.GreaterThan(decimal.NewFromInt(maxTotalCost)) || payment.TotalCost.LessThan(decimal.NewFromInt(minTotalCost)) {
	//	//	return errors.New("Недопустимое значение суммы оплаты")
	//	//}
	//}

	return
}

func (s *service) Stats(paymentID int) (stats []models2.Stat, err error) {
	stats = make([]models2.Stat, 0)
	paids, err := s.paymentStorage.SelectPaids(paymentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return stats, nil
		}
		return stats, s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	for _, paid := range paids {
		var tempStat models2.Stat
		tempStat.Paid = paid
		tempStat.Payer, err = s.accountBackend.GetUserByUid(tempStat.PayerID)
		if err != nil {
			return stats, s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		}
		stats = append(stats, tempStat)
	}

	return
}

func (s *service) Paid(paidCreate models2.PaidCreate) (err error) {
	err = s.paymentStorage.InsertPaid(paidCreate)
	if err != nil {
		return s.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}
	return
}
