import { expect, test } from '@grafana/plugin-e2e';
import * as fs from 'fs';
import * as path from 'path';

// 대시보드 버전 목록 가져오기
const getDashboardVersions = () => {
  const dashboardsDir = path.join(__dirname, '../../provisioning/dashboards');
  const files = fs.readdirSync(dashboardsDir);
  return files
    .filter(file => file.match(/^v\d+\.\d+\.\d+\.json$/))
    .map(file => file.replace('.json', ''));
};

// 각 버전별 마이그레이션 테스트
getDashboardVersions().forEach(version => {
  test(`${version} migration test`, async ({ readProvisionedDataSource, readProvisionedDashboard, gotoDashboardPage }) => {
    const dashboard = await readProvisionedDashboard({fileName: `${version}.json`});
    const dashboardPage = await gotoDashboardPage({uid: dashboard.uid});
    await dashboardPage.refreshDashboard();
    await expect(dashboardPage.waitForQueryDataResponse()).toBeOK();
  });
});

test('time series', async ({ readProvisionedDataSource, explorePage, page }) => {
  // default settings
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });

  await explorePage.datasource.set(ds.name);
  await explorePage.timeRange.set({ from: 'now-7d', to: 'now' });

  let queryMode =  explorePage.getQueryEditorRow('A').getByLabel('query-mode').getByLabel('Time Series')
  // for grafana version < 10.4.5
  if(await queryMode.count()==0){
    queryMode =  explorePage.getQueryEditorRow('A').getByText('Time Series',{exact: true})
  }
  await queryMode.check()

  // account select
  await explorePage.getQueryEditorRow('A').getByRole('button', { name: 'Account Select' }).click();
  await page.getByText('Default Account for Firebase').click();
  await page.getByText('gitblog - GA4').click();
  //  await page.waitForResponse((response) => response.url().includes('resources/property/service-level'));
  await expect(explorePage.getQueryEditorRow('A').getByLabel('account-info')).toContainText(/.*properties\/.*/);
  // metrics select
  await explorePage.getQueryEditorRow('A').getByLabel('metrics').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Active users', { exact: true }).click();
  // time dimension
  await explorePage.getQueryEditorRow('A').getByLabel('time-dimension').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Date + hour (YYYYMMDDHH)', { exact: true }).click();

  // dimensions
  await explorePage.getQueryEditorRow('A').getByLabel('dimensions').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Country', { exact: true }).click();

  await expect(explorePage.timeSeriesPanel.waitForQueryDataResponse()).toBeOK()
  // await page.getByRole('combobox', { name: 'Query type' }).click();
  // await panelEditPage.getByGrafanaSelector(selectors.components.Select.option).getByText('Table').click();
  // await expect(panelEditPage.panel.fieldNames).toHaveText(['time', 'temperature outside', 'temperature inside']);
});


test('table', async ({ readProvisionedDataSource, explorePage, page }) => {
  // default settings
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });

  await explorePage.datasource.set(ds.name);
  await explorePage.timeRange.set({ from: 'now-7d', to: 'now' });
  let queryMode =  explorePage.getQueryEditorRow('A').getByLabel('query-mode').getByLabel('Table')
  // for grafana version < 10.4.5
  if(await queryMode.count()==0){
    queryMode =  explorePage.getQueryEditorRow('A').getByText('Table',{exact: true})
  }
  await queryMode.check()  // account select
  await explorePage.getQueryEditorRow('A').getByRole('button', { name: 'Account Select' }).click();
  await page.getByText('Default Account for Firebase').click();
  await page.getByText('gitblog - GA4').click();
  //  await page.waitForResponse((response) => response.url().includes('resources/property/service-level'));
  await expect(explorePage.getQueryEditorRow('A').getByLabel('account-info')).toContainText(/.*properties\/.*/);
  // metrics select
  await explorePage.getQueryEditorRow('A').getByLabel('metrics').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Active users', { exact: true }).click();
  // time dimension
  await explorePage.getQueryEditorRow('A').getByLabel('time-dimension').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Date + hour (YYYYMMDDHH)', { exact: true }).click();

  // dimensions
  await explorePage.getQueryEditorRow('A').getByLabel('dimensions').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Country', { exact: true }).click();

  await expect(explorePage.timeSeriesPanel.waitForQueryDataResponse()).toBeOK()
  // await page.getByRole('combobox', { name: 'Query type' }).click();
  // await panelEditPage.getByGrafanaSelector(selectors.components.Select.option).getByText('Table').click();
  // await expect(panelEditPage.panel.fieldNames).toHaveText(['time', 'temperature outside', 'temperature inside']);
});


test('realtime', async ({ readProvisionedDataSource, explorePage, page }) => {
  // default settings
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });

  await explorePage.datasource.set(ds.name);
  await explorePage.timeRange.set({ from: 'now-7d', to: 'now' });
  let queryMode =  explorePage.getQueryEditorRow('A').getByLabel('query-mode').getByLabel('Realtime')
  // for grafana version < 10.4.5
  if(await queryMode.count()==0){
    queryMode =  explorePage.getQueryEditorRow('A').getByText('Realtime',{exact: true})
  }
  await queryMode.check()  // account select
  await explorePage.getQueryEditorRow('A').getByRole('button', { name: 'Account Select' }).click();
  await page.getByText('Default Account for Firebase').click();
  await page.getByText('gitblog - GA4').click();
  //  await page.waitForResponse((response) => response.url().includes('resources/property/service-level'));
  await expect(explorePage.getQueryEditorRow('A').getByLabel('account-info')).toContainText(/.*properties\/.*/);
  // metrics select
  await explorePage.getQueryEditorRow('A').getByLabel('metrics').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Active users', { exact: true }).click();
  // time dimension
  await expect(explorePage.getQueryEditorRow('A').getByLabel('time-dimension')).toBeDisabled();

  // dimensions
  await explorePage.getQueryEditorRow('A').getByLabel('dimensions').click();
  await expect(page.getByLabel('Select options menu')).toBeVisible();
  await page.getByLabel('Select options menu').getByText('Country', { exact: true }).click();

  await expect(explorePage.timeSeriesPanel.waitForQueryDataResponse()).toBeOK()
  // await page.getByRole('combobox', { name: 'Query type' }).click();
  // await panelEditPage.getByGrafanaSelector(selectors.components.Select.option).getByText('Table').click();
  // await expect(panelEditPage.panel.fieldNames).toHaveText(['time', 'temperature outside', 'temperature inside']);
});
