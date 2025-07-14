import React from 'react';
import {
  QueryEditorProps,
  SelectableValue,
} from '@grafana/data';
import {
  Input,
  InlineLabel,
  AsyncMultiSelect,
  AsyncSelect,
  Badge,
  HorizontalGroup,
  InlineFormLabel,
} from '@grafana/ui';
import _ from 'lodash';
import { DataSource } from 'DataSource';
import { GAQuery, GADataSourceOptions } from './types';

interface Props extends QueryEditorProps<DataSource, GAQuery, GADataSourceOptions> {
  query: GAQuery;
  onChange: (query: GAQuery) => void;
  onRunQuery: () => void;
  datasource: DataSource;
}

const defaultCacheDuration = 300;
const badgeMap = {
  "v3": {
    "text": "UA",
    "tootip": "2023/07/01 no more data collect"
  },
  "v4": {
    "text": "GA4(alpha)",
    "tootip": "experimental support"
  },
} as const

export class QueryEditorUA extends React.PureComponent<Props> {
  constructor(props: Readonly<Props>) {
    super(props);
    const { query, datasource, onChange } = props;

    if (!('cacheDurationSeconds' in query)) {
      query.cacheDurationSeconds = defaultCacheDuration;
      query.filtersExpression = '';
    }

    query.version = datasource.getGaVersion();
    query.displayName = new Map<string, string>();

    // Ensure modified query is propagated
    onChange(query);
  }

  onMetricChange = (items: Array<SelectableValue<string>>) => {
    const { query, onChange } = this.props;

    let metrics = [] as string[];
    items.map((item) => {
      if (item.value) {
        metrics.push(item.value);
      }
    });


    onChange({ ...query, selectedMetrics: items, metrics });
    this.willRunQuery();
  };

  onTimeDimensionChange = (item: any) => {
    const { query, onChange } = this.props;

    let timeDimension = item.value;



    onChange({ ...query, timeDimension, selectedTimeDimensions: item });
    this.willRunQuery();
  };

  onAccountIdChange = (value: string) => {
    const { query, onChange } = this.props;
    onChange({ ...query, accountId: value });
    this.refreshTimezone(value, query.webPropertyId, query.profileId);
    this.willRunQuery();
  };

  onWebPropertyIdChange = (value: string) => {
    const { query, onChange } = this.props;
    onChange({ ...query, webPropertyId: value });
    this.refreshTimezone(query.accountId, value, query.profileId);
    this.willRunQuery();
  };

  onProfileIdChange = (value: string) => {
    const { query, onChange } = this.props;
    onChange({ ...query, profileId: value });
    this.refreshTimezone(query.accountId, query.webPropertyId, value);
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


    if (
      (profileId && metrics && timeDimension && version === "v3") ||
      (webPropertyId && metrics && timeDimension && version === "v4")
    ) {

      onRunQuery();
    }
  }, 500);

  refreshTimezone = (accountId?: string, webPropertyId?: string, profileId?: string) => {
    const { query, onChange, datasource } = this.props;
    if (
      accountId &&
      webPropertyId &&
      profileId &&
      !accountId.includes('$') &&
      !webPropertyId.includes('$') &&
      !profileId.includes('$')
    ) {
      datasource.getTimezone(accountId, webPropertyId, profileId).then((timezone) => {
        onChange({ ...query, timezone });
      });
    } else {
      // 템플릿 변수를 사용하는 경우, timezone 은 쿼리 실행 시점에 결정된다.
      onChange({ ...query, timezone: '' });
    }
  };

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

    const parsedWebPropertyId = webPropertyId?.split('/')[1]
    return (
      <>
        <div className="gf-form-group">
          <div className="gf-form">
            <HorizontalGroup spacing="sm" justify='flex-start' >
              <Input
                width={20}
                value={accountId || ''}
                placeholder="$accountId"
                onChange={(e) => this.onAccountIdChange(e.currentTarget.value)}
              />
              <Input
                width={30}
                value={webPropertyId || ''}
                placeholder="$webPropertyId"
                onChange={(e) => this.onWebPropertyIdChange(e.currentTarget.value)}
              />
              <Input
                width={30}
                value={profileId || ''}
                placeholder="$profileId"
                onChange={(e) => this.onProfileIdChange(e.currentTarget.value)}
              />
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
          </div>

          <div className="gf-form" key={parsedWebPropertyId}>
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
              loadOptions={(q) => datasource.getMetrics(q, parsedWebPropertyId)}
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
              loadOptions={(q) => datasource.getDimensionsExcludeTimeDimensions(q, parsedWebPropertyId)}
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
