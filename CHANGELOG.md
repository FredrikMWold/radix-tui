# Changelog

## [1.9.0](https://github.com/FredrikMWold/radix-tui/compare/v1.8.0...v1.9.0) (2025-12-15)


### Features

* Add context switcher to application dashboard ([b19034b](https://github.com/FredrikMWold/radix-tui/commit/b19034b1c19979b98fdbcb4ac3e0465a4c1a3244))

## [1.8.0](https://github.com/FredrikMWold/radix-tui/compare/v1.7.0...v1.8.0) (2025-08-22)


### Features

* adds polling off pipeline jobs ([8df5847](https://github.com/FredrikMWold/radix-tui/commit/8df584763a28eb8746a6c112182661b37983dd62))
* cache list of application ([63bbd30](https://github.com/FredrikMWold/radix-tui/commit/63bbd303af483e9010ab806fcba8b5c1ac31f1fe))


### Bug Fixes

* imporve login flow ([2a07116](https://github.com/FredrikMWold/radix-tui/commit/2a07116cac7493c43a9a7f605a9074e2ad8cec12))
* show correct timezone ([8e1a727](https://github.com/FredrikMWold/radix-tui/commit/8e1a727c005a8b1e1cfc0468a7fc284ccc6a2075))

## [1.7.0](https://github.com/FredrikMWold/radix-tui/compare/v1.6.1...v1.7.0) (2024-07-23)


### Features

* adds error screen when terminal is to small ([15fdbca](https://github.com/FredrikMWold/radix-tui/commit/15fdbca25299fe5db215f084d8b8e69124aeeedf))
* show which context the radix cli is using ([15fdbca](https://github.com/FredrikMWold/radix-tui/commit/15fdbca25299fe5db215f084d8b8e69124aeeedf))

## [1.6.1](https://github.com/FredrikMWold/radix-tui/compare/v1.6.0...v1.6.1) (2024-07-22)


### Bug Fixes

* extra comma in the config file causing app to crash ([82f0c3d](https://github.com/FredrikMWold/radix-tui/commit/82f0c3dcd6cd255219e1db0c240633ffa9e7916c))

## [1.6.0](https://github.com/FredrikMWold/radix-tui/compare/v1.5.0...v1.6.0) (2024-07-22)


### Features

* adds ability to create apply-config pipeline ([99e978a](https://github.com/FredrikMWold/radix-tui/commit/99e978a7dc1336c1f4812309b291be450f574429))

## [1.5.0](https://github.com/FredrikMWold/radix-tui/compare/v1.4.0...v1.5.0) (2024-07-18)


### Features

* adds ability to create new build-deploy pipelines ([b84cbd3](https://github.com/FredrikMWold/radix-tui/commit/b84cbd384df98c7065e7742f4f8d172d45c0494c))


### Bug Fixes

* adds refresh help text ([cc0dff0](https://github.com/FredrikMWold/radix-tui/commit/cc0dff034a28d3123fe6e62873ec97e19bbe37e9))
* reshing in form would cause wrong help text to be displayed ([cc0dff0](https://github.com/FredrikMWold/radix-tui/commit/cc0dff034a28d3123fe6e62873ec97e19bbe37e9))

## [1.4.0](https://github.com/FredrikMWold/radix-tui/compare/v1.3.0...v1.4.0) (2024-07-18)


### Features

* adds helper text to the different pages ([05ea337](https://github.com/FredrikMWold/radix-tui/commit/05ea337db8ed8dccc797a88995baf77d5f4c9b80))

## [1.3.0](https://github.com/FredrikMWold/radix-tui/compare/v1.2.0...v1.3.0) (2024-07-18)


### Features

* removes auto update application user needs to press ctrl+r ([3aab288](https://github.com/FredrikMWold/radix-tui/commit/3aab288021823c3be68b6a72a78284b710fb21c3))

## [1.2.0](https://github.com/FredrikMWold/radix-tui/compare/v1.1.1...v1.2.0) (2024-07-16)


### Features

* adds enter and escape as naviagtion options between pipelines and aplications ([fd9c2d0](https://github.com/FredrikMWold/radix-tui/commit/fd9c2d0e1382b3a3632c50e17d7ac0f1429ed4dc))


### Bug Fixes

* bug that would spam radix api more when swapping between application ([fd9c2d0](https://github.com/FredrikMWold/radix-tui/commit/fd9c2d0e1382b3a3632c50e17d7ac0f1429ed4dc))

## [1.1.1](https://github.com/FredrikMWold/radix-tui/compare/v1.1.0...v1.1.1) (2024-07-16)


### Bug Fixes

* only start auto update tick after user selects a application ([591e8f6](https://github.com/FredrikMWold/radix-tui/commit/591e8f648b354a9e468cd501da199e6810f34ed1))

## [1.1.0](https://github.com/FredrikMWold/radix-tui/compare/v1.0.1...v1.1.0) (2024-07-16)


### Features

* adds abilty to open pipeline job in browser ([5495981](https://github.com/FredrikMWold/radix-tui/commit/549598180a1bfeeb9446bbe2f4a06de8ebf4af21))
* show selected application as pipeline table heading ([88d5247](https://github.com/FredrikMWold/radix-tui/commit/88d5247d1edacec2adc780ffe6c8d7967eacd890))


### Bug Fixes

* fix bug that stopped auto refresh when pipelineTable is focused ([0904cc6](https://github.com/FredrikMWold/radix-tui/commit/0904cc68555290eecf174fbad970c41d0f6b7e6f))
* handle case when a pipeline job has no environments ([5b79b6e](https://github.com/FredrikMWold/radix-tui/commit/5b79b6ec7ef2da1361824d79d4312cee001237ba))

## [1.0.1](https://github.com/FredrikMWold/radix-tui/compare/v1.0.0...v1.0.1) (2024-07-15)


### Bug Fixes

* handles login from rx-cli when launching the application ([8e27a35](https://github.com/FredrikMWold/radix-tui/commit/8e27a35b1a58216a7ae01d6538b6947c0cba2a9e))

## 1.0.0 (2024-07-15)


### Miscellaneous Chores

* release v1.0.0 ([8ea290f](https://github.com/FredrikMWold/radix-tui/commit/8ea290f5485b376ba764a7546620c6a70a19d7e7))
