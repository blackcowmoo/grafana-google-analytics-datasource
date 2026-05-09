import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'DataSource';
import { QueryEditorGA4 } from 'QueryEditorGA4';
import React, { PureComponent } from 'react';
import { GADataSourceOptions, GAQuery } from 'types';

type Props = QueryEditorProps<DataSource, GAQuery, GADataSourceOptions>;

export class QueryEditorCommon extends PureComponent<Props> {
  render() {
    const { query, datasource, onChange, onRunQuery } = this.props;
    return <QueryEditorGA4 datasource={datasource} onChange={onChange} onRunQuery={onRunQuery} query={query}></QueryEditorGA4>;
  }
}
