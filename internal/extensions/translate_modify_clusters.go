package extensions

import (
	"context"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	slogctx "github.com/veqryn/slog-context"
)

func (s *GatewayExtension) TranslateModifyClusters(ctx context.Context, cls []*clusterv3.Cluster) ([]*clusterv3.Cluster, error) {
	clusters := make([]*clusterv3.Cluster, 0)

	// remove clusters that has as suffix name the `openkcm`,
	for _, c := range cls {
		if IsCustomName(c.GetName()) {
			continue
		}
		clusters = append(clusters, c)
	}

	// will be added new list of the clusters
	for k, v := range s.jwtAuthClusters {
		clusterName := v.CustomName()

		slogctx.Info(ctx, "Processing cached cluster", "cluster", clusterName)

		trCtx, err := anypb.New(buildXdsUpstreamTLSSocket(v.hostname))
		if err != nil {
			return nil, err
		}

		cluster := &clusterv3.Cluster{
			Name:                 clusterName,
			ClusterDiscoveryType: &clusterv3.Cluster_Type{Type: clusterv3.Cluster_STRICT_DNS},
			ConnectTimeout:       &durationpb.Duration{Seconds: 2},
			DnsLookupFamily:      clusterv3.Cluster_V4_ONLY,
			LoadAssignment: &endpointv3.ClusterLoadAssignment{
				ClusterName: clusterName,
				Endpoints: []*endpointv3.LocalityLbEndpoints{{
					LbEndpoints: []*endpointv3.LbEndpoint{{
						HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
							Endpoint: &endpointv3.Endpoint{
								Address: &corev3.Address{
									Address: &corev3.Address_SocketAddress{
										SocketAddress: &corev3.SocketAddress{
											Address: v.hostname,
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
					TypedConfig: trCtx,
				},
			},
		}
		clusters = append(clusters, cluster)
		delete(s.jwtAuthClusters, k)
	}

	return clusters, nil
}

const (
	envoyTrustBundle = "/etc/ssl/certs/ca-certificates.crt"
)

func buildXdsUpstreamTLSSocket(sni string) *tlsv3.UpstreamTlsContext {
	return &tlsv3.UpstreamTlsContext{
		Sni: sni,
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
	}
}
