package payment

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Solar-2020/Payment-Backend/cmd/config"
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	models2 "github.com/Solar-2020/Payment-Backend/internal/models"
	paymentToken "github.com/Solar-2020/Payment-Backend/internal/payment-token"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"net/url"
	"strconv"
	"time"
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
	tokenMaker 	   paymentToken.TokenMaker
}

func NewService(paymentStorage paymentStorage, moneyClient moneyClient, groupClient groupClient,
	accountBackend accountBackend, errorWorker errorWorker, tokenMaker paymentToken.TokenMaker) *service {
	return &service{
		paymentStorage: paymentStorage,
		moneyClient:    moneyClient,
		groupClient:    groupClient,
		accountBackend: accountBackend,
		errorWorker:    errorWorker,
		tokenMaker:		tokenMaker,
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
		return paymentPage, s.errorWorker.NewError(fasthttp.StatusBadRequest, ErrorYooMoneyAccountNotExit, err)
	}

	var yoomoneyRequisite models.PaymentMethod

	for _, method := range payment.Methods {
		if method.Type == models.YoomoneyType {
			yoomoneyRequisite = method
			break
		}
	}
	if yoomoneyRequisite.ID == 0 {
		err = s.errorWorker.NewError(fasthttp.StatusBadRequest,
			fmt.Errorf("no yoomoney requisite found for payment " + strconv.Itoa(pay.PaymentID)),
			err)
		return
	}

	yandexPayment := money.Payment{
		To:        yoomoneyRequisite.AccountNumber,
		AmountDue: payment.TotalCost,
		Message:   pay.Message,
	}

	requestID, err := s.moneyClient.CreatePayment(yandexPayment)
	if err != nil {
		return paymentPage, s.errorWorker.NewError(fasthttp.StatusInternalServerError, ErrorCantCreateYooMoneyPayment, err)
	}

	token, err := s.tokenMaker.Create(paymentToken.TokenData{
		UserID:         pay.UserID,
		GroupID:        payment.GroupID,
		PostID:         payment.PostID,
		PaymentID:      pay.PaymentID,
		MethodID:      	yoomoneyRequisite.ID,
		MethodType:     3,
		Value:          payment.TotalCost,
	})
	if err != nil {
		return
	}

	successLink := fmt.Sprintf("https://%s/api/payment/confirm?token=%s", config.Config.DomainName, url.QueryEscape(token))
	paymentPage, err = s.moneyClient.CreatePaymentURLWithSuccess(requestID, successLink)
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

type PaymentToken struct {
	UserID int
	GroupID int
	PostID int
	PaymentID int
	MethodID int
	MethodType int
	Value decimal.Decimal
}

func (s *service) ConfirmYoomoney(token string, user int) (redirectUrl string, err error) {
	redirectUrl = config.Config.MoneyFailURL
	tokenDecoded, err := url.QueryUnescape(token)
	if err != nil {
		return redirectUrl, s.errorWorker.NewError(fasthttp.StatusInternalServerError, errors.New("невалидный токен"), err)
	}
	decoded, err :=s.tokenMaker.Parse(tokenDecoded)
	if err != nil {
		return redirectUrl, s.errorWorker.NewError(fasthttp.StatusInternalServerError, errors.New("невалидный токен"), err)
	}

	paid := models2.PaidCreate{
		PostID:        decoded.PostID,
		GroupID:       decoded.GroupID,
		PaymentID:     decoded.PaymentID,
		PayerID:       user,
		Message:       "",
		RequisiteType: decoded.MethodType,
		RequisiteID:   decoded.MethodID,
		PaidAt:        time.Now(),
		Cost:          decoded.Value,
	}
	err = s.paymentStorage.InsertPaid(paid)
	if err != nil {
		return redirectUrl, s.errorWorker.NewError(fasthttp.StatusInternalServerError, errors.New("не удалось зафиксировать оплату"), err)
	}

	redirectUrl = fmt.Sprintf(config.Config.YoomoneyRedirectSuccess, decoded.GroupID)
	return
}
