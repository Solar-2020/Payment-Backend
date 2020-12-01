package paymentToken

import "github.com/dgrijalva/jwt-go"

type TokenData struct {
	jwt.StandardClaims
	UserID    string
	PaymentID string
}
