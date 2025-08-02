package extensions

import (
	"fmt"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
)

const (
	defaultHTTPSPort uint64 = 443
	defaultHTTPPort  uint64 = 80
)

type EndpointType int

const (
	customSuffixName = "openkcm"

	EndpointTypeDNS EndpointType = iota
	EndpointTypeStatic
)

// urlCluster is a cluster that is created from a URL.
type urlCluster struct {
	name         string
	hostname     string
	port         uint32
	endpointType EndpointType
	tls          bool
}

// url2Cluster returns a urlCluster from the provided url.
func url2Cluster(strURL string) (*urlCluster, error) {
	epType := EndpointTypeDNS

	// The URL should have already been validated in the gateway API translator.
	u, err := url.Parse(strURL)
	if err != nil {
		return nil, err
	}

	var port uint64
	if u.Scheme == "https" {
		port = defaultHTTPSPort
	} else {
		port = defaultHTTPPort
	}

	if u.Port() != "" {
		port, err = strconv.ParseUint(u.Port(), 10, 32)
		if err != nil {
			return nil, err
		}
	}

	name := clusterName(u.Hostname(), uint32(port))

	ip, err := netip.ParseAddr(u.Hostname())
	if err == nil {
		if ip.Unmap().Is4() {
			epType = EndpointTypeStatic
		}
	}

	return &urlCluster{
		name:         name,
		hostname:     u.Hostname(),
		port:         uint32(port),
		endpointType: epType,
		tls:          u.Scheme == "https",
	}, nil
}

func clusterName(host string, port uint32) string {
	return fmt.Sprintf("%s_%d", strings.ReplaceAll(host, ".", "_"), port)
}

func (c *urlCluster) CustomName() string {
	return fmt.Sprintf("%s|%s", c.name, customSuffixName)
}
func IsCustomName(name string) bool {
	return strings.HasSuffix(name, customSuffixName)
}
