import { QueryEditorProps, SelectableValue } from '@grafana/data';
import {
  AsyncMultiSelect,
  AsyncSelect,
  Badge,
  Cascader,
  CascaderOption,
  HorizontalGroup,
  InlineFormLabel,
  InlineLabel,
  Input,
  SegmentAsync
} from '@grafana/ui';
import { DataSource } from 'DataSource';
import _ from 'lodash';
import React, { PureComponent } from 'react';
import { GADataSourceOptions, GAQuery } from 'types';

type Props = QueryEditorProps<DataSource, GAQuery, GADataSourceOptions>;

const defaultCacheDuration = 300;
const badgeMap = {
  "v3": {
    "text": "UA",
    "tootip": "2023/07/01 no more data collect"
  },
  "v4": {
    "text": "GA(alpha)",
    "tootip": "experimental support"
  },
} as const

export class QueryEditor extends PureComponent<Props> {
  options: CascaderOption[] = []
  constructor(props: Readonly<Props>) {
    super(props);
    if (!this.props.query.hasOwnProperty('cacheDurationSeconds')) {
      this.props.query.cacheDurationSeconds = defaultCacheDuration;
      this.props.query.filtersExpression = '';
    }
    this.props.query.version = props.datasource.getGaVersion()
    this.props.query.displayName = new Map<string, string>()
    this.props.datasource.getAccountSummaries().then((accountSummaries) => {
      this.options = accountSummaries
      this.props.onChange(this.props.query)
    })
  }

  onProfileIdChange = (item: SelectableValue<string>) => {
    const {
      query,
      query: { version, accountId, webPropertyId },
      onChange,
      datasource,
    } = this.props;
    const profileId = item.value as string;
    this.setDisplayName(profileId, item.label)

    if (profileId && version === "v3") {
      datasource.getProfileTimezone(accountId, webPropertyId, profileId).then((timezone) => {
        const { query, onChange } = this.props;
        console.log(`timezone`, timezone);
        onChange({ ...query, timezone });
        this.willRunQuery();
      });
    }
    onChange({ ...query, profileId });
    this.willRunQuery();
  };

  onAccountIdChange = (item: SelectableValue<string>) => {
    const { query, query: { displayName }, onChange } = this.props;
    let accountId = item.value ?? "";
    this.setDisplayName(accountId, item.label)
    onChange({ ...query, accountId, displayName });
    this.willRunQuery();
  };

  onWebPropertyIdChange = (item: any) => {
    const { query, query: { version, accountId }, onChange, datasource } = this.props;
    let webPropertyId = item.value;
    console.log('onWebPropertyIdChange', webPropertyId)
    if (webPropertyId && version === "v4") {
      console.log('gettimezone')
      datasource.getProfileTimezone(accountId, webPropertyId, "").then((timezone) => {
        const { query, onChange } = this.props;
        console.log(`timezone`, timezone);
        onChange({ ...query, timezone });
        this.willRunQuery();
      });
    }
    onChange({ ...query, webPropertyId });
    this.setDisplayName(webPropertyId, item.label)
    this.willRunQuery();
  };

  onMetricChange = (items: Array<SelectableValue<string>>) => {
    const { query, onChange } = this.props;

    let metrics = [] as string[];
    items.map((item) => {
      if (item.value) {
        metrics.push(item.value);
      }
    });
    console.log(`metrics`, metrics);

    onChange({ ...query, selectedMetrics: items, metrics });
    this.willRunQuery();
  };

  onTimeDimensionChange = (item: any) => {
    const { query, onChange } = this.props;

    let timeDimension = item.value;

    console.log(`timeDimension`, timeDimension);

    onChange({ ...query, timeDimension, selectedTimeDimensions: item });
    this.willRunQuery();
  };
  onSelect = (item: string) => {
    const { query, query: { version, accountId }, onChange, datasource } = this.props;
    let webPropertyId = item
    console.log('onWebPropertyIdChange', webPropertyId)
    if (webPropertyId && version === "v4") {
      console.log('gettimezone')
      datasource.getProfileTimezone(accountId, webPropertyId, "").then((timezone) => {
        const { query, onChange } = this.props;
        console.log(`timezone`, timezone);
        onChange({ ...query, timezone });
        this.willRunQuery();
      });
    }
    onChange({ ...query, webPropertyId });
    this.setDisplayName(webPropertyId, item)
    this.willRunQuery();
  };
  onDimensionChange = (items: Array<SelectableValue<string>>) => {
    const { query, onChange } = this.props;
    let dimensions = [] as string[];
    items.map((item) => {
      if (item.value) {
        dimensions.push(item.value);
      }
    });

    console.log(`dimensions`, dimensions);

    onChange({ ...query, selectedDimensions: items, dimensions });
    this.willRunQuery();
  };

  onFiltersExpressionChange = (item: any, ...t: any) => {
    const { query, onChange } = this.props;
    let { filtersExpression } = query;
    filtersExpression = item;

    onChange({ ...query, filtersExpression });
    this.willRunQuery();
  };

  willRunQuery = _.debounce(() => {
    const { query, onRunQuery } = this.props;
    const { version, webPropertyId, profileId, metrics, timeDimension } = query;
    console.log(`willRunQuery`);
    console.log(`query`, query);
    if (
      (profileId && metrics && timeDimension && version === "v3") ||
      (webPropertyId && metrics && timeDimension && version === "v4")
    ) {
      console.log(`onRunQuery`);
      onRunQuery();
    }
  }, 500);

