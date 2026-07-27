import { ScopedVars } from '@grafana/data';
import { TemplateSrv } from '@grafana/runtime';
import { expandVariableToArray, interpolateFilterExpression } from './interpolation';
import { GADimensionFilterType, GAFilterExpression, GANumericFilterOperation, GAStringFilterMatchType } from './types';

// Minimal TemplateSrv stub that mimics Grafana's behavior for unit tests:
// - $var / ${var} is looked up in the map.
// - If the variable value is an array and a format function is supplied,
//   the function is invoked with the array, and its return value replaces
//   the variable reference in the output string.
// - If no format function is supplied, arrays are joined with a comma.
function makeTemplateSrv(vars: Record<string, string | string[]>): TemplateSrv {
  const replace = (
    text?: string,
    _scopedVars?: ScopedVars,
    format?: string | Function
  ): string => {
    if (!text) {
      return '';
    }
    return text.replace(/\$\{?(\w+)\}?/g, (_, name: string) => {
      const raw = vars[name];
      if (raw === undefined) {
        return '';
      }
      if (typeof format === 'function') {
        return format(raw);
      }
      return Array.isArray(raw) ? raw.join(',') : String(raw);
    });
  };
  return {
    replace,
    getVariables: () => [],
    containsTemplate: (t?: string) => !!t && /\$/.test(t),
    updateTimeRange: () => undefined,
  } as unknown as TemplateSrv;
}

describe('expandVariableToArray', () => {
  it('returns literal value as single-element array when no variable is referenced', () => {
    const srv = makeTemplateSrv({});
    expect(expandVariableToArray(srv, 'foo', {})).toEqual(['foo']);
  });

  it('expands a multi-value variable into one array element per value', () => {
    const srv = makeTemplateSrv({ myVar: ['xxxx', 'yyyy'] });
    expect(expandVariableToArray(srv, '$myVar', {})).toEqual(['xxxx', 'yyyy']);
  });

  it('returns a single-element array for a single-value variable', () => {
    const srv = makeTemplateSrv({ myVar: 'xxxx' });
    expect(expandVariableToArray(srv, '$myVar', {})).toEqual(['xxxx']);
  });

  it('handles undefined variable by returning empty string replacement', () => {
    const srv = makeTemplateSrv({});
    expect(expandVariableToArray(srv, '$missing', {})).toEqual(['']);
  });
});

