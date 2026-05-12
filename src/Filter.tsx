import { GrafanaTheme2, SelectableValue } from '@grafana/data';
import { css } from '@emotion/css';
import {
  AsyncSelect,
  IconButton,
  InlineSwitch,
  Input,
  Select,
  TagsInput,
  useStyles2,
} from '@grafana/ui';
import React from 'react';
import {
  GABetweenFilter,
  GADimensionFilterType,
  GAFilter,
  GAFilterExpression,
  GAFilterExpressionList,
  GAInListFilter,
  GANumericFilter,
  GANumericFilterOperation,
  GANumericValue,
  GAStringFilter,
  GAStringFilterMatchType,
} from 'types';

export type LoadFieldsFn = (query: string) => Promise<Array<SelectableValue<string>>>;

export interface GAFilterExpressionComponentProps {
  expression: GAFilterExpression;
  onChange: (expr: GAFilterExpression) => void;
  onDelete?: () => void;
  loadFields: LoadFieldsFn;
  /** Visual nesting depth — drives indentation and border colour */
  depth?: number;
}

// ─── option arrays ────────────────────────────────────────────────────────────

const EXPR_TYPE_OPTIONS: Array<SelectableValue<string>> = [
  { label: 'No filter',     value: 'none'          },
  { label: 'Filter',        value: 'filter'        },
  { label: 'AND group',     value: 'andGroup'      },
  { label: 'OR group',      value: 'orGroup'       },
  { label: 'NOT',           value: 'notExpression' },
];

const FILTER_TYPE_OPTIONS: Array<SelectableValue<GADimensionFilterType>> = [
  { label: 'String',  value: GADimensionFilterType.STRING  },
  { label: 'In list', value: GADimensionFilterType.IN_LIST },
  { label: 'Numeric', value: GADimensionFilterType.NUMERIC },
  { label: 'Between', value: GADimensionFilterType.BETWEEN },
  { label: 'Empty',   value: GADimensionFilterType.EMPTY   },
];

const STRING_MATCH_OPTIONS: Array<SelectableValue<GAStringFilterMatchType>> = [
  { label: 'Exact',          value: GAStringFilterMatchType.EXACT          },
  { label: 'Begins with',    value: GAStringFilterMatchType.BEGINS_WITH    },
  { label: 'Ends with',      value: GAStringFilterMatchType.ENDS_WITH      },
  { label: 'Contains',       value: GAStringFilterMatchType.CONTAINS       },
  { label: 'Full regexp',    value: GAStringFilterMatchType.FULL_REGEXP    },
  { label: 'Partial regexp', value: GAStringFilterMatchType.PARTIAL_REGEXP },
];

const NUMERIC_OP_OPTIONS: Array<SelectableValue<GANumericFilterOperation>> = [
  { label: '=',  value: GANumericFilterOperation.EQUAL                 },
  { label: '<',  value: GANumericFilterOperation.LESS_THAN             },
  { label: '≤',  value: GANumericFilterOperation.LESS_THAN_OR_EQUAL   },
  { label: '>',  value: GANumericFilterOperation.GREATER_THAN          },
  { label: '≥',  value: GANumericFilterOperation.GREATER_THAN_OR_EQUAL },
];

// ─── helpers ──────────────────────────────────────────────────────────────────

function exprType(expr: GAFilterExpression): string {
  if (expr.andGroup)                    { return 'andGroup'; }
  if (expr.orGroup)                     { return 'orGroup'; }
  if (expr.notExpression !== undefined) { return 'notExpression'; }
  if (expr.filter)                      { return 'filter'; }
  return 'none';
}

function makeDefaultFilter(): GAFilter {
  return {
    fieldName: '',
    filterType: GADimensionFilterType.STRING,
    stringFilter: { matchType: GAStringFilterMatchType.EXACT, value: '', caseSensitive: false },
  };
}

function makeDefaultExpr(): GAFilterExpression {
  return { filter: makeDefaultFilter() };
}

function numericValueStr(v?: GANumericValue): string {
  if (!v) { return ''; }
  // int64Value is a string like "0", "123" — non-empty means it was explicitly set
  if (v.int64Value != null && v.int64Value !== '') { return v.int64Value; }
  if (v.doubleValue != null && v.doubleValue !== 0) { return String(v.doubleValue); }
  return '';
}

