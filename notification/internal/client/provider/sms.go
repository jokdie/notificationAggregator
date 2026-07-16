package provider

import (
	"context"
	"notification/internal/model"
)

func (c *Client) SendSms(ctx context.Context, req model.ProviderRequest) (model.NotificationResult, error) {
	return model.NotificationResult{}, nil
}
