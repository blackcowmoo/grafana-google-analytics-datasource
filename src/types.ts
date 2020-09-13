import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface GAQuery extends DataQuery {
  queryText?: string;
  constant: number;
}

export const defaultQuery: Partial<GAQuery> = {
  constant: 6.5,
};

/**
 * These are options configured for each DataSource instance
 */
export interface GADataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface GASecureJsonData {
  apiKey?: string;
}
