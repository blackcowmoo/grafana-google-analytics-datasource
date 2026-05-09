import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSource } from './DataSource';
import { GADataSourceOptions, GADimensionFilterType, GAFilterExpression, GAQuery } from './types';

jest.mock('@grafana/runtime', () => {
  // Inline the TemplateSrv stub here because jest.mock factories cannot
  // reference outer-scope variables. The tests configure variable values via
  // the exported __setVars helper before calling applyTemplateVariables.
  const vars: Record<string, string | string[]> = {};
  const replace = (text?: string, _scoped?: any, format?: string | Function): string => {
    if (!text) {
      return '';
    }
    return text.replace(/\$\{?(\w+)\}?/g, (_, name: string) => {
      const raw = vars[name];
      if (raw === undefined) {
        return '';
      }
      if (typeof format === 'function') {
        return (format as Function)(raw);
      }
      return Array.isArray(raw) ? raw.join(',') : String(raw);
    });
  };
  return {
    DataSourceWithBackend: class {
      constructor(_settings: any) {}
    },
    getTemplateSrv: () => ({ replace, getVariables: () => [], containsTemplate: () => false }),
    __setVars: (next: Record<string, string | string[]>) => {
      for (const k of Object.keys(vars)) {
        delete vars[k];
      }
      Object.assign(vars, next);
    },
  };
});

const setVars = (next: Record<string, string | string[]>) => {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  (require('@grafana/runtime') as any).__setVars(next);
};

const makeDataSource = (): DataSource => {
  const settings = {
    jsonData: { version: 'v4' },
  } as unknown as DataSourceInstanceSettings<GADataSourceOptions>;
  return new DataSource(settings);
};

const makeQuery = (overrides: Partial<GAQuery> = {}): GAQuery => {
  return {
    refId: 'A',
    accountId: '',
    webPropertyId: 'properties/12345',
    profileId: '',
    startDate: '',
    endDate: '',
    metrics: [],
    timeDimension: '',
    dimensions: [],
    selectedMetrics: [],
    selectedTimeDimensions: { value: '' },
    selectedDimensions: [],
    timezone: 'UTC',
    filtersExpression: '',
    mode: 'time series',
    serviceLevel: '',
    displayName: new Map<string, string>(),
    version: 'v4',
    ...overrides,
  } as GAQuery;
};

describe('DataSource.applyTemplateVariables', () => {
  it('does not mutate the input dimensionFilter (deep-clones before interpolation)', () => {
    setVars({ campaigns: ['a', 'b'] });
    const ds = makeDataSource();
    const inputFilter: GAFilterExpression = {
      orGroup: {
        expressions: [
          {
            filter: {
              fieldName: 'campaignName',
              filterType: GADimensionFilterType.IN_LIST,
              inListFilter: { values: ['$campaigns'], caseSensitive: true },
            },
          },
        ],
      },
    };
    const inputSnapshot = JSON.parse(JSON.stringify(inputFilter));
    const query = makeQuery({ dimensionFilter: inputFilter });

    const interpolated = ds.applyTemplateVariables(query, {});

    // Input must remain untouched — dashboards reuse the same object across
    // re-renders, mutation would silently corrupt the panel state (#148).
    expect(inputFilter).toEqual(inputSnapshot);
    // Output must reflect the expansion.
    expect(interpolated.dimensionFilter.orGroup.expressions[0].filter.inListFilter.values).toEqual(
      ['a', 'b']
    );
  });

  it('interpolates webPropertyId from a single-value variable', () => {
    setVars({ prop: '987' });
    const ds = makeDataSource();
    const query = makeQuery({ webPropertyId: 'properties/$prop' });

    const interpolated = ds.applyTemplateVariables(query, {});

    expect(interpolated.webPropertyId).toBe('properties/987');
  });

  it('passes dimensionFilter through unchanged when no variables are referenced', () => {
    setVars({});
    const ds = makeDataSource();
    const inputFilter: GAFilterExpression = {
      filter: {
        fieldName: 'campaignName',
        filterType: GADimensionFilterType.IN_LIST,
        inListFilter: { values: ['static-only'], caseSensitive: true },
      },
    };
    const query = makeQuery({ dimensionFilter: inputFilter });

    const interpolated = ds.applyTemplateVariables(query, {});

    expect(interpolated.dimensionFilter.filter.inListFilter.values).toEqual(['static-only']);
    // And input is still not the same reference (deep-clone happened either
    // way), guarding against future regressions to in-place mutation.
    expect(interpolated.dimensionFilter).not.toBe(inputFilter);
  });

  it('handles a missing dimensionFilter gracefully', () => {
    setVars({});
    const ds = makeDataSource();
    const query = makeQuery({ dimensionFilter: undefined as unknown as GAFilterExpression });

    expect(() => ds.applyTemplateVariables(query, {})).not.toThrow();
  });
});
