package business

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/openkcm/common-sdk/pkg/health"
	"github.com/samber/oops"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	pb "github.com/envoyproxy/gateway/proto/extension"
	slogctx "github.com/veqryn/slog-context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/openkcm/gateway-extension/internal/config"
	"github.com/openkcm/gateway-extension/internal/extensions"
)

const (
	Five = 5 * time.Second
)

// createGRPCServer creates the gRPC server using the given config
func createGRPCServer(_ context.Context, _ *config.Config) *grpc.Server {
	// Create the gRPC server options
	var opts []grpc.ServerOption

	// https://github.com/kubeshop/tracetest/tree/main/examples/quick-start-grpc-stream-propagation/consumer-worker
	// https://github.com/kubeshop/tracetest/tree/main/examples/quick-start-grpc-stream-propagation/producer-api
	opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	return grpc.NewServer(opts...)
}

// StartGRPCServer starts the gRPC server using the given config.
func StartGRPCServer(ctx context.Context, cfg *config.Config) error {
	// Create the gRPC server
	grpcServer := createGRPCServer(ctx, cfg)

	// Register the servers with the gRPC server
	pb.RegisterEnvoyGatewayExtensionServer(grpcServer, extensions.NewGatewayExtension())
	healthpb.RegisterHealthServer(grpcServer, &health.GRPCServer{})

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

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), Five)
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
		if cfg.Listener.UNIX == nil {
			cfg.Listener.UNIX = &config.UNIX{SocketPath: "/etc/envoy/gateway/extension.sock"}
		}
		socketPath := cfg.Listener.UNIX.SocketPath
		// remove old socket
		err := os.Remove(socketPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, oops.Wrapf(err, "Failed to remove unix socket file %s", socketPath)
		}

		return net.Listen("unix", socketPath)
	case config.TCPListener:
		if cfg.Listener.TCP == nil {
			cfg.Listener.TCP = &config.TCP{
				Address: ":9092",
			}
		}
		return net.Listen("tcp", cfg.Listener.TCP.Address)
	}

	return nil, oops.New("Something is wrong, this error should be never popup! TCP or UNIX configuration is required")
}
