package extensions

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/openkcm/common-sdk/pkg/commoncfg"

	pb "github.com/envoyproxy/gateway/proto/extension"
	slogctx "github.com/veqryn/slog-context"

	"github.com/openkcm/gateway-extension/api"
	gev1a1 "github.com/openkcm/gateway-extension/api/v1alpha1"
	"github.com/openkcm/gateway-extension/internal/flags"
)

type GatewayExtension struct {
	pb.UnimplementedEnvoyGatewayExtensionServer

	features *commoncfg.FeatureGates

	jwtAuthClustersMu sync.RWMutex
	jwtAuthClusters   map[string]*urlCluster
}

func NewGatewayExtension(features *commoncfg.FeatureGates) *GatewayExtension {
	return &GatewayExtension{
		features:          features,
		jwtAuthClustersMu: sync.RWMutex{},
		jwtAuthClusters:   make(map[string]*urlCluster),
	}
}

// PostHTTPListenerModify allows an extension to make changes to a Listener generated by Envoy Gateway before it is finalized.
// PostHTTPListenerModify is always executed when an extension is loaded. An extension may return nil
// in order to not make any changes to it.
func (s *GatewayExtension) PostHTTPListenerModify(ctx context.Context, req *pb.PostHTTPListenerModifyRequest) (*pb.PostHTTPListenerModifyResponse, error) {
	ctx = slogctx.With(ctx, "envoy-xds-hook", "PostHTTPListenerModify")

	slogctx.Info(ctx, "Calling ...")

	resp := &pb.PostHTTPListenerModifyResponse{
		Listener: req.GetListener(),
	}

	if req.GetPostListenerContext() == nil {
		slogctx.Warn(ctx, "Nil PostListenerContext")
		return resp, nil
	}

	if len(req.GetPostListenerContext().GetExtensionResources()) == 0 {
		slogctx.Info(ctx, "Empty list of extension resources")
		return resp, nil
	}

	resources := make(map[string][]any)

	for _, ext := range req.GetPostListenerContext().GetExtensionResources() {
		var generic api.Generic

		err := json.Unmarshal(ext.GetUnstructuredBytes(), &generic)
		if err != nil {
			slogctx.Error(ctx, "Failed to unmarshal the extension", "error", err)
			continue
		}

		switch generic.Kind {
		case api.JWTProviderKind:
			// Do nothing if the feature gate is set
			if s.features.IsFeatureEnabled(flags.DisableJWTProviderComputation) {
				continue
			}

			switch generic.APIVersion {
			case api.JWTProviderV1Alpha1:
				{
					slogctx.Info(ctx, "Found a resource", "yaml", generic)

					jwtProvider := &gev1a1.JWTProvider{}

					err := json.Unmarshal(ext.GetUnstructuredBytes(), jwtProvider)
					if err != nil {
						slogctx.Error(ctx, "Failed to unmarshal the v1alpha1.JWTProvider CRD", "error", err)
						continue
					}

					_, ok := resources[api.JWTProviderKind]
					if !ok {
						resources[api.JWTProviderKind] = make([]any, 0)
					}

					resources[api.JWTProviderKind] = append(resources[api.JWTProviderKind], jwtProvider)
				}
			}
		}
	}

	for key, ext := range resources {
		switch key {
		case api.JWTProviderKind:
			// Do nothing if the feature gate is set
			if s.features.IsFeatureEnabled(flags.DisableJWTProviderComputation) {
				continue
			}

			err := s.ProcessJWTProviders(ctx, req.GetListener(), ext)
			if err != nil {
				return nil, err
			}
		}
	}

	slogctx.Info(ctx, "Called successfully.")

	return resp, nil
}

// PostTranslateModify allows an extension to modify the clusters and secrets in the xDS config.
// This allows for inserting clusters that may change along with extension specific configuration to be dynamically created rather than
// using custom bootstrap config which would be sufficient for clusters that are static and not prone to have their configurations changed.
// An example of how this may be used is to inject a cluster that will be used by an ext_authz http filter created by the extension.
// The list of clusters and secrets returned by the extension are used as the final list of all clusters and secrets
// PostTranslateModify is always executed when an extension is loaded
func (s *GatewayExtension) PostTranslateModify(ctx context.Context, req *pb.PostTranslateModifyRequest) (*pb.PostTranslateModifyResponse, error) {
	ctx = slogctx.With(ctx, "envoy-xds-hook", "PostTranslateModify")

	resp := &pb.PostTranslateModifyResponse{
		Clusters: req.GetClusters(),
		Secrets:  req.GetSecrets(),
	}
	// Return response with same data if the feature gate is set
	if s.features.IsFeatureEnabled(flags.DisableJWTProviderComputation) {
		return resp, nil
	}

	slogctx.Info(ctx, "Calling ...")

	clusters, err := s.TranslateModifyClusters(ctx, req.GetClusters())
	if err != nil {
		return nil, err
	}

	slogctx.Info(ctx, "Called successfully.")

	resp.Clusters = clusters

	return resp, nil
}

// PostVirtualHostModify provides a way for extensions to modify a VirtualHost generated by Envoy Gateway before it is finalized.
// An extension can also make use of this hook to generate and insert entirely new Routes not generated by Envoy Gateway.
// PostVirtualHostModify is always executed when an extension is loaded. An extension may return nil to not make any changes
// to it.
func (s *GatewayExtension) PostVirtualHostModify(ctx context.Context, req *pb.PostVirtualHostModifyRequest) (*pb.PostVirtualHostModifyResponse, error) {
	ctx = slogctx.With(ctx, "envoy-xds-hook", "PostVirtualHostModify")

	slogctx.Info(ctx, "Calling ...")

	resp := &pb.PostVirtualHostModifyResponse{
		VirtualHost: req.GetVirtualHost(),
	}

	// Return response with same data if the feature gate is set
	if s.features.IsFeatureEnabled(flags.DisableJWTProviderComputation) {
		return resp, nil
	}

	if req.GetVirtualHost() == nil {
		slogctx.Warn(ctx, "Nil VirtualHost")
		return resp, nil
	}

	err := s.VirtualHostModifyRoutes(ctx, req.GetVirtualHost().GetRoutes())
	if err != nil {
		return nil, err
	}

	slogctx.Info(ctx, "Called successfully.")

	return resp, nil
}
