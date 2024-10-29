// cursor ai로 작성됨
import { SelectableValue } from '@grafana/data';
import { Button, FieldSet, HorizontalGroup, Input, RadioButtonGroup, Select, VerticalGroup } from '@grafana/ui';
import React from 'react';
import { GADimensionFilterType, GAFilter, GAFilterExpression, GAFilterExpressionList, GAInListFilter, GAStringFilter, GAStringFilterMatchType } from 'types';

interface Props {
  expression: GAFilterExpression;
  onChange: (expression: GAFilterExpression) => void;
  onDelete?: () => void;  // New prop added
  selectedDimensions: Array<SelectableValue<string>>; // New prop added
}

export const GAFilterExpressionComponent: React.FC<Props> = ({ expression = {}, onChange, onDelete, selectedDimensions }) => {
  const expressionTypes: Array<SelectableValue<string>> = [
    { label: 'no filter', value: 'none' },
    { label: 'AND', value: 'andGroup' },
    { label: 'OR', value: 'orGroup' },
    { label: 'NOT', value: 'notExpression' },
    { label: 'FILTER', value: 'filter' },
  ];

  const handleExpressionTypeChange = (option: SelectableValue<string>) => {
    let newExpression: GAFilterExpression;
    switch (option.value) {
      case 'none':
        newExpression = {};
        break;
      case 'andGroup':
      case 'orGroup':
        newExpression = { [option.value]: { expressions: [] } };
        break;
      case 'notExpression':
        newExpression = { notExpression: {} };
        break;
      case 'filter':
        newExpression = { 
          filter: { 
            fieldName: '', 
            filterType: GADimensionFilterType.STRING,
            stringFilter: { matchType: GAStringFilterMatchType.EXACT, value: '', caseSensitive: false }
          } 
        };
        break;
      default:
        return;
    }
    onChange(newExpression);
  };

  const renderExpressionContent = () => {
    if (Object.keys(expression).length === 0) {
      return <div>No filter set</div>;
    } else if (expression.andGroup) {
      return renderExpressionList(expression.andGroup, 'andGroup');
    } else if (expression.orGroup) {
      return renderExpressionList(expression.orGroup, 'orGroup');
    } else if (expression.notExpression) {
      return renderNotExpression();
    } else if (expression.filter) {
      return renderFilter();
    }
    return null;
  };

  const renderExpressionList = (list: GAFilterExpressionList, type: 'andGroup' | 'orGroup') => {
    return (
      <VerticalGroup>
        {list.expressions.map((expr, index) => (
          <HorizontalGroup key={index}>
            <GAFilterExpressionComponent
              expression={expr}
              selectedDimensions={selectedDimensions}
              onChange={(newExpr) => {
                const newList = { ...list, expressions: [...list.expressions] };
                newList.expressions[index] = newExpr;
                onChange({ [type]: newList });
              }}
              onDelete={() => {
                const newList = { ...list, expressions: [...list.expressions] };
                newList.expressions.splice(index, 1);
                onChange({ [type]: newList });
              }}
            />
          </HorizontalGroup>
        ))}
        <HorizontalGroup>
          <Button
            onClick={() => {
              const newList = { ...list, expressions: [...list.expressions, {}] };
              onChange({ [type]: newList });
            }}
          >
            Add Expression
          </Button>
        </HorizontalGroup>
      </VerticalGroup>
    );
  };

  const renderNotExpression = () => {
    return (
      <GAFilterExpressionComponent
        expression={expression.notExpression!}
        selectedDimensions={selectedDimensions}
        onChange={(newExpr) => onChange({ notExpression: newExpr })}
        onDelete={() => onChange({})}
      />
    );
  };

  const renderFilter = () => {
    const filter = expression.filter || { fieldName: '', filterType: GADimensionFilterType.STRING };

    const filterTypes: Array<SelectableValue<string>> = [
      { label: 'String', value: 'STRING' },
      { label: 'List', value: 'IN_LIST' },
    ];

    const handleFilterTypeChange = (option: SelectableValue<string>) => {
      let newFilter: GAFilter = {
        ...filter,
        filterType: option.value as GADimensionFilterType,
      };

      // Set initial values based on filter type
      switch (option.value) {
        case GADimensionFilterType.STRING:
          newFilter.stringFilter = { matchType: GAStringFilterMatchType.EXACT, value: '', caseSensitive: false };
          newFilter.inListFilter = undefined;
          break;
        case GADimensionFilterType.IN_LIST:
          newFilter.inListFilter = { values: [], caseSensitive: false };
          newFilter.stringFilter = undefined;
          break;
      }

      onChange({ filter: newFilter });
    };

    const renderFilterContent = () => {
      switch (filter.filterType) {
        case GADimensionFilterType.STRING:
          return renderStringFilter(filter.stringFilter!);
        case GADimensionFilterType.IN_LIST:
          return renderInListFilter(filter.inListFilter!);
        default:
          return null;
      }
    };

    return (
      <VerticalGroup>
        <Select
          options={selectedDimensions}
          value={filter.fieldName}
          onChange={(option) => onChange({ filter: { ...filter, fieldName: option.value! } })}
          placeholder="Select field"
        />
        <Select
          options={filterTypes}
          value={filter.filterType}
          onChange={handleFilterTypeChange}
        />
        {renderFilterContent()}
      </VerticalGroup>
    );
  };

  const renderStringFilter = (stringFilter: GAStringFilter) => {
    const matchTypes = Object.values(GAStringFilterMatchType).map(value => ({ label: value, value }));

    return (
      <VerticalGroup>
        <Select
          options={matchTypes}
          value={stringFilter.matchType}
          onChange={(option) => onChange({ filter: { ...expression.filter!, stringFilter: { ...stringFilter, matchType: option.value! } } })}
        />
        <Input
          value={stringFilter.value}
          onChange={(e) => onChange({ filter: { ...expression.filter!, stringFilter: { ...stringFilter, value: e.currentTarget.value } } })}
          placeholder="Value"
        />
        <RadioButtonGroup
          options={[
            { label: 'Case sensitive', value: true },
            { label: 'Case insensitive', value: false },
          ]}
          value={stringFilter.caseSensitive}
          onChange={(value) => onChange({ filter: { ...expression.filter!, stringFilter: { ...stringFilter, caseSensitive: value } } })}
        />
      </VerticalGroup>
    );
  };

  const renderInListFilter = (inListFilter: GAInListFilter = { values: [], caseSensitive: false }) => {
    return (
      <VerticalGroup>
        <Input
          value={inListFilter.values.join(', ')}
          onChange={(e) => onChange({ filter: { ...expression.filter!, inListFilter: { ...inListFilter, values: e.currentTarget.value.split(',').map(v => v.trim()) } } })}
          placeholder="Values (comma-separated)"
        />
        <RadioButtonGroup
          options={[
            { label: 'Case sensitive', value: true },
            { label: 'Case insensitive', value: false },
          ]}
          value={inListFilter.caseSensitive}
          onChange={(value) => onChange({ filter: { ...expression.filter!, inListFilter: { ...inListFilter, caseSensitive: value } } })}
        />
      </VerticalGroup>
    );
  };

  return (
    <FieldSet>
      <HorizontalGroup>
        <Select
          options={expressionTypes}
          value={Object.keys(expression).length === 0 ? 'none' : Object.keys(expression)[0]}
          onChange={handleExpressionTypeChange}
        />
        {renderExpressionContent()}
        {onDelete && (
          <Button variant="destructive" onClick={onDelete}>
            Delete
          </Button>
        )}
      </HorizontalGroup>
    </FieldSet>
  );
};
