package extensions

import (
	"context"

	"google.golang.org/protobuf/types/known/anypb"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	jwtauthnv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	slogctx "github.com/veqryn/slog-context"

	"github.com/openkcm/gateway-extension/internal/flags"
)

func (s *GatewayExtension) VirtualHostModifyRoutes(ctx context.Context, routes []*routev3.Route) error {
	for _, r := range routes {
		slogctx.Info(ctx, "Updated VirtualHost Route", "name", r.GetName())
		cleanupRoute(ctx, r, JwtAuthSecureMappingName)

		filterCfg := r.GetTypedPerFilterConfig()
		// Do nothing if the feature gate is set making empty the jwt providers
		if s.features.IsFeatureEnabled(flags.DisableJWTProviderComputation) {
			slogctx.Warn(ctx, "Skipping JWTProvider as is disabled through flags", "name", r.GetName())
			return nil
		}

		if filterCfg == nil {
			r.TypedPerFilterConfig = make(map[string]*anypb.Any)
		}

		if _, ok := filterCfg[egv1a1.EnvoyFilterJWTAuthn.String()]; !ok {
			routeCfgProto := &jwtauthnv3.PerRouteConfig{
				RequirementSpecifier: &jwtauthnv3.PerRouteConfig_RequirementName{RequirementName: JwtAuthSecureMappingName},
			}

			routeCfgAny, err := anypb.New(routeCfgProto)
			if err != nil {
				return err
			}

			r.TypedPerFilterConfig[egv1a1.EnvoyFilterJWTAuthn.String()] = routeCfgAny
		}
	}

	return nil
}

func cleanupRoute(ctx context.Context, r *routev3.Route, name string) {
	resources := r.GetTypedPerFilterConfig()
	for key, cfg := range resources {
		if cfg == nil {
			continue
		}

		perRouteConfig := new(jwtauthnv3.PerRouteConfig)

		err := cfg.UnmarshalTo(perRouteConfig)
		if err != nil {
			slogctx.Debug(ctx, "Failed to unmarshal into PerRouteConfig", "source_resource", cfg.String())
			continue
		}

		if name == perRouteConfig.GetRequirementName() {
			delete(resources, key)
		}
	}
}
