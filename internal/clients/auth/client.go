package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Client interface {
	GetUserIDByCookie(sessionToken string) (userID int, err error)
	CompareSecret(inputSecret string) (err error)
}

type client struct {
	host   string
	secret string
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret}
}

type CheckAuthRequest struct {
	SessionToken string `json:"cookie"`
}

type httpError struct {
	Error string `json:"error"`
}

type userIDRequest struct {
	UserID int `json:"uid"`
}

func (c *client) GetUserIDByCookie(sessionToken string) (userID int, err error) {
	checkAuthRequest := CheckAuthRequest{SessionToken: sessionToken}
	body, err := json.Marshal(checkAuthRequest)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, c.host+"/api/auth/cookie", bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response userIDRequest
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response.UserID, err
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return userID, errors.New(httpErr.Error)
	default:
		return userID, errors.New("Unexpected Server Error")
	}
}

func (c *client) CompareSecret(inputSecret string) (err error) {
	if !strings.EqualFold(inputSecret, c.secret) {
		return errors.New("Invalid server secret")
	}
	return
}
