package payment

import (
	"encoding/json"
	"errors"
	"github.com/Solar-2020/Payment-Backend/cmd/config"
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	models2 "github.com/Solar-2020/Payment-Backend/internal/models"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
	"github.com/valyala/fasthttp"
	"strconv"
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
	pay.UserID = ctx.UserValue("userID").(int)
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

func (t transport) StatsDecode(ctx *fasthttp.RequestCtx) (paymentID int, err error) {
	paymentIDStr := ctx.UserValue("paymentID").(string)
	paymentID, err = strconv.Atoi(paymentIDStr)
	if err != nil {
		return
	}
	return
}

func (t transport) StatsEncode(ctx *fasthttp.RequestCtx, stats []models2.Stat) (err error) {
	body, err := json.Marshal(stats)
	if err != nil {
		return
	}
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return
}

func (t transport) PaidDecode(ctx *fasthttp.RequestCtx) (paidCreate models2.PaidCreate, err error) {
	userID := ctx.UserValue("userID").(int)
	err = json.Unmarshal(ctx.Request.Body(), &paidCreate)
	if err != nil {
		return
	}
	paidCreate.PayerID = userID
	return
}

func (t transport) PaidEncode(ctx *fasthttp.RequestCtx) (err error) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	return
}

func (t transport) ConfirmYoomoneyDecode(ctx *fasthttp.RequestCtx) (token string, uid int, err error) {
	uid = ctx.UserValue("userID").(int)
	//token = ctx.UserValue("token").(string)
	tokenBytes := ctx.QueryArgs().Peek("token")
	if tokenBytes == nil {
		err = errors.New("token: empty")
		return
	}
	token = string(tokenBytes)
	return
}

func (t transport) ConfirmYoomoneyEncode(ctx *fasthttp.RequestCtx) (err error) {
	ctx.Redirect(config.Config.YoomoneyRedirectSuccess, fasthttp.StatusTemporaryRedirect)
	return
}