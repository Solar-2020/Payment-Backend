package main

import (
	"database/sql"
	asapi "github.com/Solar-2020/Account-Backend/pkg/api"
	authapi "github.com/Solar-2020/Authorization-Backend/pkg/api"
	"github.com/Solar-2020/GoUtils/context/session"
	"github.com/Solar-2020/GoUtils/http/errorWorker"
	"github.com/Solar-2020/Payment-Backend/cmd/config"
	"github.com/Solar-2020/Payment-Backend/cmd/handlers"
	paymentHandler "github.com/Solar-2020/Payment-Backend/cmd/handlers/payment"
	"github.com/Solar-2020/Payment-Backend/internal/clients/auth"
	"github.com/Solar-2020/Payment-Backend/internal/clients/group"
	"github.com/Solar-2020/Payment-Backend/internal/clients/money"
	"github.com/Solar-2020/Payment-Backend/internal/services/payment"
	"github.com/Solar-2020/Payment-Backend/internal/storages/paymentStorage"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"os"
	"os/signal"
	"syscall"
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
	if err!= nil {
		log.Fatal().Msg(err.Error())
		return
	}

	groupClient := group.NewClient(config.Config.GroupServiceAddress, config.Config.ServerSecret)

	errorWorker := errorWorker.NewErrorWorker()

	paymentStorage := paymentStorage.NewStorage(postsDB)

	paymentTransport := payment.NewTransport()

	paymentService := payment.NewService(paymentStorage, moneyClient, groupClient)

	paymentHandler := paymentHandler.NewHandler(paymentService, paymentTransport, errorWorker)

	authClient := auth.NewClient(config.Config.AuthServiceAddress, config.Config.ServerSecret)

	middlewares := handlers.NewMiddleware(&log, authClient)

	initServices()

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

func initServices() {
	authService := authapi.AuthClient{
		Addr: config.Config.AuthServiceAddress,
	}
	session.RegisterAuthService(&authService)
	accountService := asapi.AccountClient{
		Addr: config.Config.AccountServiceAddress,
	}
	session.RegisterAccountService(&accountService)
}
