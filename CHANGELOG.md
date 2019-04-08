# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

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
