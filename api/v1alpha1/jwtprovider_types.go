// Copyright Open KCM
// License-Identifier:  Version 2.0, January 2004
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=jwtproviders
//
// JWTProvider provides an example extension policy context resource.
//
//nolint:godoclint
type JWTProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec JWTProviderSpec `json:"spec"`
}

func init() {
	SchemeBuilder.Register(&JWTProvider{}, &JWTProviderList{})
}

// JWTProviderSpec defines how a JSON Web Token (JWT) can be verified.
type JWTProviderSpec struct {
	TargetRefs []gwapiv1.LocalObjectReference `json:"targetRefs"`

	// Name defines a unique name for the JWT provider. A name can have a variety of forms,
	// including RFC1123 subdomains, RFC 1123 labels, or RFC 1035 labels.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Name string `json:"name"`

	// Issuer is the principal that issued the JWT and takes the form of a URL or email address.
	// For additional details, see https://tools.ietf.org/html/rfc7519#section-4.1.1 for
	// URL format and https://rfc-editor.org/rfc/rfc5322.html for email format. If not provided,
	// the JWT issuer is not checked.
	//
	// +kubebuilder:validation:MaxLength=2048
	Issuer string `json:"issuer"`

	// Audiences is a list of JWT audiences allowed access. For additional details, see
	// https://tools.ietf.org/html/rfc7519#section-4.1.3. If not provided, JWT audiences
	// are not checked.
	//
	// +kubebuilder:validation:MaxItems=8
	// +optional
	Audiences []string `json:"audiences,omitempty"`

	// JWKS can be fetched from remote server via HTTP/HTTPS. This field specifies the remote HTTP
	// URI and how the fetched JWKS should be cached.
	//
	// +optional
	RemoteJwks *RemoteJWKS `json:"remoteJwks"`

	// Requires that the credential contains an `expiration <https://tools.ietf.org/html/rfc7519#section-4.1.4>`_.
	// For instance, this could implement JWT-SVID
	// `expiration restrictions <https://github.com/spiffe/spiffe/blob/main/standards/JWT-SVID.md#33-expiration-time>`_.
	// Unlike “max_lifetime“, this only requires that expiration is present, where “max_lifetime“ also checks the value.
	//
	// +optional
	RequireExpiration bool `json:"requireExpiration,omitempty"`

	// RecomputeRoute clears the route cache and recalculates the routing decision.
	// This field must be enabled if the headers generated from the claim are used for
	// route matching decisions. If the recomputation selects a new route, features targeting
	// the new matched route will be applied.
	//
	// +optional
	RecomputeRoute *bool `json:"recomputeRoute,omitempty"`

	// Two fields below define where to extract the JWT from an HTTP request.
	//
	// If no explicit location is specified, the following default locations are tried in order:
	//
	// 1. The Authorization header using the `Bearer schema
	// <https://tools.ietf.org/html/rfc6750#section-2.1>`_. Example::
	//
	//	Authorization: Bearer <token>.
	//
	// 2. `access_token <https://tools.ietf.org/html/rfc6750#section-2.3>`_ query parameter.
	//
	// Multiple JWTs can be verified for a request. Each JWT has to be extracted from the locations
	// its provider specified or from the default locations.
	//
	// +optional
	FromHeaders []*JWTHeader `json:"fromHeaders,omitempty"`

	// Add JWT claim to HTTP JWTHeader
	// Specify the claim name you want to copy in which HTTP header. For examples, following config:
	// The claim must be of type; string, int, double, bool. Array type claims are not supported
	//
	// +optional
	ClaimToHeaders []*JWTClaimToHeader `json:"claimToHeaders,omitempty"`

	// ExtractFrom defines different ways to extract the JWT token from HTTP request.
	// If empty, it defaults to extract JWT token from the Authorization HTTP request header using Bearer schema
	// or access_token from query parameters.
	//
	// +optional
	ExtractFrom *JWTExtractor `json:"extractFrom,omitempty"`
}

// +kubebuilder:object:root=true
//
// JWTProviderList contains a list of ListenerContext resources.
//
//nolint:godoclint
type JWTProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []JWTProvider `json:"items"`
}

type JWTHeader struct {
	// Name is the HTTP header name to retrieve the token
	//
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// ValuePrefix is the prefix that should be stripped before extracting the token.
	// The format would be used by Envoy like "{ValuePrefix}<TOKEN>".
	// For example, "Authorization: Bearer <TOKEN>", then the ValuePrefix="Bearer " with a space at the end.
	//
	// +optional
	ValuePrefix *string `json:"valuePrefix,omitempty"`
}

// JWTClaimToHeader  This message specifies a combination of header name and claim name.
type JWTClaimToHeader struct {
	// The HTTP header name to copy the claim to.
	// The header name will be sanitized and replaced.
	//
	// +kubebuilder:validation:MinLength=1
	HeaderName string `json:"headerName"`
	// The field name for the JWT Claim : it can be a nested claim of type (eg. "claim.nested.key", "sub")
	// String separated with "." in case of nested claims. The nested claim name must use dot "." to separate
	// the JSON name path.
	//
	// +kubebuilder:validation:MinLength=1
	ClaimName string `json:"claimName"`
}

type RemoteJWKS struct {
	// URI is the HTTPS URI to fetch the JWKS. Envoy's system trust bundle is used to validate the server certificate.
	// If a custom trust bundle is needed, it can be specified in a BackendTLSConfig resource and target the BackendRefs.
	//
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	URI string `json:"uri"`

	// Sets the maximum duration in milliseconds that a response can take to arrive upon request.
	//
	// +optional
	TimeoutSec int64 `json:"timeoutSec,omitempty"`

	// Duration after which the cached JWKS should be expired. If not specified, default cache
	// duration is 10 minutes.
	//
	// +optional
	CacheDuration int64 `json:"cacheDuration,omitempty"`

	// Retry define the retry policy configuration.
	//
	// +optional
	Retry *Retry `json:"retry,omitempty"`
}

// JWTExtractor defines a custom JWT token extraction from HTTP request.
// If specified, Envoy will extract the JWT token from the listed extractors (headers, cookies, or params) and validate each of them.
// If any value extracted is found to be an invalid JWT, a 401 error will be returned.
type JWTExtractor struct {
	// Headers represents a list of HTTP request headers to extract the JWT token from.
	//
	// +optional
	Headers []*JWTHeader `json:"headers,omitempty"`

	// Cookies represents a list of cookie names to extract the JWT token from.
	//
	// +optional
	Cookies []string `json:"cookies,omitempty"`

	// Params represents a list of query parameters to extract the JWT token from.
	//
	// +optional
	Params []string `json:"params,omitempty"`
}

// Retry define the retry
type Retry struct {
	// RetryOn configuration. Defaults to connect-failure,refused-stream,unavailable,cancelled,retriable-status-codes.
	//
	// +optional
	RetryOn string `json:"retryOn,omitempty"`

	// NumRetries is the number of retries to be attempted. Defaults to 2.
	//
	// +optional
	NumRetries *uint32 `json:"numRetries,omitempty"`

	// Backoff is the backoff policy to be applied per retry attempt.
	//
	// +optional
	BackOff *BackOffPolicy `json:"backOff,omitempty"`
}

type BackOffPolicy struct {
	// BaseIntervalSec is the base interval between retries.
	BaseIntervalSec int64 `json:"baseIntervalSec,omitempty"`
	// MaxIntervalSec is the maximum interval between retries.
	MaxIntervalSec int64 `json:"maxIntervalSec,omitempty"`
}
