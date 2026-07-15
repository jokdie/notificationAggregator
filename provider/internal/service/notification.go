package service

import (
	"context"
	"errors"
	"provider/internal/model"
	"provider/internal/provider"
)

var ErrUnknownProvider = errors.New("unknown notification provider")

type NotificationService struct {
	providers map[model.Channel]provider.Provider
}

func NewNotificationService(
	providers map[model.Channel]provider.Provider,
) *NotificationService {
	return &NotificationService{
		providers: providers,
	}
}

func (s *NotificationService) Send(
	ctx context.Context,
	channel model.Channel,
	req model.ProviderRequest,
) error {
	p, ok := s.providers[channel]
	if !ok {
		return ErrUnknownProvider
	}

	return p.Send(ctx, req)
}
