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
	MoneySuccessURL                 string `envconfig:"MONEY_SUCCESS_URL" default:"https://develop.pay-together.ru/pay/success"`
	MoneyFailURL                    string `envconfig:"MONEY_FAIL_URL" default:"https://develop.pay-together.ru/pay/error"`
	AuthServiceHost                 string `envconfig:"AUTH_SERVICE_HOST" default:"develop.pay-together.ru"`
	GroupServiceHost                string `envconfig:"GROUP_SERVICE_HOST" default:"develop.pay-together.ru"`
	AccountServiceHost              string `envconfig:"ACCOUNT_SERVICE_HOST" default:"develop.pay-together.ru"`
	YoomoneyRedirectSuccess			string `envconfig:"SUCCESS_PAYMENT_REDIRECT" default:"https://develop.pay-together.ru/group/%d"`
	YoomoneyRedirectCookieLifetime  int	   `envconfig:"PAYMENT_REDIRECT_COOKIE_LIFETIEME" default:300"`

	JWTPaymentTokenSecret			string	`envconfig:"JWT_PAYMENT_TOKEN_SECRET" required:"true"`
	JWTPaymentTokenLifetime			int		`envconfig:"JWT_PAYMENT_TOKEN_LIFETIME" default:"600"`	// seconds
}
