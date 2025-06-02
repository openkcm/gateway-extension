package testdata

import _ "embed"

//go:embed extension.json
var ExtensionJSON []byte

//go:embed well-known.json
var WellKnownJSON []byte

//go:embed openid-configuration.json
var OpenIDConfigurationJSON []byte
