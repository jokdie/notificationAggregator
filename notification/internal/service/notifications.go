package service

import (
	"context"
	"notification/internal/model"
)

type Client interface {
	SendEmail(ctx context.Context, req model.ProviderRequest) (model.NotificationResult, error)
	SendPush(ctx context.Context, req model.ProviderRequest) (model.NotificationResult, error)
	SendSms(ctx context.Context, req model.ProviderRequest) (model.NotificationResult, error)
}

type Notifications struct {
	client Client
}

func NewNotifications(client Client) *Notifications {
	return &Notifications{client: client}
}

func (s *Notifications) Send(ctx context.Context, req model.NotificationRequest) (model.NotificationResult, error) {
	reqDTO := model.ProviderRequest{UserID: req.UserID, Message: req.Message}
	_, _ = s.client.SendEmail(ctx, reqDTO)

	return model.NotificationResult{}, nil
}
