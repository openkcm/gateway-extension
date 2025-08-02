package extensions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/utils/ptr"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	jwtauth3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	slogctx "github.com/veqryn/slog-context"

	"github.com/openkcm/gateway-extension/api/v1alpha1"
)

const (
	JwtAuthSecureMappingName = "jwt_auth_secure"
)

// ProcessJWTProviders is called after Envoy Gateway is done generating a
// Listener xDS configuration and before that configuration is passed on to
// Envoy Proxy.
func (s *GatewayExtension) ProcessJWTProviders(ctx context.Context, listener *listenerv3.Listener, resources []any) error {
	providers := make(map[string]*jwtauth3.JwtProvider)
	reqMap := make(map[string]*jwtauth3.JwtRequirement)

	// Collect all jwt providers
	slogctx.Info(ctx, "Processing JWTProviders", "number", len(resources))

	reqs := []*jwtauth3.JwtRequirement{}

	s.jwtAuthClustersMu.Lock()
	defer s.jwtAuthClustersMu.Unlock()

	for _, resource := range resources {
		jwtp, ok := resource.(*v1alpha1.JWTProvider)
		if !ok {
			continue
		}

		slogctx.Info(ctx, "Processing JWTProvider resource...", "name", jwtp.Name)

		jwksTimeoutSec := int64(2)
		jwksCacheDuration := int64(5 * 60)

		var jwksUri string
		if jwtp.Spec.RemoteJwks != nil {
			jwksUri = jwtp.Spec.RemoteJwks.URI

			if jwtp.Spec.RemoteJwks.TimeoutSec > 0 {
				jwksTimeoutSec = jwtp.Spec.RemoteJwks.TimeoutSec
			}

			if jwtp.Spec.RemoteJwks.CacheDuration > 0 {
				jwksCacheDuration = jwtp.Spec.RemoteJwks.CacheDuration
			}
		} else {
			uri, err := extractJWKSFromWellKnownOpenIDConfiguration(ctx, jwtp.Spec.Issuer)
			if err != nil {
				return err
			}

			jwksUri = uri
		}

		if _, err := url.Parse(jwksUri); err != nil {
			slogctx.Error(ctx, "Failed to parse the remote Jwks uri", "error", err)
			continue
		}

		urlCLuster, err := url2Cluster(jwksUri)
		if err != nil {
			slogctx.Error(ctx, "Failed to translate url to cluster", "error", err)
			continue
		}

		remoteJwks := &jwtauth3.RemoteJwks{
			HttpUri: &corev3.HttpUri{
				Uri: jwksUri,
				HttpUpstreamType: &corev3.HttpUri_Cluster{
					Cluster: urlCLuster.CustomName(),
				},
				Timeout: &durationpb.Duration{Seconds: jwksTimeoutSec},
			},
			CacheDuration: &durationpb.Duration{Seconds: jwksCacheDuration},
			AsyncFetch:    &jwtauth3.JwksAsyncFetch{},
		}
		// Set the retry policy if it exists.
		if jwtp.Spec.RemoteJwks != nil && jwtp.Spec.RemoteJwks.Retry != nil {
			var rp *corev3.RetryPolicy
			if rp, err = buildNonRouteRetryPolicy(jwtp.Spec.RemoteJwks.Retry); err != nil {
				return err
			}

			remoteJwks.RetryPolicy = rp
		}

		jwt := &jwtauth3.JwtProvider{
			Issuer:            jwtp.Spec.Issuer,
			Audiences:         jwtp.Spec.Audiences,
			RequireExpiration: jwtp.Spec.RequireExpiration,
			JwksSourceSpecifier: &jwtauth3.JwtProvider_RemoteJwks{
				RemoteJwks: remoteJwks,
			},
			PayloadInMetadata: jwtp.Spec.Name,
			Forward:           true,
			NormalizePayloadInMetadata: &jwtauth3.JwtProvider_NormalizePayload{
				// Normalize the scopes to facilitate matching in Authorization.
				SpaceDelimitedClaims: []string{"scope"},
			},
		}

		if jwtp.Spec.RecomputeRoute != nil {
			jwt.ClearRouteCache = *jwtp.Spec.RecomputeRoute
		}

		if len(jwtp.Spec.FromHeaders) > 0 {
			jwt.FromHeaders = buildJwtFromHeaders(jwtp.Spec.FromHeaders)
		}

		if len(jwtp.Spec.ClaimToHeaders) > 0 {
			jwt.ClaimToHeaders = buildJwtClaimToHeader(jwtp.Spec.ClaimToHeaders)
		}

		if jwtp.Spec.ExtractFrom != nil {
			jwt.FromHeaders = buildJwtFromHeaders(jwtp.Spec.ExtractFrom.Headers)
			jwt.FromCookies = jwtp.Spec.ExtractFrom.Cookies
			jwt.FromParams = jwtp.Spec.ExtractFrom.Params
		}

		providers[jwtp.Spec.Name] = jwt
		reqs = append(reqs, &jwtauth3.JwtRequirement{
			RequiresType: &jwtauth3.JwtRequirement_ProviderName{
				ProviderName: jwtp.Spec.Name,
			},
		})
		s.jwtAuthClusters[urlCLuster.name] = urlCLuster

		slogctx.Info(ctx, "Processed JWTProvider resource", "name", jwtp.Name)
	}

	var jwtRequirement *jwtauth3.JwtRequirement
	switch len(reqs) {
	case 0:
		jwtRequirement = &jwtauth3.JwtRequirement{
			RequiresType: &jwtauth3.JwtRequirement_AllowMissingOrFailed{
				AllowMissingOrFailed: &emptypb.Empty{},
			},
		}
	case 1:
		jwtRequirement = reqs[0]
	default:
		jwtRequirement = &jwtauth3.JwtRequirement{
			RequiresType: &jwtauth3.JwtRequirement_RequiresAny{
				RequiresAny: &jwtauth3.JwtRequirementOrList{
					Requirements: reqs,
				},
			},
		}
	}

	reqMap[JwtAuthSecureMappingName] = jwtRequirement

	// First, get the filter chains from the listener
	filterChains := listener.GetFilterChains()

	defaultFC := listener.GetDefaultFilterChain()
	if defaultFC != nil {
		filterChains = append(filterChains, defaultFC)
	}
	// Go over all the chains, and add the basic authentication http filter
	for _, currChain := range filterChains {
		httpConManager, hcmIndex, err := findHCM(currChain)
		if err != nil {
			slogctx.Warn(ctx, "Failed to find an HCM in the current chain", "filter-chain", currChain.GetName())
			continue
		}

		slogctx.Info(ctx, "Processing HTTPConnectionManager", "index", hcmIndex)
		// If a jwt authentication filter already exists, update it. Otherwise, create it.
		jwtAuthFilter, baIndex, err := findJwtAuthenticationFilter(httpConManager.GetHttpFilters())
		if err != nil {
			slogctx.Error(ctx, "Failed to unmarshal the existing jwtAuthFilter filter",
				"name", currChain.GetName(), "error", err)

			continue
		}

		if baIndex == -1 {
			// Create a new jwt auth filter
			jwtAuthFilter = &jwtauth3.JwtAuthentication{
				Providers:      providers,
				RequirementMap: reqMap,
			}
		} else {
			// Update the jwt auth filter
			jwtAuthFilter.Providers = providers
			jwtAuthFilter.RequirementMap = reqMap
		}
		// Add or update the Jwt Authentication filter in the HCM
		anyFilterConfig, _ := anypb.New(jwtAuthFilter)
		if baIndex > -1 {
			httpConManager.HttpFilters[baIndex].ConfigType = &hcm.HttpFilter_TypedConfig{
				TypedConfig: anyFilterConfig,
			}
		} else {
			filters := []*hcm.HttpFilter{
				{
					Name: egv1a1.EnvoyFilterJWTAuthn.String(),
					ConfigType: &hcm.HttpFilter_TypedConfig{
						TypedConfig: anyFilterConfig,
					},
				},
			}
			filters = append(filters, httpConManager.GetHttpFilters()...)
			httpConManager.HttpFilters = filters
		}

		// Write the updated HCM back to the filter chain
		anyConnectionMgr, _ := anypb.New(httpConManager)
		currChain.Filters[hcmIndex].ConfigType = &listenerv3.Filter_TypedConfig{
			TypedConfig: anyConnectionMgr,
		}

		slogctx.Info(ctx, "Processed HTTPConnectionManager", "index", hcmIndex, "name", currChain.GetName())
	}

	return nil
}

