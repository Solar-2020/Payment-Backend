package paymentHandler

import (
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	models2 "github.com/Solar-2020/Payment-Backend/internal/models"
	"github.com/Solar-2020/Payment-Backend/internal/services/payment"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
	"github.com/valyala/fasthttp"
)

type paymentService interface {
	Create(createRequest models.CreateRequest) (createdPayments []models.Payment, err error)
	GetByPostIDs(postIDs []int) (payments []models.Payment, err error)
	Pay(pay payment.Pay) (paymentPage money.PaymentPage, err error)
	Stats(paymentID int) (stats []models2.Stat, err error)
}

type paymentTransport interface {
	CreateDecode(ctx *fasthttp.RequestCtx) (payments models.CreateRequest, err error)
	CreateEncode(ctx *fasthttp.RequestCtx, payments []models.Payment) (err error)

	GetByPostIDsDecode(ctx *fasthttp.RequestCtx) (postIDs []int, err error)
	GetByPostIDsEncode(payments []models.Payment, ctx *fasthttp.RequestCtx) (err error)

	PayDecode(ctx *fasthttp.RequestCtx) (pay payment.Pay, err error)
	PayEncode(paymentPage money.PaymentPage, ctx *fasthttp.RequestCtx) (err error)

	StatsDecode(ctx *fasthttp.RequestCtx) (paymentID int, err error)
	StatsEncode(ctx *fasthttp.RequestCtx, stats []models2.Stat) (err error)
}

type errorWorker interface {
	ServeJSONError(ctx *fasthttp.RequestCtx, serveError error)
}
