package extensions

import (
	"fmt"

	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
)

// Tries to find an HTTP connection manager in the provided filter chain.
func findHCM(filterChain *listenerv3.FilterChain) (*hcm.HttpConnectionManager, int, error) {
	for filterIndex, filter := range filterChain.GetFilters() {
		if filter.GetName() == wellknown.HTTPConnectionManager {
			h := new(hcm.HttpConnectionManager)
			if err := filter.GetTypedConfig().UnmarshalTo(h); err != nil {
				return nil, -1, err
			}
			return h, filterIndex, nil
		}
	}
	return nil, -1, fmt.Errorf("unable to find HTTPConnectionManager in FilterChain: %s", filterChain.GetName())
}