// Tries to find the JWT Authentication HTTP filter in the provided chain
func findJwtAuthenticationFilter(chain []*hcm.HttpFilter) (*jwtauth3.JwtAuthentication, int, error) {
	for i, filter := range chain {
		if filter.GetName() == "envoy.filters.http.jwt_authn" {
			ba := new(jwtauth3.JwtAuthentication)
			err := filter.GetTypedConfig().UnmarshalTo(ba)
			if err != nil {
				return nil, -1, err
			}

			return ba, i, nil
		}
	}

	return nil, -1, nil
}

// buildJwtFromHeaders returns a list of JwtHeader transformed from JWTFromHeader struct
func buildJwtFromHeaders(headers []*v1alpha1.JWTHeader) []*jwtauth3.JwtHeader {
	jwtHeaders := make([]*jwtauth3.JwtHeader, 0, len(headers))

	for _, header := range headers {
		jwtHeader := &jwtauth3.JwtHeader{
			Name:        header.Name,
			ValuePrefix: ptr.Deref(header.ValuePrefix, ""),
		}

		jwtHeaders = append(jwtHeaders, jwtHeader)
	}

	return jwtHeaders
}

// buildJwtFromHeaders returns a list of JwtHeader transformed from JWTFromHeader struct
func buildJwtClaimToHeader(headers []*v1alpha1.JWTClaimToHeader) []*jwtauth3.JwtClaimToHeader {
	jwtHeaders := make([]*jwtauth3.JwtClaimToHeader, 0, len(headers))

	for _, header := range headers {
		jwtHeader := &jwtauth3.JwtClaimToHeader{
			HeaderName: header.HeaderName,
			ClaimName:  header.ClaimName,
		}

		jwtHeaders = append(jwtHeaders, jwtHeader)
	}

	return jwtHeaders
}

