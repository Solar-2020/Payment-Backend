package money

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"net/url"
)

type client struct {
	clientID   string
	instanceID string
	successURL string
	failURL    string
}

func NewClient(clientID, successURL, failURL string) (newClient *client, err error) {
	newClient = &client{clientID: clientID, successURL: successURL, failURL: failURL}
	newClient.instanceID, err = newClient.GetInstanceID()
	if err != nil {
		return
	}
	return
}

func (c *client) GetInstanceID() (instanceID string, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("https")
	req.URI().SetHost("money.yandex.ru")
	req.URI().SetPath("api/instance-id")

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.SetMethod(fasthttp.MethodPost)

	data := url.Values{}
	data.Set("client_id", c.clientID)
	req.SetBodyString(data.Encode())

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

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.SetMethod(fasthttp.MethodPost)

	data := url.Values{}
	data.Set("pattern_id", "p2p")
	data.Set("instance_id", c.instanceID)
	data.Set("to", yandexPayment.To)
	data.Set("amount_due", yandexPayment.AmountDue.String())
	data.Set("message", yandexPayment.Message)
	req.SetBodyString(data.Encode())

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
	return c.CreatePaymentURLParametrized(requestID, c.successURL, c.failURL)
}

func (c *client) CreatePaymentURLWithSuccess(requestID string, successUrl string) (paymentPage PaymentPage, err error) {
	return c.CreatePaymentURLParametrized(requestID, successUrl, c.failURL)
}

func (c *client) CreatePaymentURLParametrized(requestID string, successUrl, failUrl string) (paymentPage PaymentPage, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("https")
	req.URI().SetHost("money.yandex.ru")
	req.URI().SetPath("api/process-external-payment")

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.SetMethod(fasthttp.MethodPost)

	data := url.Values{}
	data.Set("request_id", requestID)
	data.Set("instance_id", c.instanceID)
	data.Set("ext_auth_success_uri", successUrl)
	data.Set("ext_auth_fail_uri", failUrl)
	data.Set("request_token", "true")
	req.SetBodyString(data.Encode())

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

		if response.Status != "ext_auth_required" {
			return paymentPage, errors.New(response.Error)
		}

		paymentPage.URL = response.AcsURI
		paymentPage.OrderID = response.AcsParams.OrderId

		return paymentPage, nil
	default:
		return paymentPage, errors.New("Unexpected Server Error")
	}
}
