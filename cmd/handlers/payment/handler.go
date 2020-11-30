package paymentHandler

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	GetByPostIDs(ctx *fasthttp.RequestCtx)
	Pay(ctx *fasthttp.RequestCtx)
	Paid(ctx *fasthttp.RequestCtx)
	Stats(ctx *fasthttp.RequestCtx)
	ConfirmYoomoney(ctx *fasthttp.RequestCtx)
}

type handler struct {
	paymentService   paymentService
	paymentTransport paymentTransport
	errorWorker      errorWorker
}

func NewHandler(paymentService paymentService, paymentTransport paymentTransport, errorWorker errorWorker) Handler {
	return &handler{
		paymentService:   paymentService,
		paymentTransport: paymentTransport,
		errorWorker:      errorWorker,
	}
}

func (h *handler) Create(ctx *fasthttp.RequestCtx) {
	payments, err := h.paymentTransport.CreateDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	paymentsReturn, err := h.paymentService.Create(payments)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentTransport.CreateEncode(ctx, paymentsReturn)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) GetByPostIDs(ctx *fasthttp.RequestCtx) {
	postID, err := h.paymentTransport.GetByPostIDsDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	payments, err := h.paymentService.GetByPostIDs(postID)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentTransport.GetByPostIDsEncode(payments, ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Pay(ctx *fasthttp.RequestCtx) {
	pay, err := h.paymentTransport.PayDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	paymentPage, err := h.paymentService.Pay(pay)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentTransport.PayEncode(paymentPage, ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Stats(ctx *fasthttp.RequestCtx) {
	paymentID, err := h.paymentTransport.StatsDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	stats, err := h.paymentService.Stats(paymentID)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentTransport.StatsEncode(ctx, stats)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) Paid(ctx *fasthttp.RequestCtx) {
	createPayment, err := h.paymentTransport.PaidDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentService.Paid(createPayment)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentTransport.PaidEncode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}
}

func (h *handler) ConfirmYoomoney(ctx *fasthttp.RequestCtx) {
	token, uid, err := h.paymentTransport.ConfirmYoomoneyDecode(ctx)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentService.ConfirmYoomoney(token, uid)
	if err != nil {
		h.errorWorker.ServeJSONError(ctx, err)
		return
	}

	err = h.paymentTransport.ConfirmYoomoneyEncode(ctx)
	//if err != nil {
	//	h.errorWorker.ServeJSONError(ctx, err)
	//	return
	//}
	return
}
