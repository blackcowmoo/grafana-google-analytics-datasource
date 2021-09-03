import { DataQuery, DataSourceJsonData, SelectableValue } from '@grafana/data';

export interface GAQuery extends DataQuery {
  accountId: string;
  webPropertyId: string;
  profileId: string;
  startDate: string;
  endDate: string;
  metrics: string[];
  timeDimension: string;
  dimensions: string[];
  selectedMetrics: Array<SelectableValue<string>>;
  selectedTimeDimensions: SelectableValue<string>;
  selectedDimensions: Array<SelectableValue<string>>;
  cacheDurationSeconds?: number;
  timezone: string;
  filtersExpression: string;
}
// mapping on google-key.json
export interface JWT {
  private_key: any;
  token_uri: any;
  client_email: any;
  project_id: any;
}

export const defaultQuery: Partial<GAQuery> = {
  // constant: 6.5,
};

/**
 * These are options configured for each DataSource instance
 */
export interface GADataSourceOptions extends DataSourceJsonData {}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface GASecureJsonData {
  jwt?: string;
  profileId?: string;
  timezone?: string;
}

export interface GAMetadata {
  id: string;
  kind: string;
  attributes: GAMetadataAttribute;
}

export interface GAMetadataAttribute {
  type: string;
  dataType: string;
  group: string;
  status?: string;
  uiName: string;
  description: string;
  allowedInSegments?: string;
  addedInAPIVersion?: string;
}
