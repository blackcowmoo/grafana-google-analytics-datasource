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
  componentWillMount() {
    if (!this.props.query.hasOwnProperty('cacheDurationSeconds')) {
      this.props.query.cacheDurationSeconds = defaultCacheDuration;
    }
  }

  // onRangeChange = (event: ChangeEvent<HTMLInputElement>) => {
  //   this.props.onChange({
  //     ...this.props.query,
  //     range: event.target.value,
  //   });
  // };

  onViewIDChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;

    if (!item.value) {
      return; // ignore delete?
    }

    const v = item.value;
    // Check for pasted full URLs
    onChange({ ...query, profileId: v });
    onRunQuery();
  };

  onAccountIDChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;

    if (!item.value) {
      return; // ignore delete?
    }

    const v = item.value;
    // Check for pasted full URLs
    onChange({ ...query, accountId: v });
    onRunQuery();
  };

  onWebPropertyChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;

    if (!item.value) {
      return; // ignore delete?
    }

    const v = item.value;
    // Check for pasted full URLs
    onChange({ ...query, webPropertyId: v });
    onRunQuery();
  };
  render() {
    const { query, datasource } = this.props;
    const { accountId, webPropertyId, profileId: viewId } = query;
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
            onChange={this.onAccountIDChange}
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
                The <code>viewId</code> is used to identify which GoogleAnalytics is to be accessed or altered. This ID
                is the value between the "/d/" and the "/edit" in the URL of your GoogleAnalytics.
              </p>
            }
          >
            webPropertyId
          </InlineFormLabel>
          <SegmentAsync
            loadOptions={() => datasource.getWebPropertyIds(accountId)}
            placeholder="Enter webPropertyId"
            value={webPropertyId}
            allowCustomValue={true}
            onChange={this.onWebPropertyChange}
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
                The <code>viewId</code> is used to identify which GoogleAnalytics is to be accessed or altered. This ID
                is the value between the "/d/" and the "/edit" in the URL of your GoogleAnalytics.
              </p>
            }
          >
            viewId
          </InlineFormLabel>
          <SegmentAsync
            loadOptions={() => datasource.getProfileIds(accountId, webPropertyId)}
            placeholder="Enter viewId"
            value={viewId}
            allowCustomValue={true}
            onChange={this.onViewIDChange}
          ></SegmentAsync>
          {viewId && <LinkButton style={{ marginTop: 1 }} variant="link" icon="link" target="_blank"></LinkButton>}
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
      </>
    );
  }
}
