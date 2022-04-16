package main

import (
	"context"
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

func (h *HTTPTestConnectingStore) Get(ctx context.Context) (driver.Credentials, error) {
	return h.Refresh(ctx)
}

func (h *HTTPTestConnectingStore) Refresh(ctx context.Context) (driver.Credentials, error) {
	client := &http.Client{}
	req, err := http.NewRequest(h.method, h.url, nil)
	if err != nil {
		return nil, err
	}

	if h.headers != nil {
		req.Header = h.headers
	}

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return h.handler(resp)
}
