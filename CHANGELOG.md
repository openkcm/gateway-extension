# Changelog

## [0.4.0](https://github.com/openkcm/gateway-extension/compare/v0.3.2...v0.4.0) (2025-08-08)


### Features

* upgrade the gateway version 1.5.0 ([#64](https://github.com/openkcm/gateway-extension/issues/64)) ([948dbd5](https://github.com/openkcm/gateway-extension/commit/948dbd5e2c030888d045751f28b99ca33d46baaf))

## [0.3.2](https://github.com/openkcm/gateway-extension/compare/v0.3.1...v0.3.2) (2025-08-07)


### Bug Fixes

* make corrections on the feature flags values, as viper is lowercasing all map keys as value ([#60](https://github.com/openkcm/gateway-extension/issues/60)) ([ba3eff8](https://github.com/openkcm/gateway-extension/commit/ba3eff8e98a9c5b5e8bd0f00f666153b164940dd))

## [0.3.1](https://github.com/openkcm/gateway-extension/compare/v0.3.0...v0.3.1) (2025-08-02)


### Bug Fixes

* return same list of clusters if list of clusters is empty ([#57](https://github.com/openkcm/gateway-extension/issues/57)) ([5e2b253](https://github.com/openkcm/gateway-extension/commit/5e2b253d19ae1a548bff0b3f0f6ed13c4453c566))

## [0.3.0](https://github.com/openkcm/gateway-extension/compare/v0.2.3...v0.3.0) (2025-08-02)


### Features

* **gitub-actions:** refactor release and added new one ([#45](https://github.com/openkcm/gateway-extension/issues/45)) ([5a04271](https://github.com/openkcm/gateway-extension/commit/5a042715025b530f18f08953e4258797b2e38636))
* initial version of the gateway extension ([#1](https://github.com/openkcm/gateway-extension/issues/1)) ([0c87bf9](https://github.com/openkcm/gateway-extension/commit/0c87bf9d1db76e4923ab55d320ec67485f9ac967))
* introduce feature gates; add capability to disable the jwt providers processing ([#46](https://github.com/openkcm/gateway-extension/issues/46)) ([8d7b1fb](https://github.com/openkcm/gateway-extension/commit/8d7b1fbe57ba7b7c570f5f94afe1d94d1d953e92))
* move the github actions into build repo ([07153e4](https://github.com/openkcm/gateway-extension/commit/07153e4ae88d74914fd97e56e163982030f5b966))
* remove the usage of git submodules ([#14](https://github.com/openkcm/gateway-extension/issues/14)) ([e3abd40](https://github.com/openkcm/gateway-extension/commit/e3abd406097981cccd0a7c0072d7accc0371973e))
* update the github action ([#27](https://github.com/openkcm/gateway-extension/issues/27)) ([f941b72](https://github.com/openkcm/gateway-extension/commit/f941b721a244eaab5ef041cc916401791eeedda6))
* update the github action to sign the commits of the pull request ([#22](https://github.com/openkcm/gateway-extension/issues/22)) ([e13ea20](https://github.com/openkcm/gateway-extension/commit/e13ea20e3eb473ee072b73a71c7155080b2858c5))
* update the workflow to support the generation of sbom ([#8](https://github.com/openkcm/gateway-extension/issues/8)) ([ad41c5f](https://github.com/openkcm/gateway-extension/commit/ad41c5fef66e31fe221e55278df91018ba27c286))


### Bug Fixes

* change the items of JWTProviderList from JWTProviderSpec to JWTProvider ([#3](https://github.com/openkcm/gateway-extension/issues/3)) ([ffb4a6c](https://github.com/openkcm/gateway-extension/commit/ffb4a6c26d377409227c88eeb48e4c2a927a8dc0))
* **chart:** deployment extraPorts variable ([#32](https://github.com/openkcm/gateway-extension/issues/32)) ([291536e](https://github.com/openkcm/gateway-extension/commit/291536e5e43c486188315e42da543aebfbe5378d))
* comment out the update of the chart version in workflows ([#5](https://github.com/openkcm/gateway-extension/issues/5)) ([a4e569e](https://github.com/openkcm/gateway-extension/commit/a4e569eca3c23b9bf9a5b272002aedd754320d81))
* envoy cluster list is not updated ([#53](https://github.com/openkcm/gateway-extension/issues/53)) ([2ede102](https://github.com/openkcm/gateway-extension/commit/2ede102c63fd0e5741fc949f25fd3bd86ff5f60a))
* fetching the github token from the right location ([#17](https://github.com/openkcm/gateway-extension/issues/17)) ([4040c0d](https://github.com/openkcm/gateway-extension/commit/4040c0dae1b6ddac5445d334b8f47c955c74510b))
* fix a name into the release metadata file ([a8c1486](https://github.com/openkcm/gateway-extension/commit/a8c148620e4c46365867dd24a0ddf3ee808d3deb))
* fixing the github workflow ([#4](https://github.com/openkcm/gateway-extension/issues/4)) ([ff1e5a8](https://github.com/openkcm/gateway-extension/commit/ff1e5a82454bc69dca672298c03efa0d6884fea8))
* gRPC client configuration handling ([#49](https://github.com/openkcm/gateway-extension/issues/49)) ([eff6017](https://github.com/openkcm/gateway-extension/commit/eff601759bbf41bee92095e1d18a3ca2b84ac017))
* include the mutex to block while creating the envoy clusters resources ([#47](https://github.com/openkcm/gateway-extension/issues/47)) ([2779ca2](https://github.com/openkcm/gateway-extension/commit/2779ca28d92363293080e26cb4c714952334ef9e))
* small code rearrangement ([#51](https://github.com/openkcm/gateway-extension/issues/51)) ([8815bc8](https://github.com/openkcm/gateway-extension/commit/8815bc8fa2cc6732b606557b4639e75b22bf9815))
* updated the code to comply with lints ([#55](https://github.com/openkcm/gateway-extension/issues/55)) ([f25fed4](https://github.com/openkcm/gateway-extension/commit/f25fed4e724e16051cf6345675572a9cfc04836e))

## [0.2.3](https://github.com/openkcm/gateway-extension/compare/v0.2.2...v0.2.3) (2025-08-02)


### Bug Fixes

* envoy cluster list is not updated ([#53](https://github.com/openkcm/gateway-extension/issues/53)) ([2ede102](https://github.com/openkcm/gateway-extension/commit/2ede102c63fd0e5741fc949f25fd3bd86ff5f60a))

## [0.2.2](https://github.com/openkcm/gateway-extension/compare/v0.2.1...v0.2.2) (2025-08-02)


### Bug Fixes

* small code rearrangement ([#51](https://github.com/openkcm/gateway-extension/issues/51)) ([8815bc8](https://github.com/openkcm/gateway-extension/commit/8815bc8fa2cc6732b606557b4639e75b22bf9815))

## [0.2.1](https://github.com/openkcm/gateway-extension/compare/v0.2.0...v0.2.1) (2025-08-02)


### Bug Fixes

* gRPC client configuration handling ([#49](https://github.com/openkcm/gateway-extension/issues/49)) ([eff6017](https://github.com/openkcm/gateway-extension/commit/eff601759bbf41bee92095e1d18a3ca2b84ac017))

## [0.2.0](https://github.com/openkcm/gateway-extension/compare/v0.1.3...v0.2.0) (2025-08-02)


### Features

* **gitub-actions:** refactor release and added new one ([#45](https://github.com/openkcm/gateway-extension/issues/45)) ([5a04271](https://github.com/openkcm/gateway-extension/commit/5a042715025b530f18f08953e4258797b2e38636))
* introduce feature gates; add capability to disable the jwt providers processing ([#46](https://github.com/openkcm/gateway-extension/issues/46)) ([8d7b1fb](https://github.com/openkcm/gateway-extension/commit/8d7b1fbe57ba7b7c570f5f94afe1d94d1d953e92))
* move the github actions into build repo ([07153e4](https://github.com/openkcm/gateway-extension/commit/07153e4ae88d74914fd97e56e163982030f5b966))
* update the github action ([#27](https://github.com/openkcm/gateway-extension/issues/27)) ([f941b72](https://github.com/openkcm/gateway-extension/commit/f941b721a244eaab5ef041cc916401791eeedda6))
* update the github action to sign the commits of the pull request ([#22](https://github.com/openkcm/gateway-extension/issues/22)) ([e13ea20](https://github.com/openkcm/gateway-extension/commit/e13ea20e3eb473ee072b73a71c7155080b2858c5))


### Bug Fixes

* **chart:** deployment extraPorts variable ([#32](https://github.com/openkcm/gateway-extension/issues/32)) ([291536e](https://github.com/openkcm/gateway-extension/commit/291536e5e43c486188315e42da543aebfbe5378d))
* fetching the github token from the right location ([#17](https://github.com/openkcm/gateway-extension/issues/17)) ([4040c0d](https://github.com/openkcm/gateway-extension/commit/4040c0dae1b6ddac5445d334b8f47c955c74510b))
* fix a name into the release metadata file ([a8c1486](https://github.com/openkcm/gateway-extension/commit/a8c148620e4c46365867dd24a0ddf3ee808d3deb))
* include the mutex to block while creating the envoy clusters resources ([#47](https://github.com/openkcm/gateway-extension/issues/47)) ([2779ca2](https://github.com/openkcm/gateway-extension/commit/2779ca28d92363293080e26cb4c714952334ef9e))
