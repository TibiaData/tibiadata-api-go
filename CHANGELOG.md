# Changelog

## [3.5.0] - 2022-10-19

### Added
* fix highscores by implementing pagination ([#164](https://github.com/TibiaData/tibiadata-api-go/pull/164), [#165](https://github.com/TibiaData/tibiadata-api-go/pull/165) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed
* Bump golang from 1.19.0 to 1.19.2 ([#158](https://github.com/TibiaData/tibiadata-api-go/pull/158), [#159](https://github.com/TibiaData/tibiadata-api-go/pull/159) by [dependabot](https://github.com/dependabot))
* Bump golang.org/x/text from 0.3.7 to 0.4.0 ([#160](https://github.com/TibiaData/tibiadata-api-go/pull/160), [#161](https://github.com/TibiaData/tibiadata-api-go/pull/161) by [dependabot](https://github.com/dependabot))

## [3.4.1] - 2022-08-30

### Changed
* Bump golang from 1.18.4 to 1.19.0 ([#156](https://github.com/TibiaData/tibiadata-api-go/pull/156) by [dependabot](https://github.com/dependabot))

### Fixed
* Fix handling of special characters in boosted boss/creature name ([#157](https://github.com/TibiaData/tibiadata-api-go/pull/157) by [tiagomartines11](https://github.com/tiagomartines11))

## [3.4.0] - 2022-08-02

### Added
* Adding proxy protocol support to replace default https ([#154](https://github.com/TibiaData/tibiadata-api-go/pull/154) by [tobiasehlert](https://github.com/tobiasehlert))

### Fixed
* Parsing of spells list ([#155](https://github.com/TibiaData/tibiadata-api-go/pull/155) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.3.1] - 2022-07-29

### Changed
* Added bosspoints as additional HighscoreCategory ([#153](https://github.com/TibiaData/tibiadata-api-go/pull/153) by [Tholdrim](https://github.com/Tholdrim))

## [3.3.0] - 2022-07-27

### Added
* Add support to new boostable bosses page ([#152](https://github.com/TibiaData/tibiadata-api-go/pull/152) by [tiagomartines11](https://github.com/tiagomartines11))

### Changed
* Removing **v** in container image tag ([#144](https://github.com/TibiaData/tibiadata-api-go/pull/144) by [tobiasehlert](https://github.com/tobiasehlert))
* Updating `go build` step in dockerfile ([#148](https://github.com/TibiaData/tibiadata-api-go/pull/148) by [tobiasehlert](https://github.com/tobiasehlert))
* Some go mod and workflow build updates ([#137](https://github.com/TibiaData/tibiadata-api-go/pull/137) by [tobiasehlert](https://github.com/tobiasehlert))
* Bump github.com/gin-gonic/gin from 1.7.7 to 1.8.1 ([#136](https://github.com/TibiaData/tibiadata-api-go/pull/136), [#140](https://github.com/TibiaData/tibiadata-api-go/pull/140) by [dependabot](https://github.com/dependabot))
* Bump golang from 1.18.1 to 1.18.4 ([#135](https://github.com/TibiaData/tibiadata-api-go/pull/135), [#138](https://github.com/TibiaData/tibiadata-api-go/pull/138), [#149](https://github.com/TibiaData/tibiadata-api-go/pull/149) by [dependabot](https://github.com/dependabot))
* Bump github.com/stretchr/testify from 1.7.1 to 1.8.0 ([#139](https://github.com/TibiaData/tibiadata-api-go/pull/139), [#141](https://github.com/TibiaData/tibiadata-api-go/pull/141), [#142](https://github.com/TibiaData/tibiadata-api-go/pull/142), [#143](https://github.com/TibiaData/tibiadata-api-go/pull/143) by [dependabot](https://github.com/dependabot))
* Bump github.com/gin-contrib/gzip from 0.0.5 to 0.0.6 ([#145](https://github.com/TibiaData/tibiadata-api-go/pull/145) by [dependabot](https://github.com/dependabot))

### Fixed
* Parsing issue when guild description contains founded string ([#151](https://github.com/TibiaData/tibiadata-api-go/pull/151) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.2.2] - 2022-05-12

### Changed
* Bump various workflow versions ([#126](https://github.com/TibiaData/tibiadata-api-go/pull/126), [#127](https://github.com/TibiaData/tibiadata-api-go/pull/127), [#128](https://github.com/TibiaData/tibiadata-api-go/pull/128), [#129](https://github.com/TibiaData/tibiadata-api-go/pull/129), [#130](https://github.com/TibiaData/tibiadata-api-go/pull/130), [#131](https://github.com/TibiaData/tibiadata-api-go/pull/131), [#132](https://github.com/TibiaData/tibiadata-api-go/pull/132) by [dependabot](https://github.com/dependabot))

### Fixed
* Fix traded-string appearing in account characters info ([#133](https://github.com/TibiaData/tibiadata-api-go/pull/133) by [tiagomartines11](https://github.com/tiagomartines11))

## [3.2.1] - 2022-04-26

### Changed
* Bump github/codeql-action from 1 to 2 ([#122](https://github.com/TibiaData/tibiadata-api-go/pull/122) by [dependabot](https://github.com/dependabot))

### Fixed
* Trimming suffix on guild rank ([#121](https://github.com/TibiaData/tibiadata-api-go/pull/121) by [tobiasehlert](https://github.com/tobiasehlert))
* Change go version in go.mod ([#124](https://github.com/TibiaData/tibiadata-api-go/pull/124) by [sergot](https://github.com/sergot))
* Adjusting parsing of fansite page ([#125](https://github.com/TibiaData/tibiadata-api-go/pull/125) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.2.0] - 2022-04-24

### Changed
* Adding support for finished auctions ([#115](https://github.com/TibiaData/tibiadata-api-go/pull/115), [#118](https://github.com/TibiaData/tibiadata-api-go/pull/118) by [tobiasehlert](https://github.com/tobiasehlert))
* Enhancing deaths parsing of characters ([#116](https://github.com/TibiaData/tibiadata-api-go/pull/116) by [tobiasehlert](https://github.com/tobiasehlert))
* Bump golang from 1.17.8 to 1.18.1 ([#108](https://github.com/TibiaData/tibiadata-api-go/pull/108), [#113](https://github.com/TibiaData/tibiadata-api-go/pull/113) by [dependabot](https://github.com/dependabot))
* Bump various workflow versions ([#101](https://github.com/TibiaData/tibiadata-api-go/pull/101), [#102](https://github.com/TibiaData/tibiadata-api-go/pull/102), [#109](https://github.com/TibiaData/tibiadata-api-go/pull/109), [#110](https://github.com/TibiaData/tibiadata-api-go/pull/110), [#111](https://github.com/TibiaData/tibiadata-api-go/pull/111), [#112](https://github.com/TibiaData/tibiadata-api-go/pull/112) by [dependabot](https://github.com/dependabot))

### Fixed
* Fix regex to handle one bed in house ([#105](https://github.com/TibiaData/tibiadata-api-go/pull/105) by [tobiasehlert](https://github.com/tobiasehlert))
* Fix encoding of apostrophes in multiple places ([#106](https://github.com/TibiaData/tibiadata-api-go/pull/106) by [tobiasehlert](https://github.com/tobiasehlert))
* Fix guild description to contain guild details ([#107](https://github.com/TibiaData/tibiadata-api-go/pull/107) by [tobiasehlert](https://github.com/tobiasehlert))
* Fix characters missing marriage and some more tests ([#117](https://github.com/TibiaData/tibiadata-api-go/pull/117) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.1.1] - 2022-03-15

### Fixed
* Adding sanitize of nbsp string in death reason of players ([#100](https://github.com/TibiaData/tibiadata-api-go/pull/100) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.1.0] - 2022-03-10

### Changed
* Renaming Tibiadata to TibiaData ([#90](https://github.com/TibiaData/tibiadata-api-go/pull/90) by [tobiasehlert](https://github.com/tobiasehlert))
* Removing injection of houseType from assets ([#95](https://github.com/TibiaData/tibiadata-api-go/pull/95) by [tobiasehlert](https://github.com/tobiasehlert))
* Removing return of loyalty title on no title ([#98](https://github.com/TibiaData/tibiadata-api-go/pull/98) by [tobiasehlert](https://github.com/tobiasehlert))
* Bump docker/login-action from 1.14.0 to 1.14.1 ([#87](https://github.com/TibiaData/tibiadata-api-go/pull/87) by [dependabot](https://github.com/dependabot))
* Bump actions/checkout from 2 to 3 ([#88](https://github.com/TibiaData/tibiadata-api-go/pull/88) by [dependabot](https://github.com/dependabot))
* Bump golang from 1.17.7 to 1.17.8 ([#96](https://github.com/TibiaData/tibiadata-api-go/pull/96) by [dependabot](https://github.com/dependabot))

### Fixed
* Environment function logic fix ([#91](https://github.com/TibiaData/tibiadata-api-go/pull/91) by [tobiasehlert](https://github.com/tobiasehlert))
* Stop using ioutil as it is deprecated ([#92](https://github.com/TibiaData/tibiadata-api-go/pull/92) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
* Adding sanitize of nbsp string in death section of players ([#99](https://github.com/TibiaData/tibiadata-api-go/pull/99) by [tobiasehlert](https://github.com/tobiasehlert))

## [3.0.0] - 2022-03-01

### The release of API v3

Bumping version to *v3.0.0* so that GitHub releases match API version.

Head over to [tibiadata.com](https://tibiadata.com/2022/03/tibiadata-api-v3-released-v3-0-0/) for more information.

## [0.6.2] - 2022-02-28

### Fixed
* Adding if to prevent panic due to regex ([#84](https://github.com/TibiaData/tibiadata-api-go/pull/84) by [tobiasehlert](https://github.com/tobiasehlert))

## [0.6.1] - 2022-02-27

### Changed
* Switching assets download URL to assets.tibiadata.com
* Bump docker/build-push-action from 2.8.0 to 2.9.0 ([#80](https://github.com/TibiaData/tibiadata-api-go/pull/80) by [dependabot](https://github.com/dependabot))
* Bump golang from 1.17.6 to 1.17.7 ([#81](https://github.com/TibiaData/tibiadata-api-go/pull/81) by [dependabot](https://github.com/dependabot))
* Bump docker/login-action from 1.12.0 to 1.13.0 ([#82](https://github.com/TibiaData/tibiadata-api-go/pull/82) by [dependabot](https://github.com/dependabot))

## [0.6.0] - 2022-01-31

### Added
* Addition of some more tests ([#76](https://github.com/TibiaData/tibiadata-api-go/pull/76) by [tobiasehlert](https://github.com/tobiasehlert))
* Implement graceful shutdown ([#78](https://github.com/TibiaData/tibiadata-api-go/pull/78) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))
* Implement 404 page not found ([#77](https://github.com/TibiaData/tibiadata-api-go/pull/77) by [Pedro-Pessoa](https://github.com/Pedro-Pessoa))

### Changed
* Switching to http status codes ([#79](https://github.com/TibiaData/tibiadata-api-go/pull/79) by [tobiasehlert](https://github.com/tobiasehlert))
* Cleanup of characters deathlist ([#75](https://github.com/TibiaData/tibiadata-api-go/pull/75) by [tobiasehlert](https://github.com/tobiasehlert))
* Adjustment of endpoint namings ([#74](https://github.com/TibiaData/tibiadata-api-go/pull/74) by [tobiasehlert](https://github.com/tobiasehlert))

## [0.5.1] - 2022-01-28

### Added
* Adding gzip compression middleware ([#73](https://github.com/TibiaData/tibiadata-api-go/pull/73) by [tobiasehlert](https://github.com/tobiasehlert))
* Swagger documentation annotations ([#67](https://github.com/TibiaData/tibiadata-api-go/pull/67) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed
* Add some 'fake' unit tests to up coverage on webserver.go ([#71](https://github.com/TibiaData/tibiadata-api-go/pull/71) by [kamilon](https://github.com/kamilon))
* Add unit tests for House Overview and House APIs ([#72](https://github.com/TibiaData/tibiadata-api-go/pull/72) by [kamilon](https://github.com/kamilon))

## [0.5.0] - 2022-01-19

### Added
* Add unit tests for Spells Overview API, fix bugs in Spells Overview API ([#55](https://github.com/TibiaData/tibiadata-api-go/pull/55) by [kamilon](https://github.com/kamilon))
* Add unit tests for Spells API, fix various bugs in pulling of spell data ([#56](https://github.com/TibiaData/tibiadata-api-go/pull/56) by [kamilon](https://github.com/kamilon))
* Add unit tests for Kill Statistics API, move to goquery ([#57](https://github.com/TibiaData/tibiadata-api-go/pull/57) by [kamilon](https://github.com/kamilon))
* Add unit tests for Guild Overview and Guild API ([#58](https://github.com/TibiaData/tibiadata-api-go/pull/58) by [kamilon](https://github.com/kamilon))
* Add unit tests for Creatures Overview and Creature API, fix bug ([#59](https://github.com/TibiaData/tibiadata-api-go/pull/59) by [kamilon](https://github.com/kamilon))
* Add unit tests for Highscores API ([#61](https://github.com/TibiaData/tibiadata-api-go/pull/61) by [kamilon](https://github.com/kamilon))
* Add unit tests for Fansites API ([#62](https://github.com/TibiaData/tibiadata-api-go/pull/62) by [kamilon](https://github.com/kamilon))
* Add unit tests for News List API and News API ([#60](https://github.com/TibiaData/tibiadata-api-go/pull/60) by [kamilon](https://github.com/kamilon))
* Create and use HighscoreCategory enum ([#63](https://github.com/TibiaData/tibiadata-api-go/pull/63) by [kamilon](https://github.com/kamilon))
* Feature: Tibia Houses endpoints ([#26](https://github.com/TibiaData/tibiadata-api-go/pull/26) by [tobiasehlert](https://github.com/tobiasehlert))

### Changed
* Refactor gin server to reduce code duplication ([#64](https://github.com/TibiaData/tibiadata-api-go/pull/64) by [kamilon](https://github.com/kamilon))
* Cache new regex queries ([#66](https://github.com/TibiaData/tibiadata-api-go/pull/66) by [kamilon](https://github.com/kamilon))

### Fixed
* Fix race condition with TibiaDataRequestStruct ([#65](https://github.com/TibiaData/tibiadata-api-go/pull/65) by [kamilon](https://github.com/kamilon))
* Bump docker/build-push-action from 2.7.0 to 2.8.0 ([#68](https://github.com/TibiaData/tibiadata-api-go/pull/68) by [dependabot](https://github.com/dependabot))

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

[3.5.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.4.1...v3.5.0
[3.4.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.4.0...v3.4.1
[3.4.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.3.1...v3.4.0
[3.3.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.3.0...v3.3.1
[3.3.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.2.2...v3.3.0
[3.2.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.2.1...v3.2.2
[3.2.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.2.0...v3.2.1
[3.2.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.1.1...v3.2.0
[3.1.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.1.0...v3.1.1
[3.1.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v3.0.0...v3.1.0
[3.0.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.6.2...v3.0.0
[0.6.2]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/tibiadata/tibiadata-api-go/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/tibiadata/tibiadata-api-go/compare/30f328f...v0.1.0
