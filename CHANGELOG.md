# Changelog

## 0.3.2
- Remove deprecated GA3 (Universal Analytics) datasource — GA4 only going forward([#165](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/165))
- Complete GA4 filter UI overhaul with MetricFilter support([#169](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/169))
- Support accounts/properties as dashboard variables([#167](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/167))
- Add dashboard variable support for webPropertyId([#152](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/152))
- Fix expanding multi-value variables in filter expressions([#154](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/154))
- Fix sub-day Grafana time range not respected in time-series output([#156](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/156))
- Fix rows with unparseable time dimension causing failures([#155](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/155))

## 0.3.1
- Support or,and,not filter[#128](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/128)
- Support grafana version < 13[#136](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/136)
- Fix github actions [#136](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/136)
- Update README.md[#127](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/127)
- Remove upper bound of your grafana dependency

## 0.3.0
- Fix timezone no zoneinfo.zip([#102](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/102))
- Supprot variable at Dimensions Filter([#103](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/103))
- Support realtime Report([#107](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/107))
- Fix error message typo([#110](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/110))
- Replace deprecated Grafana(v12) SCSS styles([#114](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/114))

## 0.2.3
- Support DimensionFilter([#90](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/90))
- Fix metric type bug([#88](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/88))
- Add time series and table query mode([#87](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/87)) 

## 0.2.2
- Hotfix 0.2.1 query editor frontend bug

## 0.2.1
- Fix healthcheck fail when no profile or no property([#80](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/80))
- Update README to remove hardcoded GCP Project ID([#79](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/79))
- Support custom metrics and dimensions([#69](https://github.com/blackcowmoo/grafana-google-analytics-datasource/pull/69/files))

## 0.2.0
- support GA4
- change query editor UI
- go version 1.18 -> 1.20
- update dependencies
## 0.1.5
- support grafana version < 11
- update dependencies
## 0.1.4
- apply datasource intance management
- upgrade go dependencies
- support grafana 8.x version
## 0.1.3
- fix query editor ui
- add filter expression
## 0.1.2
- separate time dimension and other dimensions
- fix bug multi dimensions breaking graphs

## 0.1.1
- enhancement query editor ui
- add timezone label at query editor
- get metric bug fix
- add Metrics & Dimensions description
- add Metadata type for typescript
- remove unused function


## 0.1.0

Initial release. Not fit for production use.
