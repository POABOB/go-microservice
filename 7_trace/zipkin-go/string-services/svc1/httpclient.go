//go:build go1.7
// +build go1.7

package svc1

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"zipkin-go/middleware"

	"github.com/opentracing/opentracing-go"
)

// client is our actual client implementation
type client struct {
	baseURL      string
	httpClient   *http.Client
	tracer       opentracing.Tracer
	traceRequest middleware.RequestFunc
}

// Concat implements our Service interface.
func (c *client) Concat(ctx context.Context, a, b string) (string, error) {
	// 如果 context 中有 span 代表那是父 span，如果沒有，我們就根 span
	span, ctx := opentracing.StartSpanFromContext(ctx, "Concat")
	defer span.Finish()

	// 建立 HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf(
		"%s/concat/?a=%s&b=%s", c.baseURL, url.QueryEscape(a), url.QueryEscape(b),
	), nil)
	if err != nil {
		return "", err
	}

	// 使用 traceRequest 這個 Middleware 去追蹤
	req = c.traceRequest(req.WithContext(ctx))

	// 執行 HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// 標記為錯誤
		span.SetTag("error", err.Error())
		return "", err
	}
	defer resp.Body.Close()

	// 讀取回應
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 標記為錯誤
		span.SetTag("error", err.Error())
		return "", err
	}

	return string(data), nil
}

// Sum implements our Service interface.
func (c *client) Sum(ctx context.Context, a, b int64) (int64, error) {
	// 如果 context 中有 span 代表那是父 span，如果沒有，我們就根 span
	span, ctx := opentracing.StartSpanFromContext(ctx, "Sum")
	defer span.Finish()

	// 建立 HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/sum/?a=%d&b=%d", c.baseURL, a, b), nil)
	if err != nil {
		return 0, err
	}

	// 使用 traceRequest 這個 Middleware 去追蹤
	req = c.traceRequest(req.WithContext(ctx))

	// 執行 HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// annotate our span with the error condition
		span.SetTag("error", err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	// 讀取回應
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 標記為錯誤
		span.SetTag("error", err.Error())
		return 0, err
	}

	// 將 data 轉換成 int64
	result, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		// 標記為錯誤
		span.SetTag("error", err.Error())
		return 0, err
	}

	return result, nil
}

// 實例化，並挾帶 tracer 和 middleware
func NewHTTPClient(tracer opentracing.Tracer, baseURL string) Service {
	return &client{
		baseURL:      baseURL,
		httpClient:   &http.Client{},
		tracer:       tracer,
		traceRequest: middleware.ToHTTPRequest(tracer),
	}
}
