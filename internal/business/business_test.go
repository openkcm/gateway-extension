package business

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/openkcm/gateway-extension/internal/config"
)

func TestMainFunc(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Start with TCP listener",
			cfg: &config.Config{
				Listener: config.Listener{
					Type: config.TCPListener,
					TCP:  &config.TCP{Address: ":0"},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			errCh := make(chan error)
			go func() {
				errCh <- Main(ctx, tt.cfg)
			}()

			time.Sleep(1 * time.Second)
			cancel()
			err := <-errCh
			tt.wantErr(t, err, fmt.Sprintf("StartGRPCServer(ctx, %v)", tt.cfg))
		})
	}
}
