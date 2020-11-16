package money

import (
	"github.com/shopspring/decimal"
)

type Payment struct {
	PatternID  string          `json:"pattern_id"`
	InstanceID string          `json:"instance_id"`
	To         string          `json:"to"`
	AmountDue  decimal.Decimal `json:"amount_due"`
	Message    string          `json:"message"`
}

type PaymentURLRequest struct {
	RequestID    string `json:"request_id"`
	InstanceID   string `json:"instance_id"`
	SuccessURI   string `json:"ext_auth_success_uri"`
	FailURI      string `json:"ext_auth_fail_uri"`
	RequestToken bool   `json:"request_token"`
}

type PaymentPage struct {
	URL     string `json:"url"`
	OrderID string `json:"orderID"`
}
