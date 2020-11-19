package client

import (
	"encoding/json"
	"github.com/Solar-2020/GoUtils/http/errorWorker"
	"github.com/Solar-2020/Payment-Backend/pkg/models"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Client interface {
	Create(createRequest models.CreateRequest) (createdPayments []models.Payment, err error)
	GetByPostIDs(postIDs []int) (payments []models.Payment, err error)
}

type client struct {
	host        string
	secret      string
	errorWorker errorWorker.ErrorWorker
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret, errorWorker: errorWorker.NewErrorWorker()}
}

func (c *client) Create(createRequest models.CreateRequest) (createdPayments []models.Payment, err error) {
	if len(createRequest.Payments) == 0 {
		return
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/payment/payment")

	req.Header.Set("Authorization", c.secret)
	req.Header.SetMethod(fasthttp.MethodPost)

	body, err := json.Marshal(createRequest)
	if err != nil {
		return
	}

	req.SetBody(body)

	err = fasthttp.Do(req, resp)
	if err != nil {
		err = c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		err = json.Unmarshal(resp.Body(), &createdPayments)
		if err != nil {
			err = c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
			return
		}

		return
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			err = c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
			return
		}
		return createdPayments, c.errorWorker.NewError(fasthttp.StatusBadRequest, errors.New(httpErr.Error), errors.New(httpErr.Error))

	default:
		return createdPayments, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, errors.Errorf(ErrorUnknownStatusCode, resp.StatusCode()))
	}
}

func (c *client) GetByPostIDs(postIDs []int) (payments []models.Payment, err error) {
	payments = make([]models.Payment, 0)
	if len(postIDs) == 0 {
		return
	}

	ids := struct {
		PostIDs []int `json:"postIDs"`
	}{PostIDs: postIDs}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/payment/by-post-ids")

	req.Header.Set("Authorization", c.secret)
	req.Header.SetMethod(fasthttp.MethodPost)

	body, err := json.Marshal(ids)
	if err != nil {
		err = c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}

	req.SetBody(body)

	err = fasthttp.Do(req, resp)
	if err != nil {
		err = c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		err = json.Unmarshal(resp.Body(), &payments)
		if err != nil {
			err = c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
			return
		}
		return
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			err = c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
			return
		}
		return payments, c.errorWorker.NewError(fasthttp.StatusBadRequest, errors.New(httpErr.Error), errors.New(httpErr.Error))

	default:
		return payments, c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, errors.Errorf(ErrorUnknownStatusCode, resp.StatusCode()))
	}
}

func (c *client) CheckPermission(userID, groupId, actionID int) (err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/group/check-permission")

	req.URI().QueryArgs().Set("user_id", strconv.Itoa(userID))
	req.URI().QueryArgs().Set("group_id", strconv.Itoa(groupId))
	req.URI().QueryArgs().Set("action_id", strconv.Itoa(actionID))

	req.Header.Set("Authorization", c.secret)
	req.Header.SetMethod(fasthttp.MethodGet)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		return
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		}
		return c.errorWorker.NewError(fasthttp.StatusBadRequest, errors.New(httpErr.Error), errors.New(httpErr.Error))
	case fasthttp.StatusForbidden:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, err)
		}
		return c.errorWorker.NewError(fasthttp.StatusForbidden, errors.New(httpErr.Error), errors.New(httpErr.Error))
	default:
		return c.errorWorker.NewError(fasthttp.StatusInternalServerError, nil, errors.Errorf(ErrorUnknownStatusCode, resp.StatusCode()))
	}
}
