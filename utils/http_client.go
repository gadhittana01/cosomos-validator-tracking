package utils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gadhittana01/cosmos-validation-tracking/utils/types"
)

type HTTPClient interface {
	Get(ctx context.Context, url string) (*types.HTTPResponse, error)
	Post(ctx context.Context, url string, jsonBody []byte) (*types.HTTPResponse, error)
	Put(ctx context.Context, url string, jsonBody []byte) (*types.HTTPResponse, error)
	Patch(ctx context.Context, url string, jsonBody []byte) (*types.HTTPResponse, error)
	Delete(ctx context.Context, url string) (*types.HTTPResponse, error)
}

type DefaultHTTPClient struct {
	client *http.Client
}

func NewDefaultHTTPClient() HTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *DefaultHTTPClient) Get(ctx context.Context, url string) (*types.HTTPResponse, error) {
	return c.doRequestWithContext(ctx, http.MethodGet, url, nil)
}

func (c *DefaultHTTPClient) Post(ctx context.Context, url string, jsonBody []byte) (*types.HTTPResponse, error) {
	return c.doRequestWithContext(ctx, http.MethodPost, url, jsonBody)
}

func (c *DefaultHTTPClient) Put(ctx context.Context, url string, jsonBody []byte) (*types.HTTPResponse, error) {
	return c.doRequestWithContext(ctx, http.MethodPut, url, jsonBody)
}

func (c *DefaultHTTPClient) Patch(ctx context.Context, url string, jsonBody []byte) (*types.HTTPResponse, error) {
	return c.doRequestWithContext(ctx, http.MethodPatch, url, jsonBody)
}

func (c *DefaultHTTPClient) Delete(ctx context.Context, url string) (*types.HTTPResponse, error) {
	return c.doRequestWithContext(ctx, http.MethodDelete, url, nil)
}

func (c *DefaultHTTPClient) doRequestWithContext(ctx context.Context, method, url string, body []byte) (*types.HTTPResponse, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &types.HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
		Headers:    resp.Header,
	}, nil
}
