package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	paymentHandler "github.com/Solar-2020/Payment-Backend/cmd/handlers/payment"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(payment paymentHandler.Handler, middleware Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.PanicHandler = httputils.PanicHandler

	router.Handle("GET", "/health", middleware.Log(httputils.HealthCheckHandler))

	router.Handle("POST", "/api/payment/payment", middleware.Log(middleware.ExternalAuth(payment.Create)))
	router.Handle("GET", "/api/payment/payment", middleware.Log(middleware.ExternalAuth(payment.GetByPostIDs)))
	router.Handle("POST", "/api/payment/pay", middleware.Log(middleware.ExternalAuth(payment.Pay)))

	return router
}
