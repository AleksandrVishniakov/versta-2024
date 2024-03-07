package apiclient

import (
	"context"
	"net/http"
	"time"
)

type BaseClient interface {
	Send(req *http.Request) (*http.Response, error)
}

type APIClient struct {
	ctx    context.Context
	client http.Client
}

func NewAPIClient(ctx context.Context) *APIClient {
	return &APIClient{
		ctx: ctx,
		client: http.Client{
			Timeout: 7 * time.Second,
		},
	}
}

func (a *APIClient) Send(req *http.Request) (*http.Response, error) {
	req = req.WithContext(a.ctx)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
