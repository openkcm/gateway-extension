package extensions

import (
	"context"

	"google.golang.org/protobuf/types/known/anypb"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	jwtauthnv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	slogctx "github.com/veqryn/slog-context"
)

func (s *GatewayExtension) VirtualHostModifyRoutes(ctx context.Context, routes []*routev3.Route) error {
	for _, r := range routes {
		slogctx.Info(ctx, "Updated VirtualHost Route", "name", r.GetName())

		filterCfg := r.GetTypedPerFilterConfig()
		if _, ok := filterCfg[egv1a1.EnvoyFilterJWTAuthn.String()]; !ok {
			routeCfgProto := &jwtauthnv3.PerRouteConfig{
				RequirementSpecifier: &jwtauthnv3.PerRouteConfig_RequirementName{RequirementName: JwtAuthSecureMappingName},
			}

			routeCfgAny, err := anypb.New(routeCfgProto)
			if err != nil {
				return err
			}

			if filterCfg == nil {
				r.TypedPerFilterConfig = make(map[string]*anypb.Any)
			}

			r.TypedPerFilterConfig[egv1a1.EnvoyFilterJWTAuthn.String()] = routeCfgAny
		}
	}

	return nil
}
