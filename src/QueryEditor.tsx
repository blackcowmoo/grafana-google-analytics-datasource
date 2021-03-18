import { QueryEditorProps } from '@grafana/data';
import { InlineFormLabel, LinkButton, SegmentAsync } from '@grafana/ui';
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

    const v = item.value;
    onChange({ ...query, profileId: v });
  };

  onAccountIdChange = (item: any) => {
    const { query, onChange } = this.props;

    if (!item.value) {
      return;
    }

    const v = item.value;
    onChange({ ...query, accountId: v });
  };

  onWebPropertyIdChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;

    if (!item.value) {
      return;
    }

    const v = item.value;
    onChange({ ...query, webPropertyId: v });
    onRunQuery();
  };

  onMetricChange = (item: any) => {
    const { query, onChange } = this.props;

    if (!item.value) {
      return;
    }

    const v = item.value;
    onChange({ ...query, metric: v });
  };

  onDimensionChange = (item: any) => {
    const { query, onChange } = this.props;

    if (!item.value) {
      return;
    }

    const v = item.value;
    onChange({ ...query, dimension: v });
  };

  render() {
    const { query, datasource } = this.props;
    const { accountId, webPropertyId, profileId, metric, dimension } = query;
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
          {accountId && <LinkButton style={{ marginTop: 1 }} variant="link" icon="link" target="_blank"></LinkButton>}
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
          {webPropertyId && (
            <LinkButton style={{ marginTop: 1 }} variant="link" icon="link" target="_blank"></LinkButton>
          )}
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
          {profileId && <LinkButton style={{ marginTop: 1 }} variant="link" icon="link" target="_blank"></LinkButton>}
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
            Metric
          </InlineFormLabel>
          <SegmentAsync
            loadOptions={() => datasource.getMetrics()}
            placeholder={'ga:sessions'}
            value={metric}
            allowCustomValue={true}
            onChange={this.onMetricChange}
          ></SegmentAsync>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>dimension</code> ga:*
              </p>
            }
          >
            Dimension
          </InlineFormLabel>
          <SegmentAsync
            loadOptions={() => datasource.getDimensions()}
            placeholder={'ga:dateHourMinute'}
            value={dimension}
            allowCustomValue={true}
            onChange={this.onDimensionChange}
          ></SegmentAsync>
        </div>
      </>
    );
  }
}
