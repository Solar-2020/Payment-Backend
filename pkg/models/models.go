package models

import "github.com/shopspring/decimal"

type CreateRequest struct {
	CreateBy int       `json:"createBy"`
	GroupID  int       `json:"groupID"`
	PostID   int       `json:"postID"`
	Payments []Payment `json:"payments"`
}

type Payment struct {
	ID             int             `json:"id"`
	GroupID        int             `json:"groupID"`
	PostID         int             `json:"postID"`
	CreateBy       int             `json:"createBy"`
	TotalCost      decimal.Decimal `json:"totalCost"`
	PaymentAccount string          `json:"paymentAccount"`
}