describe('interpolateFilterExpression', () => {
  it('splits inListFilter values for multi-value variable (fixes #97)', () => {
    const srv = makeTemplateSrv({ campaigns: ['xxxx', 'yyyy'] });
    const expr: GAFilterExpression = {
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
    interpolateFilterExpression(srv, expr, {});
    expect(expr.orGroup!.expressions[0].filter!.inListFilter!.values).toEqual(['xxxx', 'yyyy']);
  });

  it('interpolates stringFilter value', () => {
    const srv = makeTemplateSrv({ name: 'foo' });
    const expr: GAFilterExpression = {
      filter: {
        fieldName: 'eventName',
        filterType: GADimensionFilterType.STRING,
        stringFilter: {
          matchType: GAStringFilterMatchType.EXACT,
          value: '$name',
          caseSensitive: false,
        },
      },
    };
    interpolateFilterExpression(srv, expr, {});
    expect(expr.filter!.stringFilter!.value).toBe('foo');
  });

  it('recurses into andGroup, orGroup, and notExpression', () => {
    const srv = makeTemplateSrv({ v: ['a', 'b'] });
    const expr: GAFilterExpression = {
      andGroup: {
        expressions: [
          {
            orGroup: {
              expressions: [
                {
                  notExpression: {
                    filter: {
                      fieldName: 'x',
                      filterType: GADimensionFilterType.IN_LIST,
                      inListFilter: { values: ['$v'], caseSensitive: true },
                    },
                  },
                },
              ],
            },
          },
        ],
      },
    };
    interpolateFilterExpression(srv, expr, {});
    const leaf = expr.andGroup!.expressions[0].orGroup!.expressions[0].notExpression!.filter!;
    expect(leaf.inListFilter!.values).toEqual(['a', 'b']);
  });

  it('noop on undefined expression', () => {
    const srv = makeTemplateSrv({});
    expect(() => interpolateFilterExpression(srv, undefined, {})).not.toThrow();
  });

  it('flattens inListFilter with mix of literal and multi-value variable', () => {
    const srv = makeTemplateSrv({ campaigns: ['xxxx', 'yyyy'] });
    const expr: GAFilterExpression = {
      filter: {
        fieldName: 'campaignName',
        filterType: GADimensionFilterType.IN_LIST,
        inListFilter: { values: ['static', '$campaigns', 'tail'], caseSensitive: true },
      },
    };
    interpolateFilterExpression(srv, expr, {});
    expect(expr.filter!.inListFilter!.values).toEqual(['static', 'xxxx', 'yyyy', 'tail']);
  });

  it('keeps empty inListFilter values intact', () => {
    const srv = makeTemplateSrv({ x: 'foo' });
    const expr: GAFilterExpression = {
      filter: {
        fieldName: 'campaignName',
        filterType: GADimensionFilterType.IN_LIST,
        inListFilter: { values: [], caseSensitive: true },
      },
    };
    interpolateFilterExpression(srv, expr, {});
    expect(expr.filter!.inListFilter!.values).toEqual([]);
  });

  it('does not crash on filter with neither stringFilter nor inListFilter', () => {
    const srv = makeTemplateSrv({});
    const expr: GAFilterExpression = {
      filter: {
        fieldName: 'eventName',
        filterType: GADimensionFilterType.STRING,
      },
    };
    expect(() => interpolateFilterExpression(srv, expr, {})).not.toThrow();
  });

  it('expands multi-value variable inside every nested group independently', () => {
    const srv = makeTemplateSrv({ a: ['1', '2'], b: ['3', '4'] });
    const expr: GAFilterExpression = {
      andGroup: {
        expressions: [
          {
            filter: {
              fieldName: 'x',
              filterType: GADimensionFilterType.IN_LIST,
              inListFilter: { values: ['$a'], caseSensitive: true },
            },
          },
          {
            filter: {
              fieldName: 'y',
              filterType: GADimensionFilterType.IN_LIST,
              inListFilter: { values: ['$b'], caseSensitive: true },
            },
          },
        ],
      },
    };
    interpolateFilterExpression(srv, expr, {});
    expect(expr.andGroup!.expressions[0].filter!.inListFilter!.values).toEqual(['1', '2']);
    expect(expr.andGroup!.expressions[1].filter!.inListFilter!.values).toEqual(['3', '4']);
  });

  it('does not touch numericFilter (no string interpolation needed)', () => {
    const srv = makeTemplateSrv({ n: '999' });
    const expr: GAFilterExpression = {
      filter: {
        fieldName: 'sessions',
        filterType: GADimensionFilterType.NUMERIC,
        numericFilter: {
          operation: GANumericFilterOperation.GREATER_THAN,
          value: { int64Value: '100' },
        },
      },
    };
    interpolateFilterExpression(srv, expr, {});
    // Numeric values are set via the UI as literal numbers, not interpolated.
    expect(expr.filter!.numericFilter!.value.int64Value).toBe('100');
    expect(expr.filter!.numericFilter!.operation).toBe(GANumericFilterOperation.GREATER_THAN);
  });

  it('does not crash on betweenFilter', () => {
    const srv = makeTemplateSrv({});
    const expr: GAFilterExpression = {
      filter: {
        fieldName: 'sessions',
        filterType: GADimensionFilterType.BETWEEN,
        betweenFilter: { fromValue: { int64Value: '10' }, toValue: { int64Value: '100' } },
      },
    };
    expect(() => interpolateFilterExpression(srv, expr, {})).not.toThrow();
    expect(expr.filter!.betweenFilter!.fromValue.int64Value).toBe('10');
    expect(expr.filter!.betweenFilter!.toValue.int64Value).toBe('100');
  });

  it('does not crash on emptyFilter', () => {
    const srv = makeTemplateSrv({});
    const expr: GAFilterExpression = {
      filter: {
        fieldName: 'campaignName',
        filterType: GADimensionFilterType.EMPTY,
        emptyFilter: {},
      },
    };
    expect(() => interpolateFilterExpression(srv, expr, {})).not.toThrow();
  });
});
