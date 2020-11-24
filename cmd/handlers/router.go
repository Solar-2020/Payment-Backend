package handlers

import (
	httputils "github.com/Solar-2020/GoUtils/http"
	"github.com/Solar-2020/Payment-Backend/cmd/handlers/middleware"
	paymentHandler "github.com/Solar-2020/Payment-Backend/cmd/handlers/payment"
	"github.com/buaazp/fasthttprouter"
)

func NewFastHttpRouter(payment paymentHandler.Handler, middleware middleware.Middleware) *fasthttprouter.Router {
	router := fasthttprouter.New()

	router.PanicHandler = httputils.PanicHandler

	router.Handle("GET", "/health", httputils.HealthCheckHandler)

	//router.Handle("POST", "/api/payment/pay", middleware.Log(middleware.ExternalAuth(payment.Pay)))
	router.Handle("POST", "/api/payment/pay", payment.Pay)
	router.Handle("POST", "/api/payment/paid", middleware.Log(middleware.ExternalAuth(payment.Paid)))
	router.Handle("GET", "/api/payment/stat/:paymentID", middleware.Log(middleware.ExternalAuth(payment.Stats)))
	//router.Handle("GET", "/api/payment/stat/:paymentID", payment.Stats)

	router.Handle("POST", "/api/internal/payment/payment", middleware.Log(middleware.InternalAuth(payment.Create)))
	router.Handle("POST", "/api/internal/payment/by-post-ids", middleware.Log(middleware.InternalAuth(payment.GetByPostIDs)))



	return router
}
