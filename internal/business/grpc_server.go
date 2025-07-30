package business

import (
	"context"
	"net"
	"os"

	"github.com/openkcm/common-sdk/pkg/commongrpc"
	"github.com/samber/oops"

	pb "github.com/envoyproxy/gateway/proto/extension"
	slogctx "github.com/veqryn/slog-context"

	"github.com/openkcm/gateway-extension/internal/config"
	"github.com/openkcm/gateway-extension/internal/extensions"
)

// StartGRPCServer starts the gRPC server using the given config.
func StartGRPCServer(ctx context.Context, cfg *config.Config) error {
	// Create the gRPC server
	grpcServer := commongrpc.NewServer(ctx, &cfg.Listener.TCP)

	// Register the servers with the gRPC server

	pb.RegisterEnvoyGatewayExtensionServer(grpcServer, extensions.NewGatewayExtension(&cfg.FeatureGates))

	// Create the listener
	listener, err := createListener(cfg)
	if err != nil {
		return oops.In("TCP GatewayExtension").
			WithContext(ctx).
			Wrapf(err, "Failed to create the listener")
	}

	// Serve the gRPC server
	go func() {
		slogctx.Info(ctx, "Starting gRPC GatewayExtension", "address", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			slogctx.Error(ctx, "ErrorField serving gRPC endpoint", "error", err)
		}

		slogctx.Info(ctx, "Stopped gRPC server")
	}()

	<-ctx.Done()

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, cfg.Listener.ShutdownTimeout)
	defer shutdownRelease()

	grpcServer.Stop()
	slogctx.Info(shutdownCtx, "Completed graceful shutdown of gRPC server")

	return nil
}

func createListener(cfg *config.Config) (net.Listener, error) {
	if cfg.Listener.Type == "" {
		cfg.Listener.Type = config.TCPListener
	}

	switch cfg.Listener.Type {
	case config.UNIXListener:
		socketPath := cfg.Listener.UNIX.SocketPath
		// remove old socket
		err := os.Remove(socketPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, oops.Wrapf(err, "Failed to remove unix socket file %s", socketPath)
		}

		return net.Listen("unix", socketPath)
	case config.TCPListener:
		return net.Listen("tcp", cfg.Listener.TCP.Address)
	}

	return nil, oops.New("Something is wrong, this error should be never popup! TCP or UNIX configuration is required")
}
