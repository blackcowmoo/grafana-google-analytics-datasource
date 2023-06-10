import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { CascaderOption } from '@grafana/ui';
import { GADataSourceOptions, GAMetadata, GAQuery } from './types';

export class DataSource extends DataSourceWithBackend<GAQuery, GADataSourceOptions> {
  version: string;
  constructor(instanceSettings: DataSourceInstanceSettings<GADataSourceOptions>) {
    super(instanceSettings);
    console.log('instanceSettings', instanceSettings);
    this.version = instanceSettings.jsonData.version
  }

  async getAccountIds(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('accounts').then(({ accounts }) => {
      return accounts
        ? Object.entries(accounts).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : [];
    });
  }

  async getWebPropertyIds(accountId: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('web-properties', { accountId }).then(({ webProperties }) =>
      webProperties
        ? Object.entries(webProperties).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : []
    );
  }

  async getProfileIds(accountId: string, webPropertyId: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('profiles', { accountId, webPropertyId }).then(({ profiles }) => {
      return profiles
        ? Object.entries(profiles).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : [];
    });
  }

  async getProfileTimezone(accountId: string, webPropertyId: string, profileId: string): Promise<string> {
    return this.getResource('profile/timezone', { accountId, webPropertyId, profileId }).then(({ timezone }) => {
      return timezone;
    });
  }

  async getMetrics(query: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('metrics').then(({ metrics }) => {
      return metrics.reduce((pre: Array<SelectableValue<string>>, element: GAMetadata) => {
        if (
          element.id.toLowerCase().indexOf(query) > -1 ||
          element.attributes.uiName.toLowerCase().indexOf(query) > -1
        ) {
          pre.push({
            label: element.attributes.uiName,
            value: element.id,
            description: element.attributes.description,
          } as SelectableValue<string>);
        }
        return pre;
      }, []);
    });
  }

  async getDimensions(query: string, exclude: any): Promise<Array<SelectableValue<string>>> {
    return this.getResource('dimensions').then(({ dimensions }) => {
      return dimensions.reduce((pre: Array<SelectableValue<string>>, element: GAMetadata) => {
        if (
          (element.id.toLowerCase().indexOf(query) > -1 ||
            element.attributes.uiName.toLowerCase().indexOf(query) > -1) &&
          !(
            element.id.toLowerCase().indexOf(exclude) > -1 ||
            element.attributes.uiName.toLowerCase().indexOf(exclude) > -1
          )
        ) {
          pre.push({
            label: element.attributes.uiName,
            value: element.id,
            description: element.attributes.description,
          } as SelectableValue<string>);
        }
        return pre;
      }, []);
    });
  }

  async getTimeDimensions(): Promise<Array<SelectableValue<string>>> {
    return this.getDimensions('date', null);
  }

  async getDimensionsExcludeTimeDimensions(query: string): Promise<Array<SelectableValue<string>>> {
    return await this.getDimensions(query, 'date');
  }
  getGaVersion(): string {
    return this.version
  }

  async getCascader(): Promise<CascaderOption[]>{
    let cascaderOption: CascaderOption[] = []
    let accounts = await this.getAccountIds()
    let properties
    let profiles
    accounts.map(async(account, _)=>{
      let accountId = account.value || ""
      let accountOption: CascaderOption = {
        label: account.label || "",
        value: account.value
      }
      properties = await this.getWebPropertyIds(accountId)
      properties.map(async(property,_)=>{
        let propertyId = property.value || ""
        let propertyOption: CascaderOption = {
          label: property.label || "",
          value: property.value
        }
        profiles = await this.getProfileIds(accountId,propertyId)
        profiles.map((profile)=>{
          let profleOption: CascaderOption = {
            label: profile.label || "",
            value: profile.value
          }
          propertyOption.items = []
          propertyOption.items.push(profleOption)
        })
        accountOption.items = []
        accountOption.items.push(propertyOption)
      })
      cascaderOption.push(accountOption)
    })
    console.log('cascaderOption', cascaderOption)
    return cascaderOption
  }
}
