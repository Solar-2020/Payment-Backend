package payment

import (
	"encoding/json"
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
	"github.com/valyala/fasthttp"
)

type transport struct {
}

func NewTransport() *transport {
	return &transport{}
}

func (t transport) CreateDecode(ctx *fasthttp.RequestCtx) (payments models.CreateRequest, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &payments)

	return
}

func (t transport) CreateEncode(ctx *fasthttp.RequestCtx, payments []models.Payment) (err error) {
	body, err := json.Marshal(payments)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) GetByPostIDsDecode(ctx *fasthttp.RequestCtx) (postIDs []int, err error) {
	ids := struct {
		PostIDs []int `json:"postIDs"`
	}{}
	err = json.Unmarshal(ctx.Request.Body(), &ids)
	if err != nil {
		return
	}

	return ids.PostIDs, err
}

func (t transport) GetByPostIDsEncode(payments []models.Payment, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(payments)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) PayDecode(ctx *fasthttp.RequestCtx) (pay Pay, err error) {
	err = json.Unmarshal(ctx.Request.Body(), &pay)

	return
}

func (t transport) PayEncode(paymentPage money.PaymentPage, ctx *fasthttp.RequestCtx) (err error) {
	body, err := json.Marshal(paymentPage)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}
