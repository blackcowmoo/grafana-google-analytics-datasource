import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface GAQuery extends DataQuery {
  accountId: string;
  webPropertyId: string;
  profileId: string;
  startDate: string;
  endDate: string;
  metric: string;
  dimension: string;
  cacheDurationSeconds?: number;
}
// mapping on google-key.json
export interface JWT {
  private_key: any;
  token_uri: any;
  client_email: any;
  project_id: any;
}

export enum GoogleAuthType {
  JWT = 'jwt',
  KEY = 'key',
}

export const googleAuthTypes = [
  { label: 'API Key', value: GoogleAuthType.KEY },
  { label: 'Google JWT File', value: GoogleAuthType.JWT },
];

export const defaultQuery: Partial<GAQuery> = {
  // constant: 6.5,
};

/**
 * These are options configured for each DataSource instance
 */
export interface GADataSourceOptions extends DataSourceJsonData {
  authType: GoogleAuthType;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface GASecureJsonData {
  apiKey?: string;
  jwt?: string;
  profileId?: string;
}
