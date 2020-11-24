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
	TotalCost      decimal.Decimal `json:"paymentValue"`
	PaymentAccount string          `json:"paymentAccount"`
	Methods []PaymentMethod		   `json:"methods"`
	//PaymentArrays
}

//type Payment

type PaymentType string

const (
	PhoneType    PaymentType = "phone"
	CardType     PaymentType = "card"
	YoomoneyType PaymentType = "yoomoney"
)

type PhonePayment struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type CardPayment struct {
	BankName string `json:"bankName,omitempty"`
	CardNumber string `json:"cardNumber,omitempty"`
	PhonePayment
}

type YoomoneyPayment struct {
	AccountNumber string `json:"yoomoneyAccount,omitempty"`
}

type PaymentMethod struct {
	ID int	`json:"id,omitempty"`
	Owner int	`json:"owner"`
	Type PaymentType `json:"type"`
	PaymentID int `json:"-"`

	CardPayment
	YoomoneyPayment
}