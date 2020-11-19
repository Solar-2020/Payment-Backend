package paymentHandler

import (
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	"github.com/Solar-2020/Payment-Backend/internal/services/payment"
	"github.com/Solar-2020/Payment-Backend/internal/storages/paymentStorage"
	"github.com/valyala/fasthttp"
)

type paymentService interface {
	Create(createRequest payment.CreateRequest) (createdPayments []paymentStorage.Payment, err error)
	GetByPostIDs(postIDs []int) (payments []paymentStorage.Payment, err error)
	Pay(pay payment.Pay) (paymentPage money.PaymentPage, err error)
}

type paymentTransport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (payments payment.CreateRequest, err error)
	CreateEncode(ctx *fasthttp.RequestCtx, payments []paymentStorage.Payment) (err error)

	GetByPostIDsDecode(ctx *fasthttp.RequestCtx) (postIDs []int, err error)
	GetByPostIDsEncode(payments []paymentStorage.Payment, ctx *fasthttp.RequestCtx) (err error)

	PayDecode(ctx *fasthttp.RequestCtx) (pay payment.Pay, err error)
	PayEncode(paymentPage money.PaymentPage, ctx *fasthttp.RequestCtx) (err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error)
}
