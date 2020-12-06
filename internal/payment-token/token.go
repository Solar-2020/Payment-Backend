package paymentToken

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type TokenMaker interface {
	Create(tokenData TokenData) (paymentToken string, err error)
	Parse(token string) (paymentToken TokenData, err error)
}

type tokenMaker struct {
	secret          string
	defaultLifetime time.Duration
}

func NewTokenMaker(secret string, defaultLifetime time.Duration) *tokenMaker {
	return &tokenMaker{
		secret:          secret,
		defaultLifetime: defaultLifetime,
	}
}

func (t *tokenMaker) Create(tokenData TokenData) (paymentToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := jwt.MapClaims{}
	claims["UserID"] = tokenData.UserID
	claims["GroupID"] = tokenData.GroupID
	claims["PostID"] = tokenData.PostID
	claims["PaymentID"] = tokenData.PaymentID
	claims["MethodID"] = tokenData.MethodID
	claims["MethodType"] = tokenData.MethodType
	claims["Value"] = tokenData.Value
	claims["Expire"] = time.Now().UTC().Add(t.defaultLifetime).Unix()

	token.Claims = claims
	paymentToken, err = token.SignedString([]byte(t.secret))

	return
}

func (t *tokenMaker) Parse(token string) (paymentToken TokenData, err error) {
	_, err = jwt.ParseWithClaims(token, &paymentToken, secret(t.secret))

	return
}

func secret(secret string) func(*jwt.Token) (interface{}, error) {
	return func(*jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}
