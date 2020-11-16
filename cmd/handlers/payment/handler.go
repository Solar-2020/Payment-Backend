package paymentHandler

import (
	"github.com/valyala/fasthttp"
)

type Handler interface {
	Create(ctx *fasthttp.RequestCtx)
	GetByPostIDs(ctx *fasthttp.RequestCtx)
	Pay(ctx *fasthttp.RequestCtx)
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
		h.handleError(err, ctx)
		return
	}

	paymentsReturn, err := h.paymentService.Create(payments)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.paymentTransport.CreateEncode(ctx, paymentsReturn)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) GetByPostIDs(ctx *fasthttp.RequestCtx) {
	postID, err := h.paymentTransport.GetByPostIDsDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	payments, err := h.paymentService.GetByPostIDs(postID)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.paymentTransport.GetByPostIDsEncode(payments, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) Pay(ctx *fasthttp.RequestCtx) {
	pay, err := h.paymentTransport.PayDecode(ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	paymentPage, err := h.paymentService.Pay(pay)
	if err != nil {
		h.handleError(err, ctx)
		return
	}

	err = h.paymentTransport.PayEncode(paymentPage, ctx)
	if err != nil {
		h.handleError(err, ctx)
		return
	}
}

func (h *handler) handleError(err error, ctx *fasthttp.RequestCtx) {
	err = h.errorWorker.ServeJSONError(ctx, err)
	if err != nil {
		h.errorWorker.ServeFatalError(ctx)
	}
	return
}
