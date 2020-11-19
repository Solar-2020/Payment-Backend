package middleware

import (
	"github.com/Solar-2020/GoUtils/http/errorWorker"
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
	authClient authClient
}

func NewMiddleware(log *zerolog.Logger, authClient authClient) Middleware {
	return &middleware{
		log:        log,
		authClient: authClient,
	}
}

func (m middleware) Log(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func(begin time.Time) {
			execTime := time.Since(begin).Milliseconds()
			if ctx.Value("error") != nil {
				responseError := ctx.Value("error").(errorWorker.ResponseError)
				errLog := m.log.Error()
				errLog.
					Time("request start", begin).
					Int64("request duration", execTime).
					Str("request", string(ctx.Request.Body())).
					Int("code", ctx.Response.StatusCode()).
					Str("front msg", string(ctx.Response.Body())).
					Str("full error", responseError.FullError().Error())
				errLog.Send()
			}
			// else m.log.Info() если нет ошибки и хочешь залоггировать все запросы. Реквест можно вытащить и в конце запроса,
			// если нужны доп поля, то их всегда можно в хендлере положить контекст и вытащить тут из него
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
