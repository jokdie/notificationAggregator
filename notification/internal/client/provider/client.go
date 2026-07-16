package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"notification/internal/model"
	"notification/internal/requestid"
	"time"

	"golang.org/x/sync/semaphore"
)

type ProviderError struct {
	StatusCode int
	Response   model.ProviderErrorResponse
}

func (e *ProviderError) Error() string {
	return e.Response.Message
}

type Client struct {
	httpClient *http.Client
	sem        *semaphore.Weighted
	logger     *slog.Logger
	baseURL    string
}

func NewClient(
	sem *semaphore.Weighted,
	logger *slog.Logger,
	baseURL string,
) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 2,
		},
		sem:     sem,
		logger:  logger,
		baseURL: baseURL,
	}
}

func (c *Client) do(
	ctx context.Context,
	reqDTO model.ProviderRequest,
	path string,
) error {
	if err := c.sem.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("sem acquire: %w", err)
	}
	defer c.sem.Release(1)

	body, err := json.Marshal(reqDTO)

	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+path,
		bytes.NewReader(body),
	)

	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Request-ID", requestid.Get(ctx))
	start := time.Now()

	logger := c.loggerWithContext(ctx)

	logger.Info(
		"sending request",
		slog.String("method", httpReq.Method),
		slog.String("url", path),
	)

	httpResp, err := c.httpClient.Do(httpReq)

	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer httpResp.Body.Close()

	logger.Info(
		"provider responded",
		slog.Int("status", httpResp.StatusCode),
		slog.Duration("duration", time.Since(start)),
	)

	if httpResp.StatusCode == http.StatusNoContent {
		return nil
	}

	var resp model.ProviderErrorResponse

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return fmt.Errorf("create provider error: %w", err)
	}

	return &ProviderError{
		StatusCode: httpResp.StatusCode,
		Response:   resp,
	}
}

func (c *Client) loggerWithContext(ctx context.Context) *slog.Logger {
	return c.logger.With(
		slog.String("request_id", requestid.Get(ctx)),
	)
}
