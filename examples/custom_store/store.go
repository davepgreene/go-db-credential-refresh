package main

import (
	"errors"
	"net/http"

	"github.com/davepgreene/go-db-credential-refresh/driver"
)

// ResponseHandler
type ResponseHandler func(r *http.Response) (driver.Credentials, error)

// HTTPTestConnectingStore
type HTTPTestConnectingStore struct {
	url     string
	method  string
	headers http.Header
	handler ResponseHandler
}

// NewHTTPTestConnectingStore
func NewHTTPTestConnectingStore(url, method string, headers http.Header, handler ResponseHandler) (*HTTPTestConnectingStore, error) {
	if handler == nil {
		return nil, errors.New("handler must be implemented")
	}

	return &HTTPTestConnectingStore{
		url:     url,
		method:  method,
		headers: headers,
		handler: handler,
	}, nil
}

func (h *HTTPTestConnectingStore) Get() (driver.Credentials, error) {
	return h.Refresh()
}

func (h *HTTPTestConnectingStore) Refresh() (driver.Credentials, error) {
	client := &http.Client{}
	req, err := http.NewRequest(h.method, h.url, nil)
	if err != nil {
		return nil, err
	}

	if h.headers != nil {
		req.Header = h.headers
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return h.handler(resp)
}
