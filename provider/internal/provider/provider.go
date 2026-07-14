package provider

import (
	"context"
	"provider/internal/model"
)

type Provider interface {
	Send(ctx context.Context, req model.ProviderRequest) error
}
