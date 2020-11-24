package models

import (
	account "github.com/Solar-2020/Account-Backend/pkg/models"
	"github.com/shopspring/decimal"
	"time"
)

type Stat struct {
	Payer account.User `json:"payer"`
	Paid
}

type PaidCreate struct {
	PostID        int             `json:"postID"`
	GroupID       int             `json:"groupID"`
	PaymentID     int             `json:"paymentID"`
	PayerID       int             `json:"-"`
	Message       string          `json:"message"`
	RequisiteType int             `json:"requisiteType"`
	RequisiteID   int             `json:"requisiteID"`
	PaidAt        time.Time       `json:"-"`
	Cost          decimal.Decimal `json:"cost"`
}

type Paid struct {
	PaidID        int             `json:"paidID"`
	PaymentID     int             `json:"paymentID"`
	PayerID       int             `json:"-"`
	Message       string          `json:"message"`
	RequisiteType int             `json:"requisiteType"`
	RequisiteID   int             `json:"-"`
	PaidAt        time.Time       `json:"paidAt"`
	Cost          decimal.Decimal `json:"cost"`
	Requisite     Requisite       `json:"requisite"`
}

type Requisite struct {
	*BankCard        `json:"bankCard,omitempty"`
	*YouMoneyAccount `json:"youMoneyAccount,omitempty"`
	*PhonePayment    `json:"phonePayment,omitempty"`
}

type BankCard struct {
	ID          int    `json:"id"`
	BankTitle   string `json:"bankTitle"`
	PhoneNumber string `json:"phoneNumber"`
	CardNumber  string `json:"cardNumber"`
	Owner       int    `json:"owner"`
}

type YouMoneyAccount struct {
	ID            int    `json:"id"`
	AccountNumber string `json:"accountNumber"`
	Owner         int    `json:"owner"`
}

type PhonePayment struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	Owner       int    `json:"owner"`
}
