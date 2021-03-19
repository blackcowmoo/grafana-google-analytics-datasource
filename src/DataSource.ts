import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { GADataSourceOptions, GAQuery } from './types';

export class DataSource extends DataSourceWithBackend<GAQuery, GADataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GADataSourceOptions>) {
    super(instanceSettings);
    console.log('instanceSettings', instanceSettings);
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

  async getMetrics(query?: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('metrics').then(({ metrics }) => {
      return metrics.reduce((pre: Array<SelectableValue<string>>, element: any) => {
        let id = element.id as string;
        if (query && id.toLowerCase().indexOf(query) > -1) {
          pre.push({
            label: element.id,
            value: element.id,
          } as SelectableValue<string>);
        }
        return pre;
      }, []);
    });
  }

  async getDimensions(query?: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('dimensions').then(({ dimensions }) => {
      return dimensions.reduce((pre: Array<SelectableValue<string>>, element: any) => {
        let id = element.id as string;
        if (query && id.toLowerCase().indexOf(query) > -1) {
          pre.push({
            label: element.id,
            value: element.id,
          } as SelectableValue<string>);
        }
        return pre;
      }, []);
    });
  }
}
