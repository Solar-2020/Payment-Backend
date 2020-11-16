package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	paymentHandler "github.com/Solar-2020/Payment-Backend/cmd/handlers/payment"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(payment paymentHandler.Handler, middleware Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.PanicHandler = httputils.PanicHandler

	router.Handle("GET", "/health", httputils.HealthCheckHandler)

	//router.Handle("POST", "/api/payment/pay", middleware.Log(middleware.ExternalAuth(payment.Pay)))
	router.Handle("POST", "/api/payment/pay", payment.Pay)

	router.Handle("POST", "/api/internal/payment/payment", middleware.Log(middleware.InternalAuth(payment.Create)))
	router.Handle("POST", "/api/internal/payment/by-post-ids", middleware.Log(middleware.InternalAuth(payment.GetByPostIDs)))

	return router
}
