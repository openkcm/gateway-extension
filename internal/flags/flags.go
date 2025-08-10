package flags

const (
	DisableJWTProviderComputation = "disable-jwt-provider-computation"
	// EnableAllowMissingJwtAuthenticationEnvoy If the flag has true value will result to add AllowMissingOrFailed on
	//Envoy JWT Authentication as configuration on missing JWT Providers
	EnableAllowMissingJwtAuthenticationEnvoy = "enable-allow-missing-jwt-authentication-envoy"
)
