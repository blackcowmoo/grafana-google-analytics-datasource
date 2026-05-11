import { DataSourceInstanceSettings, MetricFindValue, ScopedVars, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { CascaderOption } from '@grafana/ui';
import { AccountSummary, GADataSourceOptions, GAMetadata, GAQuery } from './types';

export class DataSource extends DataSourceWithBackend<GAQuery, GADataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<GADataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: GAQuery, scopedVars: ScopedVars): Record<string, any> {
    const templateSrv = getTemplateSrv();
    let dimensionFilter = query.dimensionFilter;
    dimensionFilter?.orGroup?.expressions.map((expression) => {
      if (expression.filter?.stringFilter) {
        expression.filter.stringFilter.value = templateSrv.replace(expression.filter.stringFilter.value, scopedVars);
      }
      if (expression.filter?.inListFilter) {
        expression.filter.inListFilter.values = expression.filter.inListFilter.values.map((value) => {
          value = templateSrv.replace(value, scopedVars);
          return value;
        });
      }
      return expression;
    });

    // Apply template variable interpolation to webPropertyId
    const webPropertyId = templateSrv.replace(query.webPropertyId, scopedVars);

    return {
      ...query,
      webPropertyId,
      dimensionFilter,
    };
  }
  async getAccountSummaries(): Promise<CascaderOption[]> {
    let accountSummaries = (await this.getResource('account-summaries')).accountSummaries as AccountSummary[];
    let accounts: CascaderOption[] = [];
    for (const accountSummary of accountSummaries) {
      let accountCascader: CascaderOption = {
        label: accountSummary.DisplayName,
        value: accountSummary.Account,
      };
      let properties: CascaderOption[] = [];
      for (const propertySummary of accountSummary.PropertySummaries) {
        let propertyCascader: CascaderOption = {
          label: propertySummary.DisplayName,
          value: propertySummary.Property,
        };
        properties.push(propertyCascader);
        let profiles: CascaderOption[] = [];

        if (!propertySummary.ProfileSummaries) {
          continue;
        }
        for (const profileSummary of propertySummary.ProfileSummaries) {
          let profileCascader: CascaderOption = {
            label: profileSummary.DisplayName,
            value: profileSummary.Profile,
          };
          profiles.push(profileCascader);
        }
        propertyCascader.children = profiles;
        propertyCascader.items = profiles;
      }
      accountCascader.children = properties;
      accountCascader.items = properties;
      accounts.push(accountCascader);
    }
    return accounts;
  }

  async getTimezone(accountId: string, webPropertyId: string, profileId: string): Promise<string> {
    return this.getResource('profile/timezone', { accountId, webPropertyId, profileId }).then(({ timezone }) => {
      return timezone;
    });
  }

  async getServiceLevel(accountId: string, webPropertyId: string): Promise<string> {
    return this.getResource('property/service-level', { accountId, webPropertyId }).then(({ serviceLevel }) => {
      return serviceLevel;
    });
  }

  async getMetrics(query: string, webPropertyId: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('metrics', { webPropertyId }).then(({ metrics }) => {
      let test = metrics.reduce((pre: Array<SelectableValue<string>>, element: GAMetadata) => {
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
      return test;
    });
  }

  async getDimensions(query: string, exclude: any, webPropertyId: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('dimensions', { webPropertyId }).then(({ dimensions }) => {
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

  async getRealtimeMetrics(query: string, webPropertyId: string): Promise<Array<SelectableValue<string>>> {
    return this.getResource('realtime-metrics', { webPropertyId }).then(({ metrics }) => {
      let test = metrics.reduce((pre: Array<SelectableValue<string>>, element: GAMetadata) => {
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
      return test;
    });
  }

  async getRealtimeDimensions(
    query: string,
    exclude: any,
    webPropertyId: string
  ): Promise<Array<SelectableValue<string>>> {
    return this.getResource('realtime-dimensions', { webPropertyId }).then(({ dimensions }) => {
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

  async metricFindQuery(query: string, options?: { scopedVars?: ScopedVars }): Promise<MetricFindValue[]> {
    const templateSrv = getTemplateSrv();
    const interpolated = templateSrv.replace(query.trim(), options?.scopedVars);

    const accountSummaries = (await this.getResource('account-summaries')).accountSummaries as AccountSummary[];

    const propertiesMatch = interpolated.match(/^properties\(([^)]+)\)$/i);
    if (propertiesMatch) {
      const accountId = propertiesMatch[1].trim();
      const account = accountSummaries.find(
        (a) => a.Account === accountId || a.Account === `accounts/${accountId}`
      );
      return (account?.PropertySummaries ?? []).map((p) => ({
        text: p.DisplayName,
        value: p.Property,
      }));
    }

    if (/^properties$/i.test(interpolated)) {
      return accountSummaries.flatMap((a) =>
        (a.PropertySummaries ?? []).map((p) => ({
          text: p.DisplayName,
          value: p.Property,
        }))
      );
    }

    if (/^accounts$/i.test(interpolated)) {
      return accountSummaries.map((a) => ({
        text: a.DisplayName,
        value: a.Account,
      }));
    }

    return [];
  }

  async getTimeDimensions(): Promise<Array<SelectableValue<string>>> {
    return this.getDimensions('date', null, '');
  }

  async getDimensionsExcludeTimeDimensions(
    query: string,
    webPropertyId: string
  ): Promise<Array<SelectableValue<string>>> {
    return await this.getDimensions(query, 'date', webPropertyId);
  }
}