function strToNumericValue(s: string): GANumericValue {
  const trimmed = s.trim();
  if (trimmed === '') { return {}; }
  const n = Number(trimmed);
  if (Number.isInteger(n)) { return { int64Value: trimmed }; }
  return { doubleValue: isNaN(n) ? 0 : n };
}

// ─── styles ───────────────────────────────────────────────────────────────────

const getStyles = (theme: GrafanaTheme2) => {
  const borderColors = [
    theme.colors.primary.main,
    theme.colors.warning.main,
    theme.colors.success.main,
  ];

  return {
    row: css`
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: ${theme.spacing(0.5)};
      min-height: ${theme.spacing(4)};
    `,
    group: (depth: number) => css`
      border-left: 3px solid ${borderColors[depth % borderColors.length]};
      padding-left: ${theme.spacing(1.5)};
      margin-top: ${theme.spacing(0.5)};
      display: flex;
      flex-direction: column;
      gap: ${theme.spacing(0.5)};
    `,
    groupHeader: css`
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: ${theme.spacing(0.5)};
    `,
    groupLabel: css`
      color: ${theme.colors.text.secondary};
      font-size: ${theme.typography.bodySmall.fontSize};
      padding: 0 ${theme.spacing(0.5)};
      align-self: center;
    `,
    emptyLabel: css`
      color: ${theme.colors.text.disabled};
      font-style: italic;
      font-size: ${theme.typography.bodySmall.fontSize};
      align-self: center;
      padding: 0 ${theme.spacing(0.5)};
    `,
    betweenSep: css`
      color: ${theme.colors.text.secondary};
      align-self: center;
      padding: 0 ${theme.spacing(0.25)};
    `,
    tagInput: css`
      flex: 1;
      min-width: 160px;
    `,
    caseLabel: css`
      color: ${theme.colors.text.secondary};
      font-size: ${theme.typography.bodySmall.fontSize};
      align-self: center;
      white-space: nowrap;
    `,
  };
};

// ─── leaf filter editors ──────────────────────────────────────────────────────

interface StringEditorProps {
  filter: GAStringFilter;
  onChange: (f: GAStringFilter) => void;
  styles: ReturnType<typeof getStyles>;
}
const StringEditor: React.FC<StringEditorProps> = ({ filter, onChange, styles }) => (
  <>
    <Select<GAStringFilterMatchType>
      options={STRING_MATCH_OPTIONS}
      value={filter.matchType}
      onChange={(o) => onChange({ ...filter, matchType: o.value! })}
      width={16}
      menuPlacement="bottom"
    />
    <Input
      value={filter.value}
      onChange={(e) => onChange({ ...filter, value: e.currentTarget.value })}
      placeholder="value or $variable"
      width={20}
    />
    <span className={styles.caseLabel}>Aa</span>
    <InlineSwitch
      value={filter.caseSensitive}
      onChange={(e) => onChange({ ...filter, caseSensitive: e.currentTarget.checked })}
      showLabel={false}
      title="Case sensitive"
    />
  </>
);

interface InListEditorProps {
  filter: GAInListFilter;
  onChange: (f: GAInListFilter) => void;
  styles: ReturnType<typeof getStyles>;
}
const InListEditor: React.FC<InListEditorProps> = ({ filter, onChange, styles }) => (
  <>
    <div className={styles.tagInput}>
      <TagsInput
        tags={filter.values}
        onChange={(tags) => onChange({ ...filter, values: tags })}
        placeholder="Add value or $variable, press Enter"
        addOnBlur
      />
    </div>
    <span className={styles.caseLabel}>Aa</span>
    <InlineSwitch
      value={filter.caseSensitive}
      onChange={(e) => onChange({ ...filter, caseSensitive: e.currentTarget.checked })}
      showLabel={false}
      title="Case sensitive"
    />
  </>
);

