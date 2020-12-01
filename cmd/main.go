package main

import (
	"database/sql"
	account "github.com/Solar-2020/Account-Backend/pkg/client"
	auth "github.com/Solar-2020/Authorization-Backend/pkg/client"
	"github.com/Solar-2020/GoUtils/http/errorWorker"
	group "github.com/Solar-2020/Group-Backend/pkg/client"
	"github.com/Solar-2020/Payment-Backend/cmd/config"
	"github.com/Solar-2020/Payment-Backend/cmd/handlers"
	"github.com/Solar-2020/Payment-Backend/cmd/handlers/middleware"
	paymentHandler "github.com/Solar-2020/Payment-Backend/cmd/handlers/payment"
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	paymentToken "github.com/Solar-2020/Payment-Backend/internal/payment-token"
	"github.com/Solar-2020/Payment-Backend/internal/services/payment"
	"github.com/Solar-2020/Payment-Backend/internal/storages/paymentStorage"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	err := envconfig.Process("", &config.Config)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	postsDB, err := sql.Open("postgres", config.Config.PaymentDataBaseConnectionString)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	postsDB.SetMaxIdleConns(5)
	postsDB.SetMaxOpenConns(10)

	moneyClient, err := money.NewClient(config.Config.MoneyClientID, config.Config.MoneySuccessURL, config.Config.MoneyFailURL)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	accountClient := account.NewClient(config.Config.AccountServiceHost, config.Config.ServerSecret)
	groupClient := group.NewClient(config.Config.GroupServiceHost, config.Config.ServerSecret)

	errorWorker := errorWorker.NewErrorWorker()

	paymentStorage := paymentStorage.NewStorage(postsDB)

	paymentTransport := payment.NewTransport()
	jwtTokenMaker := paymentToken.NewTokenMaker(config.Config.JWTPaymentTokenSecret,
		time.Duration(config.Config.JWTPaymentTokenLifetime) * time.Second)

	paymentService := payment.NewService(paymentStorage, moneyClient, groupClient, accountClient, errorWorker, jwtTokenMaker)

	paymentHandler := paymentHandler.NewHandler(paymentService, paymentTransport, errorWorker)

	authClient := auth.NewClient(config.Config.AuthServiceHost, config.Config.ServerSecret)

	middlewares := middleware.NewMiddleware(&log, authClient)

	server := fasthttp.Server{
		Handler: handlers.NewFastHttpRouter(paymentHandler, middlewares).Handler,
	}

	go func() {
		log.Info().Str("msg", "start server").Str("port", config.Config.Port).Send()
		if err := server.ListenAndServe(":" + config.Config.Port); err != nil {
			log.Error().Str("msg", "server run failure").Err(err).Send()
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	defer func(sig os.Signal) {

		log.Info().Str("msg", "received signal, exiting").Str("signal", sig.String()).Send()

		if err := server.Shutdown(); err != nil {
			log.Error().Str("msg", "server shutdown failure").Err(err).Send()
		}

		//dbConnection.Shutdown()
		log.Info().Str("msg", "goodbye").Send()
	}(<-c)
}
