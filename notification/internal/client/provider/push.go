package provider

import (
	"context"
	"notification/internal/model"
)

func (c *Client) SendPush(ctx context.Context, req model.ProviderRequest) (model.NotificationResult, error) {
	return model.NotificationResult{}, nil
}
