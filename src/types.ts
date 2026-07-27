import { DataQuery, SelectableValue } from '@grafana/data';
import type {
  DataSourceOptions as GoogleDataSourceOptions,
  DataSourceSecureJsonData as GoogleDataSourceSecureJsonData,
} from '@grafana/google-sdk';

export interface GAQuery extends DataQuery {
  displayName: Map<string, string>
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
  metricFilter?: GAFilterExpression;
  serviceLevel: string;
}

// mapping on google-key.json
export interface JWT {
  private_key: any;
  token_uri: any;
  client_email: any;
  project_id: any;
}

export const defaultQuery: Partial<GAQuery> = {};

/**
 * These are options configured for each DataSource instance.
 * Extends @grafana/google-sdk's DataSourceOptions so the shared
 * <ConnectionConfig /> component can read/write the same fields.
 */
export interface GADataSourceOptions extends GoogleDataSourceOptions {}

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

// https://developers.google.com/analytics/devguides/reporting/data/v1/rest/v1beta/FilterExpression

export interface GAFilterExpression {
  andGroup?: GAFilterExpressionList;
  orGroup?: GAFilterExpressionList;
  notExpression?: GAFilterExpression;
  filter?: GAFilter;
}

export interface GAFilterExpressionList {
  expressions: GAFilterExpression[];
}

export interface GAFilter {
  fieldName: string;
  filterType: GADimensionFilterType;
  stringFilter?: GAStringFilter;
  inListFilter?: GAInListFilter;
  numericFilter?: GANumericFilter;
  betweenFilter?: GABetweenFilter;
  emptyFilter?: GAEmptyFilter;
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

export interface GANumericFilter {
  operation: GANumericFilterOperation;
  value: GANumericValue;
}

export interface GABetweenFilter {
  fromValue: GANumericValue;
  toValue: GANumericValue;
}

// int64Value serialises as a JSON string per the GA4 API spec
export interface GANumericValue {
  int64Value?: string;
  doubleValue?: number;
}

export interface GAEmptyFilter {}

export enum GADimensionFilterType {
  STRING  = 'STRING',
  IN_LIST = 'IN_LIST',
  NUMERIC = 'NUMERIC',
  BETWEEN = 'BETWEEN',
  EMPTY   = 'EMPTY',
}

export enum GAStringFilterMatchType {
  EXACT           = 'EXACT',
  BEGINS_WITH     = 'BEGINS_WITH',
  ENDS_WITH       = 'ENDS_WITH',
  CONTAINS        = 'CONTAINS',
  FULL_REGEXP     = 'FULL_REGEXP',
  PARTIAL_REGEXP  = 'PARTIAL_REGEXP',
}

export enum GANumericFilterOperation {
  OPERATION_UNSPECIFIED   = 'OPERATION_UNSPECIFIED',
  EQUAL                   = 'EQUAL',
  LESS_THAN               = 'LESS_THAN',
  LESS_THAN_OR_EQUAL      = 'LESS_THAN_OR_EQUAL',
  GREATER_THAN            = 'GREATER_THAN',
  GREATER_THAN_OR_EQUAL   = 'GREATER_THAN_OR_EQUAL',
}

// Legacy aliases kept for any persisted dashboard JSON that might reference them
/** @deprecated use GANumericFilter */
export type GANumbericFilter = GANumericFilter;
/** @deprecated use GANumericFilterOperation */
export type GANumbericFilterOperation = GANumericFilterOperation;
/** @deprecated use GANumericValue */
export type GANumbericValue = GANumericValue;
