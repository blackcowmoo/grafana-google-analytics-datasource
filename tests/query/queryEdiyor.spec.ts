import { expect, test } from '@grafana/plugin-e2e';

test('time series', async ({ readProvisionedDataSource, explorePage, page }) => {
  // default settings
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });

  await explorePage.datasource.set(ds.name);
  await explorePage.timeRange.set({ from: 'now-7d', to: 'now' });

  // account select
  await explorePage.getQueryEditorRow('A').getByRole('button', { name: 'Account Select' }).click();
  await page.getByText('Default Account for Firebase').click();
  await page.getByText('gitblog - GA4').click();
  await page.waitForResponse((response) => response.url().includes('resources/property/service-level'));
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
  await explorePage.getQueryEditorRow('A').getByText('Table',{exact: true})
  // account select
  await explorePage.getQueryEditorRow('A').getByRole('button', { name: 'Account Select' }).click();
  await page.getByText('Default Account for Firebase').click();
  await page.getByText('gitblog - GA4').click();
  await page.waitForResponse((response) => response.url().includes('resources/property/service-level'));
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
  await explorePage.getQueryEditorRow('A').getByText('Realtime',{exact: true})
  // account select
  await explorePage.getQueryEditorRow('A').getByRole('button', { name: 'Account Select' }).click();
  await page.getByText('Default Account for Firebase').click();
  await page.getByText('gitblog - GA4').click();
  await page.waitForResponse((response) => response.url().includes('resources/property/service-level'));
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
