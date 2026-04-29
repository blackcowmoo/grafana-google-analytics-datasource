import { expect, test } from '@grafana/plugin-e2e';
import type { Locator } from '@playwright/test';
import * as fs from 'fs';
import * as path from 'path';
import { dismissWhatsNewModal } from '../utils';

// Grafana 10.4.5+ renders the RadioButtonGroup with a visible <label> and a
// stacked opacity:0 <input>, so Playwright's `.check()` on the native input
// times out (waits for visibility) and a plain `.click()` on the label is
// blocked by the input's pointer-event interception. Use `force: true` to
// click the visible label directly. Fall back to the old plain-text path for
// pre-10.4.5 Grafana (no label-wrapped radios).
//
// The `count()` check is a snapshot of the current DOM and does not wait -
// plugin-e2e >= 2.0 no longer uses networkidle, so the datasource editor can
// still be rendering when we reach this helper. Wait for either the radio
// group or the legacy text to actually appear before branching.
const selectQueryMode = async (row: Locator, mode: string) => {
  const label = row
    .getByLabel('query-mode')
    .locator('label')
    .filter({ hasText: new RegExp(`^${mode}\\s*$`) });
  const legacyText = row.getByText(mode, { exact: true });
  await expect(label.or(legacyText).first()).toBeVisible({ timeout: 15000 });
  if ((await label.count()) > 0) {
    await label.click({ force: true });
    return;
  }
  await legacyText.check();
};

// plugin-e2e >= 2.0 dropped `waitUntil: 'networkidle'` from GrafanaPage.navigate()
// in favor of `'load'`, so when gotoDashboardPage/explorePage returns the initial
// panel queries may still be in flight and the scenes-based dashboard (Grafana 11+)
// may not be ready to react to refresh/datasource selection. Replicate the v1.6.x
// behavior by explicitly waiting for the network to settle; the catch swallows the
// harmless timeout on dashboards that keep polling (e.g. live data).
const waitForPageSettle = async (page: import('@playwright/test').Page) => {
  await page.waitForLoadState('networkidle', { timeout: 10000 }).catch(() => {});
};

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
  test(`${version} migration test`, async ({ readProvisionedDataSource, readProvisionedDashboard, gotoDashboardPage, page }) => {
    const dashboard = await readProvisionedDashboard({fileName: `${version}.json`});
    const dashboardPage = await gotoDashboardPage({uid: dashboard.uid});
    await dismissWhatsNewModal(page);
    // Wait for the initial dashboard render/hydration to settle. Without this,
    // Grafana 11/12 scenes are not yet ready to react to the refresh click and
    // no /api/ds/query request is emitted.
    await waitForPageSettle(page);
    // Attach the listener BEFORE clicking refresh so we don't miss the response
    // to a fast-firing query.
    const responsePromise = dashboardPage.waitForQueryDataResponse();
    await dashboardPage.refreshDashboard();
    await expect(responsePromise).toBeOK();
  });
});

test('time series', async ({ readProvisionedDataSource, explorePage, page }) => {
  // default settings
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });

  await dismissWhatsNewModal(page);
  await explorePage.datasource.set(ds.name);
  await waitForPageSettle(page);
  await explorePage.timeRange.set({ from: 'now-7d', to: 'now' });

  await selectQueryMode(explorePage.getQueryEditorRow('A'), 'Time Series');

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

  await dismissWhatsNewModal(page);
  await explorePage.datasource.set(ds.name);
  await waitForPageSettle(page);
  await explorePage.timeRange.set({ from: 'now-7d', to: 'now' });
  await selectQueryMode(explorePage.getQueryEditorRow('A'), 'Table');
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


test('realtime', async ({ readProvisionedDataSource, explorePage, page }) => {
  // default settings
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });

  await dismissWhatsNewModal(page);
  await explorePage.datasource.set(ds.name);
  await waitForPageSettle(page);
  await explorePage.timeRange.set({ from: 'now-7d', to: 'now' });
  await selectQueryMode(explorePage.getQueryEditorRow('A'), 'Realtime');
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
