package email

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"provider/internal/model"
	"provider/internal/requestid"
	"time"
)

var ErrSendFailed = errors.New("send email failed")

type Provider struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Provider {
	return &Provider{
		logger: logger,
	}
}

func (p *Provider) loggerWithContext(ctx context.Context) *slog.Logger {
	return p.logger.With(
		slog.String("request_id", requestid.Get(ctx)),
	)
}

func (p *Provider) Send(
	ctx context.Context,
	req model.ProviderRequest,
) error {
	logger := p.loggerWithContext(ctx)
	logger.Info(
		"sending email",
		slog.Int("user_id", req.UserID),
		slog.String("message", req.Message),
	)

	select {
	case <-ctx.Done():
		return ctx.Err()

	case <-time.After(250 * time.Millisecond):
		if rand.Intn(100) > 85 {
			return ErrSendFailed
		}
	}

	logger.Info(
		"email sent",
		slog.Int("user_id", req.UserID),
		slog.String("message", req.Message),
	)

	return nil
}
