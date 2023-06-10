import { QueryEditorProps } from "@grafana/data";
import { Cascader, CascaderOption } from "@grafana/ui";
import { DataSource } from "DataSource";
import React from "react";
import { GADataSourceOptions, GAQuery } from "types";

type Props = QueryEditorProps<DataSource, GAQuery, GADataSourceOptions>;

export const GACascader = (props: Props) => {
 const [updatedOption, setOptions ]  =  React.useState<CascaderOption[]>([])
 React.useEffect(()=>{
  Loading()
    // eslint-disable-next-line react-hooks/exhaustive-deps
 },[])

 const onSelect = (value: string) =>{
  console.log('value', value)
  props.onChange(props.query)
 }
const Loading = async () =>{
    const v =  await props.datasource.getCascader()
    console.log('v', v)
    // const t = [...v]
    setOptions(v)
  //   setOptions([
  //     {
  //         "label": "Default Account for Firebase",
  //         "value": "accounts/145710468",
  //         "items": [
  //             {
  //                 "label": "gitblog - GA4",
  //                 "value": "properties/323466308"
  //             }
  //         ]
  //     }
  // ])
    // props.onChange(props.query)
}
console.log('updatedOption', updatedOption)

 return <Cascader options={updatedOption} onSelect={onSelect} changeOnSelect={false} />
}

