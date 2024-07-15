import { QueryEditorProps } from '@grafana/data';
import { DataSource } from 'DataSource';
import { QueryEditorGA4 } from 'QueryEditorGA4';
import { QueryEditorUA } from 'QueryEditorUA';
import React, { PureComponent } from 'react';
import { GADataSourceOptions, GAQuery } from 'types';

type Props = QueryEditorProps<DataSource, GAQuery, GADataSourceOptions>;


export class QueryEditorCommon extends PureComponent<Props> {
  constructor(props: Readonly<Props>) {
    super(props);
    this.props.query.version = props.datasource.getGaVersion()
  }
  render() {
    const { query, datasource, onChange, onRunQuery } = this.props;
    const { version } = query
    if (version === "v4") {
      return <QueryEditorGA4 datasource={datasource} onChange={onChange} onRunQuery={onRunQuery} query={query}></QueryEditorGA4>
    } else {
      return <QueryEditorUA datasource={datasource} onChange={onChange} onRunQuery={onRunQuery} query={query}></QueryEditorUA>
    }
  }
}
