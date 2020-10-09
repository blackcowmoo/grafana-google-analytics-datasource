import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { GADataSourceOptions, GAQuery } from './types';

export class DataSource extends DataSourceWithBackend<GAQuery, GADataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GADataSourceOptions>) {
    super(instanceSettings);
    console.log('instanceSettings', instanceSettings);
  }

  async getAccountId(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('accounts').then(({ accounts }) => {
      console.log(accounts);
      return accounts
        ? Object.entries(accounts).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : [];
    });
  }

  async getWebProperties(accountId: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('web-properties', { accountId }).then(({ webProperties }) =>
      webProperties
        ? Object.entries(webProperties).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : []
    );
  }

  async getViewId(): Promise<Array<SelectableValue<string>>> {
    // let test = { aa: '123', bb: '456' };
    // let abc = Object.entries(test).map(([value, label]) => ({ label, value } as SelectableValue<string>));
    // return abc;
    return this.getResource('profiles').then(({ profiles }) => {
      return profiles
        ? Object.entries(profiles).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : [];
    });
  }
}
