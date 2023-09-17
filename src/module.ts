import { DataSourcePlugin } from '@grafana/data';
import { QueryEditorCommon } from 'QueryEditorCommon';
import { ConfigEditor } from './ConfigEditor';
import { DataSource } from './DataSource';
import { GADataSourceOptions, GAQuery } from './types';

export const plugin = new DataSourcePlugin<DataSource, GAQuery, GADataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditorCommon);
