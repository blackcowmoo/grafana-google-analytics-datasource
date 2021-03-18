import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { AsyncMultiSelect, InlineFormLabel, SegmentAsync } from '@grafana/ui';
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
    const { query, onChange } = this.props;

    if (!item.value) {
      return;
    }

    onChange({ ...query, profileId: item.value });
    this.willRunQuery();
  };

  onAccountIdChange = (item: any) => {
    const { query, onChange } = this.props;

    if (!item.value) {
      return;
    }

    onChange({ ...query, accountId: item.value });
    this.willRunQuery();
  };

  onWebPropertyIdChange = (item: any) => {
    const { query, onChange } = this.props;

    if (!item.value) {
      return;
    }

    onChange({ ...query, webPropertyId: item.value });
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

  willRunQuery = () => {
    const { query, onRunQuery } = this.props;
    const { profileId, metrics, dimensions } = query;
    console.log(`willRunQuery`);
    console.log(`query`, query);
    if (profileId && metrics && dimensions) {
      console.log(`onRunQuery`);
      onRunQuery();
    }
  };

  render() {
    const { query, datasource } = this.props;
    const {
      accountId,
      webPropertyId,
      profileId,
      selectedMetrics: selectMetrics,
      selectedDimensions: selectDimensions,
    } = query;
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
              <p>
                The <code>profileId</code> is used to identify which GoogleAnalytics is to be accessed or altered. This
              </p>
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
            value={selectMetrics}
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
            value={selectDimensions}
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
