import { QueryEditorProps } from '@grafana/data';
import { Alert, InlineFormLabel, LinkButton, SegmentAsync, SegmentInput } from '@grafana/ui';
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

export const checkDateForm = (date: string): boolean => {
  const exp = new RegExp('^[0-9]{4}-[0-9]{2}-[0-9]{2}$|^today$|^yesterday$|^[0-9]+(daysAgo)$');
  return exp.test(date);
};

export class QueryEditor extends PureComponent<Props> {
  componentWillMount() {
    if (!this.props.query.hasOwnProperty('cacheDurationSeconds')) {
      this.props.query.cacheDurationSeconds = defaultCacheDuration;
    }
  }

  onStartDateChange = (event: React.ReactText) => {
    const { query, onChange } = this.props;
    const date = event.toString();
    const result = checkDateForm(date);
    console.log('start', date, result);
    if (result) {
      onChange({ ...query, startDate: date });
    }
  };

  onEndDateChange = (event: React.ReactText) => {
    const { query, onChange } = this.props;
    const date = event.toString();
    const result = checkDateForm(date);
    console.log('end', date, result);
    if (result) {
      onChange({ ...query, endDate: date });
    }
  };

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

  onMetricsChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;
    const v = item.split(',');

    onChange({ ...query, metrics: v });
    onRunQuery();
  };

  onDimensionsChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;
    const v = item.split(',');

    onChange({ ...query, dimensions: v });
    onRunQuery();
  };

  onSortChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;
    const v = item.split(',');

    onChange({ ...query, sort: v });
    onRunQuery();
  };

  render() {
    const { query, datasource } = this.props;
    const { accountId, webPropertyId, profileId, startDate, endDate, metrics, dimensions, sort } = query;
    return (
      <>
        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>accountId</code> is used to identify which GoogleAnalytics is to be accessed or altered. This
                ID is the value between the "/d/" and the "/edit" in the URL of your GoogleAnalytics.
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
                The <code>profileId</code> is used to identify which GoogleAnalytics is to be accessed or altered. This
                ID is the value between the "/d/" and the "/edit" in the URL of your GoogleAnalytics.
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
                ID is the value between the "/d/" and the "/edit" in the URL of your GoogleAnalytics.
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
                The <code>StartDate</code> data format regular expression is
                <code>
                  [0-9]{4}-[0-9]{2}-[0-9]{2}|today|yesterday|[0-9]+(daysAgo)
                </code>
              </p>
            }
          >
            Start Date
          </InlineFormLabel>
          <SegmentInput onChange={this.onStartDateChange} value={startDate} placeholder={'YYYY-MM-DD'}></SegmentInput>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>EndDate</code> data format regular expression is
                <code>
                  [0-9]{4}-[0-9]{2}-[0-9]{2}|today|yesterday|[0-9]+(daysAgo)
                </code>
              </p>
            }
          >
            End Date
          </InlineFormLabel>
          <SegmentInput onChange={this.onEndDateChange} value={endDate} placeholder={'YYYY-MM-DD'}></SegmentInput>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>metrics</code> ga:*
              </p>
            }
          >
            Metrics
          </InlineFormLabel>
          <SegmentInput
            onChange={this.onMetricsChange}
            value={metrics ? metrics.toString() : ''}
            placeholder={'ga:sessions'}
          ></SegmentInput>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>dimensions</code> ga:*
              </p>
            }
          >
            Dimensions
          </InlineFormLabel>
          <SegmentInput
            onChange={this.onDimensionsChange}
            value={dimensions ? dimensions.toString() : ''}
            placeholder={'ga:dateHourMinute'}
          ></SegmentInput>
        </div>

        <div className="gf-form-inline">
          <InlineFormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>sort</code> asc = ga:* , desc = -ga:*
              </p>
            }
          >
            Sort
          </InlineFormLabel>
          <SegmentInput
            onChange={this.onSortChange}
            value={sort ? sort.toString() : ''}
            placeholder={'ga:dateHourMinute'}
          ></SegmentInput>
        </div>
        {!startDate && <Alert title={'Start Date Invalid'}></Alert>}
        {!endDate && <Alert title={'End Date Invalid'}></Alert>}
      </>
    );
  }
}
