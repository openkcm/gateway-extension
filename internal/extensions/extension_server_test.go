package extensions

import (
	"fmt"
	"maps"
	"net/http"
	"testing"
	"time"

	"github.com/envoyproxy/gateway/proto/extension"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	jwtauth3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"

	"github.com/openkcm/gateway-extension/internal/extensions/testdata"
)

func mustNewAny(src proto.Message) *anypb.Any {
	v, err := anypb.New(src)
	if err != nil {
		panic(err)
	}

	return v
}

func TestGatewayExtension_PostHTTPListenerModify(t *testing.T) {
	tests := []struct {
		name    string
		req     *extension.PostHTTPListenerModifyRequest
		want    *extension.PostHTTPListenerModifyResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Modify request",
			req: &extension.PostHTTPListenerModifyRequest{
				Listener: &listenerv3.Listener{
					FilterChains: []*listenerv3.FilterChain{{
						FilterChainMatch:              nil,
						Filters:                       nil,
						Metadata:                      nil,
						TransportSocket:               nil,
						TransportSocketConnectTimeout: nil,
						Name:                          "",
					}},
					DefaultFilterChain: &listenerv3.FilterChain{
						FilterChainMatch: nil,
						Filters: []*listenerv3.Filter{
							{
								Name: wellknown.HTTPConnectionManager,
								ConfigType: &listenerv3.Filter_TypedConfig{
									TypedConfig: mustNewAny(&hcm.HttpConnectionManager{
										HttpFilters: []*hcm.HttpFilter{
											{
												Name: "envoy.filters.http.jwt_authn",
												ConfigType: &hcm.HttpFilter_TypedConfig{
													TypedConfig: mustNewAny(&jwtauth3.JwtAuthentication{
														Providers: map[string]*jwtauth3.JwtProvider{
															"Provider": {
																Issuer:                     "",
																Audiences:                  []string{"one", "two"},
																RequireExpiration:          false,
																JwksSourceSpecifier:        nil,
																Forward:                    true,
																FromHeaders:                nil,
																FromParams:                 nil,
																FromCookies:                nil,
																PayloadInMetadata:          "",
																NormalizePayloadInMetadata: nil,
															},
														},
														RequirementMap: map[string]*jwtauth3.JwtRequirement{
															"jwt_auth_secure": {
																RequiresType: &jwtauth3.JwtRequirement_ProviderName{
																	ProviderName: "Provider",
																},
															},
														},
													}),
												},
												IsOptional: false,
												Disabled:   false,
											},
										},
									}),
								},
							},
						},
					},
				},
				PostListenerContext: &extension.PostHTTPListenerExtensionContext{
					ExtensionResources: []*extension.ExtensionResource{
						{
							UnstructuredBytes: testdata.ExtensionJSON,
						},
					},
				},
			},
			want: &extension.PostHTTPListenerModifyResponse{
				Listener: &listenerv3.Listener{
					FilterChains: []*listenerv3.FilterChain{{
						FilterChainMatch:              nil,
						Filters:                       nil,
						Metadata:                      nil,
						TransportSocket:               nil,
						TransportSocketConnectTimeout: nil,
						Name:                          "",
					}},
					DefaultFilterChain: &listenerv3.FilterChain{
						FilterChainMatch: nil,
						Filters: []*listenerv3.Filter{
							{
								Name: wellknown.HTTPConnectionManager,
								ConfigType: &listenerv3.Filter_TypedConfig{
									TypedConfig: mustNewAny(&hcm.HttpConnectionManager{
										HttpFilters: []*hcm.HttpFilter{
											{
												Name: "envoy.filters.http.jwt_authn",
												ConfigType: &hcm.HttpFilter_TypedConfig{
													TypedConfig: mustNewAny(&jwtauth3.JwtAuthentication{
														Providers: map[string]*jwtauth3.JwtProvider{
															"Provider": {
																Issuer:    "",
																Audiences: []string{"one", "two"},
																ClaimToHeaders: []*jwtauth3.JwtClaimToHeader{{
																	HeaderName: "X-Custom-Header",
																	ClaimName:  "claim",
																}},
																RequireExpiration: false,
																JwksSourceSpecifier: &jwtauth3.JwtProvider_RemoteJwks{
																	RemoteJwks: &jwtauth3.RemoteJwks{
																		HttpUri: &corev3.HttpUri{
																			Uri:              "https://example.com/jwks",
																			HttpUpstreamType: &corev3.HttpUri_Cluster{Cluster: "example_com_443|openkcm"},
																			Timeout:          durationpb.New(2 * time.Second),
																		},
																		AsyncFetch:    &jwtauth3.JwksAsyncFetch{},
																		CacheDuration: durationpb.New(300 * time.Second),
																		RetryPolicy: &corev3.RetryPolicy{
																			RetryBackOff: &corev3.BackoffStrategy{
																				BaseInterval: durationpb.New(1 * time.Second),
																				MaxInterval:  durationpb.New(1 * time.Second),
																			},
																			NumRetries:                    wrapperspb.UInt32(2),
																			RetryOn:                       "connect-failure,refused-stream,unavailable,cancelled,retriable-status-codes",
																			RetryPriority:                 nil,
																			RetryHostPredicate:            nil,
																			HostSelectionRetryMaxAttempts: 0,
																		},
																	},
																},
																Forward:                    true,
																FromHeaders:                []*jwtauth3.JwtHeader{{Name: "X-Custom-Header", ValuePrefix: "prefix"}},
																FromParams:                 []string{"Param one", "Param two"},
																FromCookies:                nil,
																PayloadInMetadata:          "Provider",
																NormalizePayloadInMetadata: &jwtauth3.JwtProvider_NormalizePayload{SpaceDelimitedClaims: []string{"scope"}},
															},
														},
														RequirementMap: map[string]*jwtauth3.JwtRequirement{
															"jwt_auth_secure": {
																RequiresType: &jwtauth3.JwtRequirement_ProviderName{
																	ProviderName: "Provider",
																},
															},
														},
													}),
												},
											},
										},
									}),
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGatewayExtension()
			req := proto.CloneOf(tt.req)
			got, err := s.PostHTTPListenerModify(t.Context(), req)
			if !tt.wantErr(t, err, fmt.Sprintf("PostHTTPListenerModify(%v)", req)) {
				return
			}
			diff := cmp.Diff(tt.want, got, protocmp.Transform(), protocmp.IgnoreDefaultScalars())
			if diff != "" {
				assert.Fail(t, fmt.Sprintf("Not equal: \n"+
					"expected: %s\n"+
					"actual  : %s%s", tt.want, got, diff), "PostHTTPListenerModify(%v)", req)
			}
			if len(s.jwtAuthClusters) == 0 {
				assert.Fail(t, "No jwt auth clusters processed")
			}
		})
	}
}

func startWellKnownServer() {
	if err := http.ListenAndServe(":4543", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if _, err := w.Write(testdata.OpenIDConfigurationJSON); err != nil {
			panic(err)
		}
	})); err != nil {
		panic(err)
	}
}

func TestGatewayExtension_PostHTTPListenerModify_WellKnown(t *testing.T) {
	go startWellKnownServer()

	tests := []struct {
		name    string
		req     *extension.PostHTTPListenerModifyRequest
		want    *extension.PostHTTPListenerModifyResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Well known JWT authentication",
			req: &extension.PostHTTPListenerModifyRequest{
				Listener: &listenerv3.Listener{
					FilterChains: []*listenerv3.FilterChain{{
						FilterChainMatch:              nil,
						Filters:                       nil,
						Metadata:                      nil,
						TransportSocket:               nil,
						TransportSocketConnectTimeout: nil,
						Name:                          "",
					}},
					DefaultFilterChain: &listenerv3.FilterChain{
						FilterChainMatch: nil,
						Filters: []*listenerv3.Filter{
							{
								Name: wellknown.HTTPConnectionManager,
								ConfigType: &listenerv3.Filter_TypedConfig{
									TypedConfig: mustNewAny(&hcm.HttpConnectionManager{
										HttpFilters: []*hcm.HttpFilter{
											{
												Name: "envoy.filters.http.jwt_authn",
												ConfigType: &hcm.HttpFilter_TypedConfig{
													TypedConfig: mustNewAny(&jwtauth3.JwtAuthentication{
														Providers: map[string]*jwtauth3.JwtProvider{
															"Provider": {
																Issuer:                     "",
																Audiences:                  []string{"one", "two"},
																RequireExpiration:          false,
																JwksSourceSpecifier:        nil,
																Forward:                    true,
																FromHeaders:                nil,
																FromParams:                 nil,
																FromCookies:                nil,
																PayloadInMetadata:          "",
																NormalizePayloadInMetadata: nil,
															},
														},
														RequirementMap: map[string]*jwtauth3.JwtRequirement{
															"jwt_auth_secure": {
																RequiresType: &jwtauth3.JwtRequirement_ProviderName{
																	ProviderName: "Provider",
																},
															},
														},
													}),
												},
												IsOptional: false,
												Disabled:   false,
											},
										},
									}),
								},
							},
						},
					},
				},
				PostListenerContext: &extension.PostHTTPListenerExtensionContext{
					ExtensionResources: []*extension.ExtensionResource{
						{
							UnstructuredBytes: testdata.WellKnownJSON,
						},
					},
				},
			},
			want: &extension.PostHTTPListenerModifyResponse{
				Listener: &listenerv3.Listener{
					FilterChains: []*listenerv3.FilterChain{{
						FilterChainMatch:              nil,
						Filters:                       nil,
						Metadata:                      nil,
						TransportSocket:               nil,
						TransportSocketConnectTimeout: nil,
						Name:                          "",
					}},
					DefaultFilterChain: &listenerv3.FilterChain{
						FilterChainMatch: nil,
						Filters: []*listenerv3.Filter{
							{
								Name: wellknown.HTTPConnectionManager,
								ConfigType: &listenerv3.Filter_TypedConfig{
									TypedConfig: mustNewAny(&hcm.HttpConnectionManager{
										HttpFilters: []*hcm.HttpFilter{
											{
												Name: "envoy.filters.http.jwt_authn",
												ConfigType: &hcm.HttpFilter_TypedConfig{
													TypedConfig: mustNewAny(&jwtauth3.JwtAuthentication{
														Providers: map[string]*jwtauth3.JwtProvider{
															"Well Known": {
																Issuer:    "http://localhost:4543",
																Audiences: []string{"one", "two"},
																ClaimToHeaders: []*jwtauth3.JwtClaimToHeader{{
																	HeaderName: "X-Custom-Header",
																	ClaimName:  "claim",
																}},
																RequireExpiration: false,
																JwksSourceSpecifier: &jwtauth3.JwtProvider_RemoteJwks{
																	RemoteJwks: &jwtauth3.RemoteJwks{
																		HttpUri: &corev3.HttpUri{
																			Uri:              "http://www.localhost/oauth2/v3/certs",
																			HttpUpstreamType: &corev3.HttpUri_Cluster{Cluster: "www_localhost_80|openkcm"},
																			Timeout:          durationpb.New(2 * time.Second),
																		},
																		AsyncFetch:    &jwtauth3.JwksAsyncFetch{},
																		CacheDuration: durationpb.New(300 * time.Second),
																	},
																},
																Forward:                    true,
																FromHeaders:                []*jwtauth3.JwtHeader{{Name: "X-Custom-Header", ValuePrefix: "prefix"}},
																FromParams:                 []string{"Param one", "Param two"},
																FromCookies:                nil,
																PayloadInMetadata:          "Well Known",
																NormalizePayloadInMetadata: &jwtauth3.JwtProvider_NormalizePayload{SpaceDelimitedClaims: []string{"scope"}},
															},
														},
														RequirementMap: map[string]*jwtauth3.JwtRequirement{
															"jwt_auth_secure": {
																RequiresType: &jwtauth3.JwtRequirement_ProviderName{
																	ProviderName: "Well Known",
																},
															},
														},
													}),
												},
											},
										},
									}),
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGatewayExtension()
			req := proto.CloneOf(tt.req)
			got, err := s.PostHTTPListenerModify(t.Context(), req)
			if !tt.wantErr(t, err, fmt.Sprintf("PostHTTPListenerModify(%v)", req)) {
				return
			}
			diff := cmp.Diff(tt.want, got, protocmp.Transform(), protocmp.IgnoreDefaultScalars())
			if diff != "" {
				assert.Fail(t, fmt.Sprintf("Not equal: \n"+
					"expected: %s\n"+
					"actual  : %s%s", tt.want, got, diff), "PostHTTPListenerModify(%v)", req)
			}
			if len(s.jwtAuthClusters) == 0 {
				assert.Fail(t, "No jwt auth clusters processed")
			}
		})
	}
}

func TestGatewayExtension_PostTranslateModify(t *testing.T) {
	tests := []struct {
		name            string
		jwtAuthClusters map[string]*urlCluster
		req             *extension.PostTranslateModifyRequest
		want            *extension.PostTranslateModifyResponse
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "Post Translate Modify",
			jwtAuthClusters: map[string]*urlCluster{
				"example_com_443": {
					name:         "example_com_443",
					hostname:     "example.com",
					port:         443,
					endpointType: EndpointTypeDNS,
					tls:          true,
				},
			},
			req: &extension.PostTranslateModifyRequest{
				PostTranslateContext: &extension.PostTranslateExtensionContext{},
				Clusters: []*clusterv3.Cluster{{
					Name: "localhost_80|openkcm",
				}, {
					Name: "www_localhost_80",
				}},
				Secrets: []*tlsv3.Secret{
					{
						Name: "secret",
						Type: &tlsv3.Secret_TlsCertificate{
							TlsCertificate: &tlsv3.TlsCertificate{
								CertificateChain: &corev3.DataSource{
									Specifier: &corev3.DataSource_EnvironmentVariable{EnvironmentVariable: "TLS_CERT"},
								},
								PrivateKey: &corev3.DataSource{
									Specifier: &corev3.DataSource_EnvironmentVariable{EnvironmentVariable: "TLS_KEY"},
								},
							},
						},
					},
				},
			},
			want: &extension.PostTranslateModifyResponse{
				Clusters: []*clusterv3.Cluster{{
					Name: "www_localhost_80",
				}, {
					Name:                 "example_com_443|openkcm",
					ClusterDiscoveryType: &clusterv3.Cluster_Type{Type: clusterv3.Cluster_STRICT_DNS},
					ConnectTimeout:       &durationpb.Duration{Seconds: 2},
					DnsLookupFamily:      clusterv3.Cluster_V4_ONLY,
					LoadAssignment: &endpointv3.ClusterLoadAssignment{
						ClusterName: "example_com_443|openkcm",
						Endpoints: []*endpointv3.LocalityLbEndpoints{{
							LbEndpoints: []*endpointv3.LbEndpoint{{
								HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
									Endpoint: &endpointv3.Endpoint{
										Address: &corev3.Address{
											Address: &corev3.Address_SocketAddress{
												SocketAddress: &corev3.SocketAddress{
													Address: "example.com",
													PortSpecifier: &corev3.SocketAddress_PortValue{
														PortValue: uint32(443),
													},
												},
											},
										},
									},
								},
							}},
						}},
					},
					TransportSocket: &corev3.TransportSocket{
						Name: wellknown.TransportSocketTls,
						ConfigType: &corev3.TransportSocket_TypedConfig{
							TypedConfig: mustNewAny(&tlsv3.UpstreamTlsContext{
								Sni: "example.com",
								CommonTlsContext: &tlsv3.CommonTlsContext{
									ValidationContextType: &tlsv3.CommonTlsContext_ValidationContext{
										ValidationContext: &tlsv3.CertificateValidationContext{
											TrustedCa: &corev3.DataSource{
												Specifier: &corev3.DataSource_Filename{
													Filename: envoyTrustBundle,
												},
											},
										},
									},
								},
							}),
						},
					},
				}},
				Secrets: []*tlsv3.Secret{
					{
						Name: "secret",
						Type: &tlsv3.Secret_TlsCertificate{
							TlsCertificate: &tlsv3.TlsCertificate{
								CertificateChain: &corev3.DataSource{
									Specifier: &corev3.DataSource_EnvironmentVariable{EnvironmentVariable: "TLS_CERT"},
								},
								PrivateKey: &corev3.DataSource{
									Specifier: &corev3.DataSource_EnvironmentVariable{EnvironmentVariable: "TLS_KEY"},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &GatewayExtension{jwtAuthClusters: maps.Clone(tt.jwtAuthClusters)}
			got, err := s.PostTranslateModify(t.Context(), tt.req)
			if !tt.wantErr(t, err, fmt.Sprintf("PostTranslateModify(%v)", tt.req)) {
				return
			}
			diff := cmp.Diff(tt.want, got, protocmp.Transform(), protocmp.IgnoreDefaultScalars())
			if diff != "" {
				assert.Fail(t, fmt.Sprintf("Not equal: \n"+
					"expected: %s\n"+
					"actual  : %s%s", tt.want, got, diff), "PostTranslateModify(%v)", tt.req)
			}
			if len(s.jwtAuthClusters) >= len(tt.jwtAuthClusters) {
				assert.Fail(t, "Expected read jwtAuthClusters")
			}
		})
	}
}

func TestGatewayExtension_PostVirtualHostModify(t *testing.T) {
	tests := []struct {
		name    string
		req     *extension.PostVirtualHostModifyRequest
		want    *extension.PostVirtualHostModifyResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Post Virtual Host Modify",
			req: &extension.PostVirtualHostModifyRequest{
				VirtualHost: &routev3.VirtualHost{
					Name: "example_com_443|openkcm",
					Routes: []*routev3.Route{{
						Name: "route-with-jwt-auth",
						TypedPerFilterConfig: map[string]*anypb.Any{
							egv1a1.EnvoyFilterJWTAuthn.String(): mustNewAny(&jwtauth3.PerRouteConfig{
								RequirementSpecifier: &jwtauth3.PerRouteConfig_RequirementName{RequirementName: JwtAuthSecureMappingName},
							}),
						},
					}, {
						Name:                 "route-without-jwt-auth",
						TypedPerFilterConfig: make(map[string]*anypb.Any),
					}},
				},
				PostVirtualHostContext: &extension.PostVirtualHostExtensionContext{},
			},
			want: &extension.PostVirtualHostModifyResponse{VirtualHost: &routev3.VirtualHost{
				Name: "example_com_443|openkcm",
				Routes: []*routev3.Route{{
					Name: "route-with-jwt-auth",
					TypedPerFilterConfig: map[string]*anypb.Any{
						egv1a1.EnvoyFilterJWTAuthn.String(): mustNewAny(&jwtauth3.PerRouteConfig{
							RequirementSpecifier: &jwtauth3.PerRouteConfig_RequirementName{RequirementName: JwtAuthSecureMappingName},
						}),
					},
				}, {
					Name: "route-without-jwt-auth",
					TypedPerFilterConfig: map[string]*anypb.Any{
						egv1a1.EnvoyFilterJWTAuthn.String(): mustNewAny(&jwtauth3.PerRouteConfig{
							RequirementSpecifier: &jwtauth3.PerRouteConfig_RequirementName{RequirementName: JwtAuthSecureMappingName},
						}),
					},
				}},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "Nil virtual host",
			req: &extension.PostVirtualHostModifyRequest{
				VirtualHost: nil,
			},
			want: &extension.PostVirtualHostModifyResponse{
				VirtualHost: nil,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGatewayExtension()
			got, err := s.PostVirtualHostModify(t.Context(), tt.req)
			if !tt.wantErr(t, err, fmt.Sprintf("PostVirtualHostModify(ctx, %v)", tt.req)) {
				return
			}
			diff := cmp.Diff(tt.want, got, protocmp.Transform(), protocmp.IgnoreDefaultScalars())
			if diff != "" {
				assert.Fail(t, fmt.Sprintf("Not equal: \n"+
					"expected: %s\n"+
					"actual  : %s%s", tt.want, got, diff), "PostVirtualHostModify(%v)", tt.req)
			}
		})
	}
}
