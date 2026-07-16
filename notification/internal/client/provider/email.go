package provider

import (
	"context"
	"notification/internal/model"
)

func (c *Client) SendEmail(ctx context.Context, req model.ProviderRequest) (model.NotificationResult, error) {
	err := c.do(ctx, req, "email")

	if err != nil {
		return model.NotificationResult{}, err
	}

	return model.NotificationResult{}, nil
}
