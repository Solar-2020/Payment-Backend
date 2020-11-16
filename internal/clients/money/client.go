package money

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"net/url"
)

type Client interface {
	GetInstanceID() (instanceID string, err error)
	CreatePayment(yandexPayment Payment) (requestID string, err error)
	CreatePaymentURL(requestID string) (paymentPage PaymentPage, err error)
}

type client struct {
	clientID   string
	instanceID string
	successURL string
	failURL    string
}

func NewClient(clientID, successURL, failURL string) Client {
	return &client{clientID: clientID, successURL: successURL, failURL: failURL}
}

func (c *client) GetInstanceID() (instanceID string, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("https")
	req.URI().SetHost("money.yandex.ru")
	req.URI().SetPath("api/instance-id")

	data := url.Values{}
	data.Set("client_id", c.clientID)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		var response struct {
			InstanceID string `json:"instance_id"`
			Status     string `json:"status"`
			Error      string `json:"error"`
		}

		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			return
		}

		if response.Status != "success" {
			return "", errors.New(response.Error)
		}

		return response.InstanceID, nil
	default:
		return "", errors.New("Unexpected Server Error")
	}
}

func (c *client) CreatePayment(yandexPayment Payment) (requestID string, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("https")
	req.URI().SetHost("money.yandex.ru")
	req.URI().SetPath("api/request-external-payment")

	data := url.Values{}
	data.Set("pattern_id", "p2p")
	data.Set("instance_id", c.instanceID)
	data.Set("to", yandexPayment.To)
	data.Set("amount_due", yandexPayment.AmountDue.String())
	data.Set("message", yandexPayment.Message)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		var response struct {
			RequestID string `json:"request_id"`
			Status    string `json:"status"`
			Error     string `json:"error"`
		}

		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			return
		}

		if response.Status != "success" {
			return "", errors.New(response.Error)
		}

		return response.RequestID, nil
	default:
		return "", errors.New("Unexpected Server Error")
	}
}

func (c *client) CreatePaymentURL(requestID string) (paymentPage PaymentPage, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("https")
	req.URI().SetHost("money.yandex.ru")
	req.URI().SetPath("api/process-external-payment")

	data := url.Values{}
	data.Set("request_id", requestID)
	data.Set("instance_id", c.instanceID)
	data.Set("ext_auth_success_uri", c.successURL)
	data.Set("ext_auth_fail_uri", c.failURL)
	data.Set("request_token", "true")

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		var response struct {
			AcsURI    string `json:"acs_uri"`
			AcsParams struct {
				OrderId string `json:"orderId"`
			} `json:"acs_params"`
			Status string `json:"status"`
			Error  string `json:"error"`
		}

		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			return
		}

		if response.Status != "success" {
			return paymentPage, errors.New(response.Error)
		}

		paymentPage.URL = response.AcsURI
		paymentPage.OrderID = response.AcsParams.OrderId

		return paymentPage, nil
	default:
		return paymentPage, errors.New("Unexpected Server Error")
	}
}
