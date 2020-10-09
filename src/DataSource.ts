import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { GADataSourceOptions, GAQuery } from './types';

export class DataSource extends DataSourceWithBackend<GAQuery, GADataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GADataSourceOptions>) {
    super(instanceSettings);
    console.log('instanceSettings', instanceSettings);
  }

  async getAccountId(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('accounts').then(({ account }) => {
      console.log(account);
      return account
        ? Object.entries(account).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : [];
    });
  }

  async getWebProperties(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('spreadsheets').then(({ viewIds }) =>
      viewIds ? Object.entries(viewIds).map(([value, label]) => ({ label, value } as SelectableValue<string>)) : []
    );
  }

  async getViewId(): Promise<Array<SelectableValue<string>>> {
    // let test = { aa: '123', bb: '456' };
    // let abc = Object.entries(test).map(([value, label]) => ({ label, value } as SelectableValue<string>));
    // return abc;
    return this.getResource('spreadsheets').then(({ viewIds }) => {
      console.log(viewIds);
      return viewIds
        ? Object.entries(viewIds).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : [];
    });
  }
}
