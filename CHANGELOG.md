# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 1.10.1  (2023-04-03)
### Changed
- Fix log example file path in unix tarballs
- Disable CGO

## 1.10.0  (2023-03-08)
### Changed
- Upgrade Go to 1.19 and bump dependencies

## 1.9.1 (2022-06-27)
### Added
 - Added support for more distributions:
    RHEL(EL) 9
    Ubuntu 22.04

## 1.9.0 (2022-04-27)
### Added
- Log config examples files.

### Changed
- Use Go 1.18.
- Bump integration and tools dependencies.
- Remove unused Dockerfile.
## 1.8.0 (2022-02-08)
### Added
- Added `binary_path` config setting. Allows user to specify the a custom location of the Apache binary file for inventory collection purposes.

## 1.7.1 (2021-10-20)
### Added
Added support for more distributions:
- Debian 11
- Ubuntu 20.10
- Ubuntu 21.04
- SUSE 12.15
- SUSE 15.1
- SUSE 15.2
- SUSE 15.3
- Oracle Linux 7
- Oracle Linux 8

## 1.7.0 (2021-06-08)
### Changed

Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 1.6.1 (2021-06-08)
### Changed
- Support for ARM

## 1.6.0 (2021-05-06)
## Changed
- Update Go to v1.16.
- Migrate to Go Modules
- Update Infrastracture SDK to v3.6.7.
- Update other dependecies.

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
