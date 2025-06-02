package business

import (
	"context"

	"github.com/openkcm/gateway-extension/internal/config"
)

// Main Application Business Logic
func Main(ctx context.Context, cfg *config.Config) error {
	return StartGRPCServer(ctx, cfg)
}
