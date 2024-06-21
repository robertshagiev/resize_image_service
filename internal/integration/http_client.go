package integration

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
)

// NewHTTPClient создает и возвращает нового HTTP-клиента с заданными настройками.
func NewHTTPClient(timeout time.Duration) *resty.Client {
	client := resty.New()
	client.SetTimeout(timeout)
	// Здесь можно добавить другие настройки клиента, если это необходимо.
	return client
}

// NewRequestContext создает новый GET-запрос с контекстом.
func NewRequestContext(ctx context.Context, client *resty.Client, url string) (*resty.Request, error) {
	req := client.R().SetContext(ctx).SetQueryParam("url", url)
	return req, nil
}
