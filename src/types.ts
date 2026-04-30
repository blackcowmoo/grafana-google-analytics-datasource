import { DataQuery, SelectableValue } from '@grafana/data';
import type {
  DataSourceOptions as GoogleDataSourceOptions,
  DataSourceSecureJsonData as GoogleDataSourceSecureJsonData,
} from '@grafana/google-sdk';

export interface GAQuery extends DataQuery {
  displayName: Map<string, string>
  version: string;
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
  mode: string;
  dimensionFilter: GAFilterExpression;
  serviceLevel: string;
  // metricFilter: Array<GAFilter>;
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
 * These are options configured for each DataSource instance.
 * Extends @grafana/google-sdk's DataSourceOptions so the shared
 * <ConnectionConfig /> component can read/write the same fields.
 */
export interface GADataSourceOptions extends GoogleDataSourceOptions {
  // Optional so the type stays structurally compatible with the SDK's
  // DataSourceOptions (covariance through <ConnectionConfig />). The
  // ConfigEditor fills in the default 'v4' on first mount, and the
  // DataSource constructor coerces missing values to 'v4'.
  version?: string;
}

/**
 * Secret values stored on the backend. Extends the SDK's
 * DataSourceSecureJsonData (which carries `privateKey`) with the legacy
 * `jwt` blob so existing datasources continue to authenticate via the
 * backend's dual-read fallback.
 */
export interface GASecureJsonData extends GoogleDataSourceSecureJsonData {
  jwt?: string;
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

export interface AccountSummary {
  Account: string
  DisplayName: string
  PropertySummaries: PropertySummary[]
}

export interface PropertySummary {
  Property: string
  DisplayName: string
  Parent: string
  ProfileSummaries: ProfileSummary[]
}

export interface ProfileSummary {
  Profile: string
  DisplayName: string
  Parent: string
  Type: string
}

// https://developers.google.com/analytics/devguides/reporting/data/v1/rest/v1beta/FilterExpression?hl=ko#Filter

export interface GAFilterExpression {
  andGroup?: GAFilterExpressionList
  orGroup?: GAFilterExpressionList
  notExpression?: GAFilterExpression
  filter?: GAFilter
}

export interface GAFilterExpressionList {
  expressions: GAFilterExpression[]
}

export interface GAFilter {
  fieldName: string;
  
  // filterType: GAMetricFilterType | GADimensionFilterType | undefined;
  filterType: GADimensionFilterType;
  stringFilter?: GAStringFilter;
  inListFilter?: GAInListFilter;
  numberFilter?: GANumbericFilter;
  betweenFilter?: GABetweenFilter;
}

export interface GAStringFilter {
  matchType: GAStringFilterMatchType;
  value: string;
  caseSensitive: boolean;
}

export interface GAInListFilter {
  values: string[];
  caseSensitive: boolean;
}

export interface GANumbericFilter {
  operation: GANumbericFilterOperation;
  value: GANumbericValue;
  caseSensitive: boolean;
}

export interface GABetweenFilter {
  fromValue: GANumbericValue;
  toValue: GANumbericValue;
}

export interface GANumbericValue {
  int64Value: string;
  doubleValue: number;
}

export enum GADimensionFilterType {
  STRING = "STRING",
  IN_LIST = "IN_LIST"
}

export enum GAMetricFilterType {
  NUMBERIC,
  BETWEEN
}


export enum GAStringFilterMatchType {
  MATCH_TYPE_UNSPECIFIED = "MATCH_TYPE_UNSPECIFIED",
  EXACT = "EXACT",
  BEGINS_WITH = "BEGINS_WITH",
  ENDS_WITH = "ENDS_WITH",
  CONTAINS = "CONTAINS",
  FULL_REGEXP = "FULL_REGEXP",
  PARTIAL_REGEXP = "PARTIAL_REGEXP"
}

export enum GANumbericFilterOperation {
  OPERATION_UNSPECIFIED,
  EQUAL,
  LESS_THAN,
  LESS_THAN_OR_EQUAL,
  GREATER_THAN,
  GREATER_THAN_OR_EQUAL,
}
