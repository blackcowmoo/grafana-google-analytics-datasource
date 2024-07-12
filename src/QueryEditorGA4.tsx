import { QueryEditorProps, SelectableValue } from '@grafana/data';
import {
  AsyncMultiSelect,
  AsyncSelect,
  Badge,
  ButtonCascader,
  CascaderOption,
  HorizontalGroup,
  InlineFormLabel,
  InlineLabel,
  RadioButtonGroup,
} from '@grafana/ui';
import { DataSource } from 'DataSource';
import { DimensionFilter } from 'Filter';
import _ from 'lodash';
import React, { PureComponent } from 'react';
import { GADataSourceOptions, GAQuery } from 'types';
type Props = QueryEditorProps<DataSource, GAQuery, GADataSourceOptions>;

const defaultCacheDuration = 300;
const gaVersionBadge = {
  v4: {
    text: 'GA4(alpha)',
    tootip: 'experimental support',
  },
} as const;
const queryMode = [
  { label: 'Time Series', value: 'timeSeries' },
  { label: 'Table', value: 'table' },
  { label: 'Realtime', value: 'realtime' },
] as Array<SelectableValue<string>>;

const gaServiceLevelBadge = {
  GOOGLE_ANALYTICS_STANDARD: {
    text: 'Standard',
    tootip: 'max realtime query 30min',
  },
  GOOGLE_ANALYTICS_360: {
    text: 'Premium',
    tootip: 'max realtime query 60min',
  },
};

export const GAServiceLevel = {
  ServiceLevelStandard: 'GOOGLE_ANALYTICS_STANDARD',
  ServiceLevelPremium: 'GOOGLE_ANALYTICS_360',
  ServiceLevelUnspecified: 'SERVICE_LEVEL_UNSPECIFIED',
};

export class QueryEditorGA4 extends PureComponent<Props> {
  options: CascaderOption[] = [];
  constructor(props: Readonly<Props>) {
    super(props);
    const { query } = this.props;
    console.log('query.mode', query.mode);
    if (!query.hasOwnProperty('cacheDurationSeconds')) {
      this.props.query.cacheDurationSeconds = defaultCacheDuration;
    }
    this.props.query.version = props.datasource.getGaVersion();
    this.props.query.displayName = new Map<string, string>();
    this.props.datasource.getAccountSummaries().then((accountSummaries) => {
      this.options = accountSummaries;
      this.props.onChange(this.props.query);
    });
    if (query.mode === undefined || query.mode === '') {
      query.mode = 'timeSeries';
    }
    if (query.dimensionFilter === undefined) {
      query.dimensionFilter = {};
    }
  }

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

  onTimeDimensionChange = (value: SelectableValue<string>) => {
    console.log('value', value);
    const { query, onChange } = this.props;

    let timeDimension = value?.value || '';

    console.log(`timeDimension`, timeDimension);

    onChange({ ...query, timeDimension, selectedTimeDimensions: value });
    this.willRunQuery();
  };

