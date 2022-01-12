# Changelog

## [0.4.0] - 2022-01-12

### Added
* Enhance data collection and debugging information ([#31](https://github.com/TibiaData/tibiadata-api-go/pull/31) by [tobiasehlert](https://github.com/tobiasehlert))
* Adding rate-limit detection and status code switching ([#36](https://github.com/TibiaData/tibiadata-api-go/pull/36) by [tobiasehlert](https://github.com/tobiasehlert))
* Adding TIBIADATA_HOST env ([#39](https://github.com/TibiaData/tibiadata-api-go/pull/39) by [tobiasehlert](https://github.com/tobiasehlert))
* Rewrite of characters and adding of unit testing ([#40](https://github.com/TibiaData/tibiadata-api-go/pull/40) by [kamilon](https://github.com/kamilon))
* Adding codecov.io ([#48](https://github.com/TibiaData/tibiadata-api-go/pull/48) by [tobiasehlert](https://github.com/tobiasehlert))
* Adding new Tibia News endpoints ([#32](https://github.com/TibiaData/tibiadata-api-go/pull/32) by [tobiasehlert](https://github.com/tobiasehlert))
* Worlds Overview unit tests ([#53](https://github.com/TibiaData/tibiadata-api-go/pull/53) by [kamilon](https://github.com/kamilon))
* World unit tests ([#54](https://github.com/TibiaData/tibiadata-api-go/pull/54) by [kamilon](https://github.com/kamilon))

### Changed
* Refactor to use pure goquery and no manual regex on guilds ([#29](https://github.com/TibiaData/tibiadata-api-go/pull/29) by [JorgeMag96](https://github.com/JorgeMag96))
* Performing go mod tidy ([#43](https://github.com/TibiaData/tibiadata-api-go/pull/43) by [tobiasehlert](https://github.com/tobiasehlert))
* Use new assert instead of passing t many times ([#45](https://github.com/TibiaData/tibiadata-api-go/pull/45) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
* Updating workflows ([#51](https://github.com/TibiaData/tibiadata-api-go/pull/51) by [tobiasehlert](https://github.com/tobiasehlert))
* Bump golang from 1.17.5 to 1.17.6 ([#50](https://github.com/TibiaData/tibiadata-api-go/pull/50) by [dependabot](https://github.com/dependabot))
* Reduce dependencies ([#44](https://github.com/TibiaData/tibiadata-api-go/pull/44) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
* Clean up webserver.go and a few utility funcs ([#37](https://github.com/TibiaData/tibiadata-api-go/pull/37), [#52](https://github.com/TibiaData/tibiadata-api-go/pull/52) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

### Fixed
* Maintenance mode detection with error message return ([#34](https://github.com/TibiaData/tibiadata-api-go/pull/34) by [tobiasehlert](https://github.com/tobiasehlert))
* Guilds name upper casing correction ([#35](https://github.com/TibiaData/tibiadata-api-go/pull/35) by [tobiasehlert](https://github.com/tobiasehlert))
* Handle timezone information during DateTime parsing ([#46](https://github.com/TibiaData/tibiadata-api-go/pull/46), [#49](https://github.com/TibiaData/tibiadata-api-go/pull/49) by [kamilon](https://github.com/kamilon), [tobiasehlert](https://github.com/tobiasehlert))

## [0.3.0] - 2022-01-05

### Added
* Addition of two guild endpoints ([#20](https://github.com/TibiaData/tibiadata-api-go/pull/20) by [tobiasehlert](https://github.com/tobiasehlert))
* Adding proxy support to replace default URL ([#25](https://github.com/TibiaData/tibiadata-api-go/pull/25) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed
* Clear lint errors ([#18](https://github.com/TibiaData/tibiadata-api-go/pull/18) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
* Idiomatic go ([#19](https://github.com/TibiaData/tibiadata-api-go/pull/19) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
* Removing duplicate function ([#27](https://github.com/TibiaData/tibiadata-api-go/pull/27) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed
* Highscores endpoints redirect ([#22](https://github.com/TibiaData/tibiadata-api-go/pull/22) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

## [0.2.0] - 2022-01-02

### Changed
* Sanitize certain data ([#14](https://github.com/TibiaData/tibiadata-api-go/pull/14) by [darrentaytay](https://github.com/darrentaytay))
* Removing some code smell from SonarCloud ([#10](https://github.com/TibiaData/tibiadata-api-go/pull/10) by [tobiasehlert](https://github.com/tobiasehlert))
* Moving vocation logic to separate function ([#17](https://github.com/TibiaData/tibiadata-api-go/pull/17) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed
* Fix decoding issues and save some allocations ([#15](https://github.com/TibiaData/tibiadata-api-go/pull/15), [#16](https://github.com/TibiaData/tibiadata-api-go/pull/16) by [kamilon](https://github.com/kamilon), [tobiasehlert](https://github.com/tobiasehlert))

## [0.1.1] - 2021-12-31

### Changed
- Bump actions/cache from 2.1.6 to 2.1.7 ([#2](https://github.com/TibiaData/tibiadata-api-go/pull/2) by [dependabot](https://github.com/dependabot))
- Bump docker/metadata-action from 3.6.0 to 3.6.2 ([#1](https://github.com/TibiaData/tibiadata-api-go/pull/1) by [dependabot](https://github.com/dependabot))
- Update README.md ([#3](https://github.com/TibiaData/tibiadata-api-go/pull/3), [#5](https://github.com/TibiaData/tibiadata-api-go/pull/5) by [tobiasehlert](https://github.com/tobiasehlert))
- Implementation of response handler ([#4](https://github.com/TibiaData/tibiadata-api-go/pull/4) by [tobiasehlert](https://github.com/tobiasehlert))
- Changing building and releasing ([#6](https://github.com/TibiaData/tibiadata-api-go/pull/6), [#7](https://github.com/TibiaData/tibiadata-api-go/pull/7) by [tobiasehlert](https://github.com/tobiasehlert))
- Shrink of TibiaFansitesV3 ([#8](https://github.com/TibiaData/tibiadata-api-go/pull/8) by [tobiasehlert](https://github.com/tobiasehlert))
- Updating build workflow with enhancements ([#11](https://github.com/TibiaData/tibiadata-api-go/pull/11) by [tobiasehlert](https://github.com/tobiasehlert))


## [0.1.0] - 2021-12-23

Initial commit

[0.4.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/tibiadata/tibiadata-api-go/compare/30f328f...v0.1.0