interface NumericEditorProps {
  filter: GANumericFilter;
  onChange: (f: GANumericFilter) => void;
}
const NumericEditor: React.FC<NumericEditorProps> = ({ filter, onChange }) => (
  <>
    <Select<GANumericFilterOperation>
      options={NUMERIC_OP_OPTIONS}
      value={filter.operation ?? GANumericFilterOperation.EQUAL}
      onChange={(o) => onChange({ ...filter, operation: o.value! })}
      width={8}
      menuPlacement="bottom"
    />
    <Input
      value={numericValueStr(filter.value)}
      onChange={(e) => onChange({ ...filter, value: strToNumericValue(e.currentTarget.value) })}
      placeholder="number"
      width={14}
      type="number"
    />
  </>
);

interface BetweenEditorProps {
  filter: GABetweenFilter;
  onChange: (f: GABetweenFilter) => void;
  styles: ReturnType<typeof getStyles>;
}
const BetweenEditor: React.FC<BetweenEditorProps> = ({ filter, onChange, styles }) => (
  <>
    <Input
      value={numericValueStr(filter.fromValue)}
      onChange={(e) => onChange({ ...filter, fromValue: strToNumericValue(e.currentTarget.value) })}
      placeholder="from"
      width={10}
      type="number"
    />
    <span className={styles.betweenSep}>–</span>
    <Input
      value={numericValueStr(filter.toValue)}
      onChange={(e) => onChange({ ...filter, toValue: strToNumericValue(e.currentTarget.value) })}
      placeholder="to"
      width={10}
      type="number"
    />
  </>
);

// ─── single leaf filter row ───────────────────────────────────────────────────

interface LeafFilterProps {
  filter: GAFilter;
  onChange: (f: GAFilter) => void;
  loadFields: LoadFieldsFn;
  styles: ReturnType<typeof getStyles>;
}

const LeafFilterEditor: React.FC<LeafFilterProps> = ({ filter, onChange, loadFields, styles }) => {
  const handleTypeChange = (type: GADimensionFilterType) => {
    const base: GAFilter = { fieldName: filter.fieldName, filterType: type };
    switch (type) {
      case GADimensionFilterType.STRING:
        return onChange({ ...base, stringFilter: { matchType: GAStringFilterMatchType.EXACT, value: '', caseSensitive: false } });
      case GADimensionFilterType.IN_LIST:
        return onChange({ ...base, inListFilter: { values: [], caseSensitive: false } });
      case GADimensionFilterType.NUMERIC:
        return onChange({ ...base, numericFilter: { operation: GANumericFilterOperation.EQUAL, value: {} } });
      case GADimensionFilterType.BETWEEN:
        return onChange({ ...base, betweenFilter: { fromValue: {}, toValue: {} } });
      case GADimensionFilterType.EMPTY:
        return onChange({ ...base, emptyFilter: {} });
    }
  };

  const renderParams = () => {
    switch (filter.filterType) {
      case GADimensionFilterType.STRING:
        return (
          <StringEditor
            filter={filter.stringFilter ?? { matchType: GAStringFilterMatchType.EXACT, value: '', caseSensitive: false }}
            onChange={(f) => onChange({ ...filter, stringFilter: f })}
            styles={styles}
          />
        );
      case GADimensionFilterType.IN_LIST:
        return (
          <InListEditor
            filter={filter.inListFilter ?? { values: [], caseSensitive: false }}
            onChange={(f) => onChange({ ...filter, inListFilter: f })}
            styles={styles}
          />
        );
      case GADimensionFilterType.NUMERIC:
        return (
          <NumericEditor
            filter={filter.numericFilter ?? { operation: GANumericFilterOperation.EQUAL, value: {} }}
            onChange={(f) => onChange({ ...filter, numericFilter: f })}
          />
        );
      case GADimensionFilterType.BETWEEN:
        return (
          <BetweenEditor
            filter={filter.betweenFilter ?? { fromValue: {}, toValue: {} }}
            onChange={(f) => onChange({ ...filter, betweenFilter: f })}
            styles={styles}
          />
        );
      case GADimensionFilterType.EMPTY:
        return <span className={styles.emptyLabel}>is empty / not set</span>;
      default:
        return null;
    }
  };

  return (
    <>
      <AsyncSelect
        loadOptions={loadFields}
        value={filter.fieldName ? { label: filter.fieldName, value: filter.fieldName } : null}
        onChange={(o) => onChange({ ...filter, fieldName: o?.value ?? '' })}
        placeholder="field name"
        allowCustomValue
        width={22}
        defaultOptions
        menuPlacement="bottom"
        isClearable
      />
      <Select<GADimensionFilterType>
        options={FILTER_TYPE_OPTIONS}
        value={filter.filterType}
        onChange={(o) => handleTypeChange(o.value!)}
        width={12}
        menuPlacement="bottom"
      />
      {renderParams()}
    </>
  );
};