  onIdSelect = async (value: string[], selectedOptions: CascaderOption[]) => {
    const [account, property, profile] = value;
    const { query, onChange, datasource } = this.props;
    const timezone = await datasource.getTimezone(account, property, profile);
    const serviceLevel = await datasource.getServiceLevel(account, property);

    onChange({ ...query, accountId: account, webPropertyId: property, profileId: profile, timezone, serviceLevel });
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

  onModeChange = (value: string) => {
    const { query, onChange } = this.props;
    switch (value) {
      case 'realtime':
        query.timeDimension = '';
        query.selectedTimeDimensions = {};
    }
    onChange({ ...query, mode: value });
    this.willRunQuery();
  };

  willRunQuery = _.debounce(() => {
    const { query, onRunQuery } = this.props;
    const { webPropertyId, metrics, timeDimension, mode } = query;
    console.log(`willRunQuery`);
    console.log(`query`, query);
    if (webPropertyId && metrics && (mode === 'table' || mode === 'realtime' || timeDimension)) {
      console.log(`onRunQuery`);
      onRunQuery();
    }
  }, 500);

  setDisplayName = (key: string, value = '') => {
    const {
      query: { displayName },
    } = this.props;
    displayName.set(key, value);
  };
  getDisplayName = (key: string) => {
    const {
      query: { displayName },
    } = this.props;
    if (displayName.has(key)) {
      return displayName.get(key);
    }
    return '';
  };
  render() {
    const { query, datasource } = this.props;
    const {
      accountId,
      webPropertyId,
      selectedTimeDimensions,
      selectedMetrics,
      selectedDimensions,
      timezone,
      mode,
      serviceLevel,
    } = query;
    const parsedWebPropertyId = webPropertyId?.split('/')[1];
    let serviceLevelBadge;
    switch (serviceLevel) {
      case GAServiceLevel.ServiceLevelPremium:
        serviceLevelBadge = (
          <Badge
            color="orange"
            text={gaServiceLevelBadge.GOOGLE_ANALYTICS_360.text}
            tooltip={gaServiceLevelBadge.GOOGLE_ANALYTICS_360.tootip}
            icon="google"
          ></Badge>
        );
        break;
      case GAServiceLevel.ServiceLevelStandard:
        serviceLevelBadge = (
          <Badge
            color="orange"
            text={gaServiceLevelBadge.GOOGLE_ANALYTICS_STANDARD.text}
            tooltip={gaServiceLevelBadge.GOOGLE_ANALYTICS_STANDARD.tootip}
            icon="google"
          ></Badge>
        );
        break;
      default:
        serviceLevelBadge = <Badge color="orange" text="Unknown" tooltip="Unknown" icon="google"></Badge>;
    }
    return (
      <>
        <div className="gf-form-group">
          <div className="gf-form">
            <HorizontalGroup spacing="sm" justify="flex-start">
              <ButtonCascader options={this.options} onChange={this.onIdSelect}>
                Account Select
              </ButtonCascader>
              <InlineLabel aria-label='account-info'>{`Account: ${accountId || ''},Property: ${webPropertyId || ''}`}</InlineLabel>
              <InlineLabel className="query-keyword" width={'auto'} tooltip={<>GA timeZone</>}>
                Timezone
              </InlineLabel>
              <InlineLabel width="auto">{timezone ? timezone : 'determined by profileId'}</InlineLabel>
              <Badge
                color="orange"
                text={gaVersionBadge.v4.text}
                tooltip={gaVersionBadge.v4.tootip}
                icon="google"
              ></Badge>
              {serviceLevelBadge}
            </HorizontalGroup>
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
              loadOptions={(q) => {
                console.log('load')
                if(mode === "realtime"){
                  return datasource.getRealtimeMetrics(q,parsedWebPropertyId)
                }
               return datasource.getMetrics(q, parsedWebPropertyId)}
              }
              placeholder={'ga:sessions'}
              value={selectedMetrics}
              onChange={this.onMetricChange}
              backspaceRemovesValue
              noOptionsMessage={'Search Metrics'}
              defaultOptions
              menuPlacement="bottom"
              isClearable
              key={mode+webPropertyId+"metrics"}
              aria-label='metrics'
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
              noOptionsMessage={'Search Dimension'}
              defaultOptions
              menuPlacement="bottom"
              isClearable
              disabled={mode === 'realtime'}
              aria-label='time-dimension'
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
              loadOptions={(q) => {
                if(mode === "realtime"){
                  return datasource.getRealtimeDimensions(q,null,parsedWebPropertyId)
                }
                return datasource.getDimensionsExcludeTimeDimensions(q, parsedWebPropertyId);
              }}
              placeholder={'ga:country'}
              value={selectedDimensions}
              onChange={this.onDimensionChange}
              backspaceRemovesValue
              noOptionsMessage={'Search Dimension'}
              defaultOptions
              menuPlacement="bottom"
              isClearable
              key={mode+parsedWebPropertyId+"dimensions"}
              aria-label='dimensions'
            />
          </div>
          <div className="gf-form">
            <InlineFormLabel
              className="query-keyword"
              tooltip={
                <>
                  Currently, only <code>or groups</code> are supported.
                </>
              }
            >
              DimensionFilter
            </InlineFormLabel>
            <DimensionFilter props={this.props}></DimensionFilter>
          </div>
          <div className="gf-form">
            <InlineFormLabel className="query-keyword">Query Mode</InlineFormLabel>
            <RadioButtonGroup options={queryMode} onChange={this.onModeChange} value={mode} />
          </div>
        </div>
      </>
    );
  }
}
