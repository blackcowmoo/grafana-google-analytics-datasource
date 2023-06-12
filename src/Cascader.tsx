import { QueryEditorProps } from "@grafana/data";
import { Cascader, CascaderOption } from "@grafana/ui";
import { DataSource } from "DataSource";
import React from "react";
import { GADataSourceOptions, GAQuery } from "types";

type Props = QueryEditorProps<DataSource, GAQuery, GADataSourceOptions>;

export const GACascader = (props: Props) => {
  const [updatedOption, setOptions] = React.useState<CascaderOption[]>([])
  const [select, setSelect] = React.useState<string>("")
  React.useEffect(() => {
    Loading()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])
  React.useEffect(() => {
    props.query.webPropertyId = select
    props.onChange(props.query)
  }, [select])
  const onSelect = (value: string) => {
    console.log('value', value)
    setSelect(value)

  }
  const Loading = async () => {
    const v = await props.datasource.getAccountSummaries()
    console.log('v', v)
    setOptions(v)

  }
  console.log('updatedOption', updatedOption)

  return <Cascader options={updatedOption} onSelect={onSelect} changeOnSelect={true} displayAllSelectedLevels={true} />
}

