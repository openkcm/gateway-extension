package business

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/openkcm/common-sdk/pkg/commoncfg"
	"github.com/stretchr/testify/assert"

	"github.com/openkcm/gateway-extension/internal/config"
)

func TestStartGRPCServer(t *testing.T) {
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
					TCP:  commoncfg.GRPCServer{Address: ":0"},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Start with Unix listener",
			cfg: &config.Config{
				Listener: config.Listener{
					Type: config.UNIXListener,
					UNIX: config.UNIX{
						SocketPath: "extension.sock",
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Unexpected listener type",
			cfg: &config.Config{
				Listener: config.Listener{
					Type: "Unexpected",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Default to TCP listener",
			cfg: &config.Config{
				Listener: config.Listener{
					Type: "",
					TCP:  commoncfg.GRPCServer{},
					UNIX: config.UNIX{},
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
				errCh <- StartGRPCServer(ctx, tt.cfg)
			}()

			time.Sleep(1 * time.Second)
			cancel()
			err := <-errCh
			tt.wantErr(t, err, fmt.Sprintf("StartGRPCServer(ctx, %v)", tt.cfg))
		})
	}
}
