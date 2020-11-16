package handlers

import (
	"github.com/Solar-2020/GoUtils/log"
	"github.com/Solar-2020/Payment-Backend/internal/clients/auth"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"time"
)

type Middleware interface {
	Log(next fasthttp.RequestHandler) fasthttp.RequestHandler
	ExternalAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler
	InternalAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler
}

type middleware struct {
	log        *zerolog.Logger
	authClient auth.Client
}

func NewMiddleware(log *zerolog.Logger, authClient auth.Client) Middleware {
	return &middleware{
		log:        log,
		authClient: authClient,
	}
}

func (m middleware) Log(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		logger := log.NewLog()
		log.Set(ctx, &logger)
		logger.Println(ctx, "Start new request: ", ctx.Request.URI())
		if len(ctx.Request.String()) < 1024 {
			logger.Println(ctx, ctx.Request.String())
		}

		defer func(begin time.Time) {
			logger.Printf(
				ctx,
				"End: %s, status: %d, time: %d ms",
				ctx.Request.URI().String(),
				ctx.Response.StatusCode(),
				time.Since(begin).Milliseconds(),
			)
		}(time.Now())

		next(ctx)
	}
}

func (m middleware) ExternalAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		cookie := ctx.Request.Header.Cookie("SessionToken")

		userID, err := m.authClient.GetUserIDByCookie(string(cookie))
		if err != nil {
			ctx.SetUserValue("error", err)
			ServeUnAuthorizationError(ctx)
			return
		}
		ctx.SetUserValue("userID", userID)
		next(ctx)
	}
}

func (m middleware) InternalAuth(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authSecret := ctx.Request.Header.Peek("Authorization")

		err := m.authClient.CompareSecret(string(authSecret))
		if err != nil {
			ctx.SetUserValue("error", err)
			ServeUnAuthorizationError(ctx)
			return
		}
		next(ctx)
	}
}

func ServeUnAuthorizationError(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetStatusCode(fasthttp.StatusUnauthorized)
	return
}
