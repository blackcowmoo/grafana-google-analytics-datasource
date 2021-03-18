import { QueryEditorProps, SelectableValue } from '@grafana/data';
import {
  AlphaNotice,
  AsyncMultiSelect,
  InlineFormLabel,
  SegmentAsync,
} from '@grafana/ui';
import { DataSource } from 'DataSource';
import React, { PureComponent } from 'react';
import { GADataSourceOptions, GAQuery } from 'types';

type Props = QueryEditorProps<DataSource, GAQuery, GADataSourceOptions>;

const defaultCacheDuration = 300;

export const formatCacheTimeLabel = (s: number = defaultCacheDuration) => {
  if (s < 60) {
    return s + 's';
  } else if (s < 3600) {
    return s / 60 + 'm';
  }

  return s / 3600 + 'h';
};

export class QueryEditor extends PureComponent<Props> {
  constructor(props: Readonly<Props>) {
    super(props);
    if (!this.props.query.hasOwnProperty('cacheDurationSeconds')) {
      this.props.query.cacheDurationSeconds = defaultCacheDuration;
    }
  }

  onProfileIdChange = (item: any) => {
    const {
      query,
      query: { metrics, dimensions, accountId, webPropertyId },
      onChange,
      datasource,
    } = this.props;
    const profileId = item.value as string;

    if (profileId)
      datasource.getProfileTimezone(accountId, webPropertyId, profileId).then((timezone) => {
        const { query, onChange } = this.props;
        console.log(`timezone`, timezone)
        onChange({ ...query, timezone });
        this.willRunQuery(profileId, metrics, dimensions);
      });
    onChange({ ...query, profileId });
    this.willRunQuery(profileId, metrics, dimensions);
  };

  onAccountIdChange = (item: any) => {
    const {
      query,
      query: { profileId, metrics, dimensions },
      onChange,
    } = this.props;
    let accountId = item.value;

    onChange({ ...query, accountId });
    this.willRunQuery(profileId, metrics, dimensions);
  };

  onWebPropertyIdChange = (item: any) => {
    const {
      query,
      query: { profileId, metrics, dimensions },
      onChange,
    } = this.props;
    let webPropertyId = item.value;

    onChange({ ...query, webPropertyId });
    this.willRunQuery(profileId, metrics, dimensions);
  };

  onMetricChange = (items: Array<SelectableValue<string>>) => {
    const {
      query,
      query: { profileId, dimensions },
      onChange,
    } = this.props;

    let metrics = [] as string[];
    items.map((item) => {
      if (item.value) {
        metrics.push(item.value);
      }
    });
    console.log(`metrics`, metrics);

    onChange({ ...query, selectedMetrics: items, metrics });
    this.willRunQuery(profileId, metrics, dimensions);
  };

  onDimensionChange = (items: Array<SelectableValue<string>>) => {
    const {
      query,
      query: { profileId, metrics },
      onChange,
    } = this.props;
    let dimensions = [] as string[];
    items.map((item) => {
      if (item.value) {
        dimensions.push(item.value);
      }
    });

    console.log(`dimensions`, dimensions);

    onChange({ ...query, selectedDimensions: items, dimensions });
    this.willRunQuery(profileId, metrics, dimensions);
  };

  willRunQuery = (profileId: string, metrics: string[], dimensions: string[]) => {
    const { query, onRunQuery } = this.props;
    console.log(`willRunQuery`);
    console.log(`query`, query);
    if (profileId && metrics && dimensions) {
      console.log(`onRunQuery`);
      onRunQuery();
    }
  };

  render() {
    const { query, datasource } = this.props;
    const { accountId, webPropertyId, profileId, selectedMetrics, selectedDimensions, timezone } = query;
    return (
      <>
        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>accountId</code> is used to identify which GoogleAnalytics is to be accessed or altered.
              </p>
            }
          >
            Account ID
          </InlineFormLabel>
          <SegmentAsync
            loadOptions={() => datasource.getAccountIds()}
            placeholder="Enter Account ID"
            value={accountId}
            allowCustomValue={true}
            onChange={this.onAccountIdChange}
          ></SegmentAsync>
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>webPropertyId</code> is used to identify which GoogleAnalytics is to be accessed or altered.
              </p>
            }
          >
            Web Property ID
          </InlineFormLabel>
          <SegmentAsync
            loadOptions={() => datasource.getWebPropertyIds(accountId)}
            placeholder="Enter Web Property ID"
            value={webPropertyId}
            allowCustomValue={true}
            onChange={this.onWebPropertyIdChange}
          ></SegmentAsync>
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <div>
                The <code>profileId</code> is used to identify which GoogleAnalytics is to be accessed or altered. This
              </div>
            }
          >
            Profile ID
          </InlineFormLabel>
          <SegmentAsync
            loadOptions={() => datasource.getProfileIds(accountId, webPropertyId)}
            placeholder="Enter Profile ID"
            value={profileId}
            allowCustomValue={true}
            onChange={this.onProfileIdChange}
          ></SegmentAsync>

          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>metric</code> ga:*
              </p>
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
          ></AsyncMultiSelect>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code> dimensions </code> At least one ga:date* is required.
              </p>
            }
          >
            Dimension
          </InlineFormLabel>
          <AsyncMultiSelect
            loadOptions={(q) => datasource.getDimensions(q)}
            placeholder={'ga:dateHour'}
            value={selectedDimensions}
            onChange={this.onDimensionChange}
            backspaceRemovesValue
            cacheOptions
            noOptionsMessage={'Search Dimension'}
          ></AsyncMultiSelect>
        </div>
      </>
    );
  }
}
