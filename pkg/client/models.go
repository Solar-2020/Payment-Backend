package client

var (
	ErrorUnknownStatusCode = "Unknown status code %v"
)

type httpError struct {
	Error string `json:"error"`
}
