package config

import "github.com/Solar-2020/GoUtils/common"

var (
	Config config
)

type config struct {
	common.SharedConfig
	PaymentDataBaseConnectionString string `envconfig:"PAYMENT_DB_CONNECTION_STRING" default:"-"`
	DomainName                      string `envconfig:"DOMAIN_NAME" default:"solar.ru"` //for static file prefix
	ServerSecret                    string `envconfig:"SERVER_SECRET" default:"Basic secret"`
	MoneyClientID                   string `envconfig:"MONEY_CLIENT_ID"`
	MoneySuccessURL                 string `envconfig:"MONEY_SUCCESS_URL"`
	MoneyFailURL                    string `envconfig:"MONEY_FAIL_URL"`
}
