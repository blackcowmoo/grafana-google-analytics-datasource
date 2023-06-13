import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { CascaderOption } from '@grafana/ui';
import { AccountSummary, GADataSourceOptions, GAMetadata, GAQuery } from './types';

export class DataSource extends DataSourceWithBackend<GAQuery, GADataSourceOptions> {
  version: string;
  constructor(instanceSettings: DataSourceInstanceSettings<GADataSourceOptions>) {
    super(instanceSettings);
    console.log('instanceSettings', instanceSettings);
    this.version = instanceSettings.jsonData.version
  }

  async getAccountSummaries(): Promise<CascaderOption[]> {
    let accountSummaries = (await this.getResource('account-summaries')).accountSummaries as AccountSummary[]
    let accounts: CascaderOption[] = [];
    for (const accountSummary of accountSummaries) {
      let accountCascader: CascaderOption = {
        label: accountSummary.DisplayName,
        value: accountSummary.Account,
      }
      let properties: CascaderOption[] = [];
      for (const propertySummary of accountSummary.PropertySummaries) {
        let propertyCascader: CascaderOption = {
          label: propertySummary.DisplayName,
          value: propertySummary.Property,
        }
        properties.push(propertyCascader);
        let profiles: CascaderOption[] = [];

        if (!propertySummary.ProfileSummaries) {
          continue
        }
        for (const profileSummary of propertySummary.ProfileSummaries) {
          let profileCascader: CascaderOption = {
            label: profileSummary.DisplayName,
            value: profileSummary.Profile,
          }
          profiles.push(profileCascader);
        }
        propertyCascader.children = profiles;
        propertyCascader.items = profiles;
      }
      accountCascader.children = properties
      accountCascader.items = properties
      accounts.push(accountCascader);
    }
    return accounts;
  }

  async getTimezone(accountId: string, webPropertyId: string, profileId: string): Promise<string> {
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
}
