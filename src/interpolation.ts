import { ScopedVars } from '@grafana/data';
import { TemplateSrv } from '@grafana/runtime';
import { GAFilter, GAFilterExpression } from './types';

// Expand a value that may reference a template variable.
// Multi-value variables return one element per selected value; literal or
// single-value inputs return a one-element array.
export function expandVariableToArray(
  templateSrv: TemplateSrv,
  value: string,
  scopedVars: ScopedVars
): string[] {
  const out: string[] = [];
  let multi = false;
  const replaced = templateSrv.replace(value, scopedVars, (v: unknown) => {
    if (Array.isArray(v)) {
      multi = true;
      for (const item of v) {
        out.push(String(item));
      }
      return v.join(',');
    }
    return v === undefined || v === null ? '' : String(v);
  });
  return multi ? out : [replaced];
}

export function interpolateFilter(
  templateSrv: TemplateSrv,
  filter: GAFilter,
  scopedVars: ScopedVars
): void {
  if (filter.stringFilter) {
    filter.stringFilter.value = templateSrv.replace(filter.stringFilter.value, scopedVars);
  }
  if (filter.inListFilter) {
    filter.inListFilter.values = filter.inListFilter.values.flatMap((value) =>
      expandVariableToArray(templateSrv, value, scopedVars)
    );
  }
}

export function interpolateFilterExpression(
  templateSrv: TemplateSrv,
  expression: GAFilterExpression | undefined,
  scopedVars: ScopedVars
): void {
  if (!expression) {
    return;
  }
  if (expression.filter) {
    interpolateFilter(templateSrv, expression.filter, scopedVars);
  }
  if (expression.andGroup) {
    for (const e of expression.andGroup.expressions) {
      interpolateFilterExpression(templateSrv, e, scopedVars);
    }
  }
  if (expression.orGroup) {
    for (const e of expression.orGroup.expressions) {
      interpolateFilterExpression(templateSrv, e, scopedVars);
    }
  }
  if (expression.notExpression) {
    interpolateFilterExpression(templateSrv, expression.notExpression, scopedVars);
  }
}
