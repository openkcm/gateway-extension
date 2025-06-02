package extensions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCustomName(t *testing.T) {
	tests := []struct {
		name        string
		clusterName string
		want        bool
	}{
		{
			name:        "Custom name",
			clusterName: "cluster|openkcm",
			want:        true,
		},
		{
			name:        "Non-custom name",
			clusterName: "cluster",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, IsCustomName(tt.clusterName), "IsCustomName(%v)", tt.clusterName)
		})
	}
}

func Test_url2Cluster(t *testing.T) {
	tests := []struct {
		name    string
		strURL  string
		want    *urlCluster
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "Hostname",
			strURL: "https://example.com/jwks",
			want: &urlCluster{
				name:         "example_com_443",
				hostname:     "example.com",
				port:         443,
				endpointType: EndpointTypeDNS,
				tls:          true,
			},
			wantErr: assert.NoError,
		},
		{
			name:   "IPv4",
			strURL: "https://127.0.0.1:443/jwks",
			want: &urlCluster{
				name:         "127_0_0_1_443",
				hostname:     "127.0.0.1",
				port:         443,
				endpointType: EndpointTypeStatic,
				tls:          true,
			},
			wantErr: assert.NoError,
		},
		{
			name:   "No TLS",
			strURL: "http://example.com/jwks",
			want: &urlCluster{
				name:         "example_com_80",
				hostname:     "example.com",
				port:         80,
				endpointType: EndpointTypeDNS,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := url2Cluster(tt.strURL)
			if !tt.wantErr(t, err, fmt.Sprintf("url2Cluster(%v)", tt.strURL)) {
				return
			}
			assert.Equalf(t, tt.want, got, "url2Cluster(%v)", tt.strURL)
		})
	}
}
