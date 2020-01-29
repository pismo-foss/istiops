# Changelog
All notable changes to this project will be documented in this file.

## [2.1.0] - 2020-01-28
### Feature
- add granular route `clear` modes: `soft` (default) & `hard`. More details at [documentation](https://github.com/pismo/istiops/blob/master/README.md). - [#19](https://github.com/pismo/istiops/issues/19)

### Fixes
- istio and kubernetes clients are now supported by `InClusterConfig` for authentication
- CLI logs are now correctly forwarded

## [2.0.0] - 2019-12-10
### Break
- add count of routable pods for each istio's routes at `show` command

## [1.1.0] - 2019-10-11
### Feature
- add support to regexp values for canary release at `shift` command
