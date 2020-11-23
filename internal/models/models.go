package models

import (
	account "github.com/Solar-2020/Account-Backend/pkg/models"
	"github.com/shopspring/decimal"
	"time"
)

type Stat struct {
	Payer account.User
	Paid
}

type Paid struct {
	PaidID        int             `json:"paidID"`
	PaymentID     int             `json:"paymentID"`
	PayerID       int             `json:"-"`
	RequisiteType int             `json:"requisiteType"`
	RequisiteID   int             `json:"-"`
	PaidAt        time.Time       `json:"paidAt"`
	Cost          decimal.Decimal `json:"cost"`
	Requisite     `json:"requisite"`
}

type Requisite struct {
	*BankCard        `json:"bankCard,requisite"`
	*YouMoneyAccount `json:"youMoneyAccount,requisite"`
	*PhonePayment    `json:"phonePayment,requisite"`
}

type BankCard struct {
	ID          int    `json:"requisite"`
	BankTitle   string `json:"requisite"`
	PhoneNumber string `json:"requisite"`
	CardNumber  string `json:"requisite"`
	Owner       int    `json:"requisite"`
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