// ─── main recursive component ─────────────────────────────────────────────────

export const GAFilterExpressionComponent: React.FC<GAFilterExpressionComponentProps> = ({
  expression,
  onChange,
  onDelete,
  loadFields,
  depth = 0,
}) => {
  const styles = useStyles2(getStyles);
  const currentType = exprType(expression);

  const handleTypeChange = (newType: string) => {
    switch (newType) {
      case 'none':
        return onChange({});
      case 'filter':
        return onChange(makeDefaultExpr());
      case 'andGroup':
        return onChange({ andGroup: { expressions: [makeDefaultExpr()] } });
      case 'orGroup':
        return onChange({ orGroup: { expressions: [makeDefaultExpr()] } });
      case 'notExpression':
        return onChange({ notExpression: makeDefaultExpr() });
    }
  };

  const typeSelector = (
    <Select
      options={EXPR_TYPE_OPTIONS}
      value={currentType}
      onChange={(o) => handleTypeChange(o.value!)}
      width={14}
      menuPlacement="bottom"
    />
  );

  // ── none ──
  if (currentType === 'none') {
    return (
      <div className={styles.row}>
        {typeSelector}
        {onDelete && (
          <IconButton name="times" tooltip="Remove" size="sm" variant="destructive" onClick={onDelete} />
        )}
      </div>
    );
  }

  // ── single filter ──
  if (currentType === 'filter') {
    const filter = expression.filter!;
    return (
      <div className={styles.row}>
        {typeSelector}
        <LeafFilterEditor
          filter={filter}
          onChange={(f) => onChange({ filter: f })}
          loadFields={loadFields}
          styles={styles}
        />
        {onDelete && (
          <IconButton name="times" tooltip="Remove" size="sm" variant="destructive" onClick={onDelete} />
        )}
      </div>
    );
  }

  // ── AND / OR group ──
  if (currentType === 'andGroup' || currentType === 'orGroup') {
    const list: GAFilterExpressionList = (expression as any)[currentType];

    const addChild = () => {
      const newList = { expressions: [...list.expressions, makeDefaultExpr()] };
      onChange({ [currentType]: newList });
    };

    const updateChild = (index: number, child: GAFilterExpression) => {
      const expressions = [...list.expressions];
      expressions[index] = child;
      onChange({ [currentType]: { expressions } });
    };

    const deleteChild = (index: number) => {
      const expressions = list.expressions.filter((_, i) => i !== index);
      onChange({ [currentType]: { expressions } });
    };

    return (
      <>
        <div className={styles.groupHeader}>
          {typeSelector}
          <IconButton name="plus" tooltip="Add expression" size="sm" onClick={addChild} />
          {onDelete && (
            <IconButton name="times" tooltip="Remove group" size="sm" variant="destructive" onClick={onDelete} />
          )}
        </div>
        <div className={styles.group(depth)}>
          {list.expressions.map((child, i) => (
            <GAFilterExpressionComponent
              key={i}
              expression={child}
              onChange={(c) => updateChild(i, c)}
              onDelete={() => deleteChild(i)}
              loadFields={loadFields}
              depth={depth + 1}
            />
          ))}
        </div>
      </>
    );
  }

  // ── NOT ──
  if (currentType === 'notExpression') {
    return (
      <>
        <div className={styles.groupHeader}>
          {typeSelector}
          {onDelete && (
            <IconButton name="times" tooltip="Remove" size="sm" variant="destructive" onClick={onDelete} />
          )}
        </div>
        <div className={styles.group(depth)}>
          <GAFilterExpressionComponent
            expression={expression.notExpression!}
            onChange={(child) => onChange({ notExpression: child })}
            loadFields={loadFields}
            depth={depth + 1}
          />
        </div>
      </>
    );
  }

  return null;
};
