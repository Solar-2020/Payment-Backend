package group

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type Client interface {
	GetUserRole(userID, groupID int) (roleID int, err error)
}

type client struct {
	host   string
	secret string
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret}
}

type httpError struct {
	Error string `json:"error"`
}

type UserRole struct {
	UserID   int    `json:"userID"`
	GroupID  int    `json:"groupID"`
	RoleID   int    `json:"roleID"`
	RoleName string `json:"roleName"`
}

func (c *client) GetUserRole(userID, groupID int) (roleID int, err error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/group/permission")

	req.URI().QueryArgs().Add("user_id", strconv.Itoa(userID))
	req.URI().QueryArgs().Add("group_id", strconv.Itoa(groupID))

	req.Header.Set("Authorization", c.secret)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		var response UserRole
		err = json.Unmarshal(resp.Body(), &response)
		return response.RoleID, err
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return
		}
		return roleID, errors.New(httpErr.Error)
	default:
		return roleID, errors.New("Unexpected Server Error")
	}
}

func (c *client) CompareSecret(inputSecret string) (err error) {
	if !strings.EqualFold(inputSecret, c.secret) {
		return errors.New("Invalid server secret")
	}
	return
}