  setDisplayName = (key: string, value = "") => {
    const { query: { displayName } } = this.props
    displayName.set(key, value)
  }
  getDisplayName = (key: string) => {
    const { query: { displayName } } = this.props
    if (displayName.has(key)) {
      return displayName.get(key)
    }
    return ""
  }
  render() {
    const { query, datasource } = this.props;
    console.log('query.version', query.version)
    const {
      accountId,
      webPropertyId,
      profileId,
      selectedTimeDimensions,
      selectedMetrics,
      selectedDimensions,
      timezone,
      filtersExpression,
      version
    } = query;
    console.log('webPropertyId', webPropertyId)
    return (
      <>
        <Cascader options={this.options} onSelect={this.onSelect}></Cascader>
        <div className="gf-form-group">
          <div className="gf-form">
            <HorizontalGroup spacing="none">
              <InlineFormLabel
                className="query-keyword"
                tooltip={
                  <>
                    The <code>accountId</code> is used to identify which GoogleAnalytics is to be accessed or altered.
                  </>
                }
              >
                Account ID
              </InlineFormLabel>
              <SegmentAsync
                loadOptions={() => datasource.getAccountIds()}
                placeholder="Enter Account ID"
                value={this.getDisplayName(accountId)}
                allowCustomValue
                onChange={this.onAccountIdChange}
              />
              <InlineFormLabel
                className="query-keyword"
                tooltip={
                  <>
                    The <code>webPropertyId</code> is used to identify which GoogleAnalytics is to be accessed or
                    altered.
                  </>
                }
              >
                Web Property ID
              </InlineFormLabel>
              <SegmentAsync
                loadOptions={() => datasource.getWebPropertyIds(accountId)}
                placeholder="Enter Web Property ID"
                value={this.getDisplayName(webPropertyId)}
                allowCustomValue
                onChange={this.onWebPropertyIdChange}
              />
              {version === "v3" &&
                <>
                  <InlineFormLabel
                    className="query-keyword"
                    tooltip={
                      <>
                        The <code>profileId</code> is used to identify which GoogleAnalytics is to be accessed or altered.
                      </>
                    }
                  >
                    Profile ID
                  </InlineFormLabel>
                  <SegmentAsync
                    loadOptions={() => datasource.getProfileIds(accountId, webPropertyId)}
                    placeholder="Enter Profile ID"
                    value={this.getDisplayName(profileId)}
                    allowCustomValue
                    onChange={this.onProfileIdChange}
                  />
                </>
              }
              <InlineLabel className="query-keyword" width={'auto'} tooltip={<>GA timeZone</>}>
                Timezone
              </InlineLabel>
              <InlineLabel width="auto">{timezone ? timezone : 'determined by profileId'}</InlineLabel>
              {
                version === "v3"
                &&
                <Badge color='red' text={badgeMap.v3.text} tooltip={badgeMap.v3.tootip} icon='google'></Badge>
              }
              {
                version === "v4"
                &&
                <Badge color='orange' text={badgeMap.v4.text} tooltip={badgeMap.v4.tootip} icon='google'></Badge>
              }
            </HorizontalGroup>

            <div className="gf-form-label gf-form-label--grow" />
          </div>

          <div className="gf-form">
            <InlineFormLabel
              className="query-keyword"
              tooltip={
                <>
                  The <code>metric</code> ga:*
                </>
              }
            >
              Metrics
            </InlineFormLabel>
            <AsyncMultiSelect
              loadOptions={(q) => datasource.getMetrics(q)}
              placeholder={'ga:sessions'}
              value={selectedMetrics}
              onChange={this.onMetricChange}
              backspaceRemovesValue
              cacheOptions
              noOptionsMessage={'Search Metrics'}
              defaultOptions
              menuPlacement="bottom"
              isClearable
            />

            <InlineFormLabel
              className="query-keyword"
              tooltip={
                <>
                  The <code>time dimensions</code> At least one ga:date* is required.
                </>
              }
            >
              Time Dimension
            </InlineFormLabel>
            <AsyncSelect
              loadOptions={() => datasource.getTimeDimensions()}
              placeholder={'ga:dateHour'}
              value={selectedTimeDimensions}
              onChange={this.onTimeDimensionChange}
              backspaceRemovesValue
              cacheOptions
              noOptionsMessage={'Search Dimension'}
              defaultOptions
              menuPlacement="bottom"
              isClearable
            />

            <InlineFormLabel
              className="query-keyword"
              tooltip={
                <>
                  The <code>dimensions</code> exclude time dimensions
                </>
              }
            >
              Dimensions
            </InlineFormLabel>
            <AsyncMultiSelect
              loadOptions={(q) => datasource.getDimensionsExcludeTimeDimensions(q)}
              placeholder={'ga:country'}
              value={selectedDimensions}
              onChange={this.onDimensionChange}
              backspaceRemovesValue
              cacheOptions
              noOptionsMessage={'Search Dimension'}
              defaultOptions
              menuPlacement="bottom"
              isClearable
            />
          </div>
          <div className="gf-form">
            <InlineFormLabel
              className="query-keyword"
              tooltip={
                <>
                  The <code>filter</code> dimensions and metrics
                </>
              }
            >
              Filters Expressions
            </InlineFormLabel>
            <Input
              value={filtersExpression}
              onChange={(e) => this.onFiltersExpressionChange(e.currentTarget.value)}
              placeholder="ga:pagePath==/path/to/page"
            />
          </div>
        </div>
      </>
    );
  }
}
