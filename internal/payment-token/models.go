package paymentToken

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/shopspring/decimal"
)

type TokenData struct {
	jwt.StandardClaims
	UserID     int
	GroupID    int
	PostID     int
	PaymentID  int
	MethodID   int
	MethodType int
	Value      decimal.Decimal
}
