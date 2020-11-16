package paymentStorage

import "github.com/shopspring/decimal"

type Payment struct {
	ID             int             `json:"id"`
	GroupID        int             `json:"groupID"`
	PostID         int             `json:"postID"`
	CreateBy       int             `json:"createBy"`
	TotalCost      decimal.Decimal `json:"totalCost"`
	PaymentAccount string          `json:"paymentAccount"`
}
