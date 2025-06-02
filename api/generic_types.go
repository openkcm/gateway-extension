package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gev1a1 "github.com/openkcm/gateway-extension/api/v1alpha1"
)

type Generic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

const (
	JWTProviderKind = "JWTProvider"
)

var (
	JWTProviderV1Alpha1 = gev1a1.GroupVersion.String()
)
