import { ScopedVars } from '@grafana/data';
import { TemplateSrv } from '@grafana/runtime';
import { expandVariableToArray, interpolateFilterExpression } from './interpolation';
import { GADimensionFilterType, GAFilterExpression, GAStringFilterMatchType } from './types';

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
});
