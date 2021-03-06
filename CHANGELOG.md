# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 1.5.1 (2020-06-16)
## Fixed
- Updated the configuration sample to include the status_url for inventory required
  for entity naming.

## 1.5.0 (2019-12-10)
## Added
- Added `validate_certs` configuration option (default: `true`). Set it to `false` if you have self-signed certificates
  and want to avoid the integration to fail.

## 1.4.0 (2019-11-15)
### Changed
- Renamed the integration executable from nr-apache to nri-apache in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.

## 1.3.0 (2019-04-29)
### Added
- Upgraded to SDK v3.1.5. This version implements [the aget/integrations
  protocol v3](https://github.com/newrelic/infra-integrations-sdk/blob/cb45adacda1cd5ff01544a9d2dad3b0fedf13bf1/docs/protocol-v3.md),
  which enables [name local address replacement](https://github.com/newrelic/infra-integrations-sdk/blob/cb45adacda1cd5ff01544a9d2dad3b0fedf13bf1/docs/protocol-v3.md#name-local-address-replacement).
  and could change your entity names and alarms. For more information, refer
  to:

  - https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-executable-file-specifications#h2-loopback-address-replacement-on-entity-names
  - https://docs.newrelic.com/docs/remote-monitoring-host-integration://docs.newrelic.com/docs/remote-monitoring-host-integrations

## 1.2.0 (2019-04-08)
### Added
- Upgraded to SDKv3
- Remote monitoring option. It enables monitoring multiple Apache instances,
  more information can be found at the [official documentation page](https://docs.newrelic.com/docs/remote-monitoring-host-integrations).

## 1.1.2 (2018-10-18)
### Fixed
- The release process was incorrectly triggered, fixing tags and versioning. None change in the integration.

## 1.1.1 (2018-10-17)
### Fixed
- Error on weird modules output with message: `slice bounds out of range`.

## 1.1.0 (2018-02-08)
### Added
- Allow working with own's Apache Certificate Authority through the `ca_bundle_file` and `ca_bundle_dir` configuration
  options.

## 1.0.0 (2017-11-29)
### Added
- Initial release, which contains inventory and metrics data