const (
	retryDefaultRetryOn = "connect-failure,refused-stream,unavailable,cancelled,retriable-status-codes"
)

func buildNonRouteRetryPolicy(rr *v1alpha1.Retry) (*corev3.RetryPolicy, error) {
	retryOn := rr.RetryOn
	if retryOn == "" {
		retryOn = retryDefaultRetryOn
	}

	rp := &corev3.RetryPolicy{
		RetryOn: retryOn,
	}

	if rr.BackOff != nil {
		baseInterval := rr.BackOff.BaseIntervalSec
		if baseInterval <= 0 {
			baseInterval = 1
		}

		maxInterval := rr.BackOff.MaxIntervalSec
		if maxInterval <= 0 {
			maxInterval = 1
		}

		rp.RetryBackOff = &corev3.BackoffStrategy{
			BaseInterval: &durationpb.Duration{
				Seconds: baseInterval,
			},
			MaxInterval: &durationpb.Duration{
				Seconds: maxInterval,
			},
		}
	}

	if rr.NumRetries != nil {
		retries := *rr.NumRetries
		if retries == 0 {
			retries = 2
		}

		rp.NumRetries = &wrapperspb.UInt32Value{
			Value: retries,
		}
	}

	return rp, nil
}

type wellKnownOpenIDConfiguration struct {
	Issuer string `json:"issuer"`
	JURIS  string `json:"jwks_uri"`
}

func extractJWKSFromWellKnownOpenIDConfiguration(ctx context.Context, issuer string) (string, error) {
	wkoc := wellKnownOpenIDConfiguration{}

	parsedURL, err := url.Parse(issuer)
	if err != nil {
		return "", err
	}

	wkocURI := parsedURL.JoinPath(".well-known/openid-configuration")

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, wkocURI.String(), nil)
	if err != nil {
		return "", fmt.Errorf("could not build request to get well known OpenID configuration: %w", err)
	}

	client := http.DefaultClient

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("could not get well known OpenID configuration: %w", err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			slogctx.Error(ctx, "could not close response body", "error", err)
		}
	}()

	// decode the well known OpenID configuration
	if err := json.NewDecoder(response.Body).Decode(&wkoc); err != nil {
		return "", fmt.Errorf("could not decode well known OpenID configuration: %w", err)
	}

	return wkoc.JURIS, nil
}
