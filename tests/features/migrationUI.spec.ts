import { expect, test } from '@grafana/plugin-e2e';
import * as fs from 'fs';
import * as path from 'path';

// Migration UI tests: verify that old dashboard JSON files render without
// JavaScript errors and without a "panel error" when upgrading from v0.3.1.
// These do NOT require real GA credentials - they check that the plugin
// handles old query formats gracefully (no crash / no unhandled JS error).

const waitForPageSettle = async (page: import('@playwright/test').Page) => {
  await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {});
};

const getDashboardVersions = () => {
  const dashboardsDir = path.join(__dirname, '../../provisioning/dashboards');
  return fs
    .readdirSync(dashboardsDir)
    .filter((f) => f.match(/^v\d+\.\d+\.\d+\.json$/))
    .map((f) => f.replace('.json', ''));
};

getDashboardVersions().forEach((version) => {
  test(`${version}: dashboard loads without JS crash`, async ({
    readProvisionedDashboard,
    gotoDashboardPage,
    page,
  }) => {
    // Collect console errors
    const jsErrors: string[] = [];
    page.on('console', (msg) => {
      if (msg.type() === 'error') {
        jsErrors.push(msg.text());
      }
    });

    const dashboard = await readProvisionedDashboard({ fileName: `${version}.json` });
    await gotoDashboardPage({ uid: dashboard.uid });
    await waitForPageSettle(page);

    // No unhandled React/JS errors from the plugin (ignore Grafana internal network errors)
    const reactErrors = jsErrors.filter(
      (e) =>
        (e.includes('Uncaught') || e.includes('TypeError') || e.includes('Cannot read')) &&
        !e.includes('Error loading recent dashboard actions') &&
        !e.includes('Failed to fetch')
    );
    expect(reactErrors, `JS errors on ${version} dashboard: ${reactErrors.join(', ')}`).toHaveLength(0);
  });

  test(`${version}: panels render (no "panel plugin not found" error)`, async ({
    readProvisionedDashboard,
    gotoDashboardPage,
    page,
  }) => {
    const dashboard = await readProvisionedDashboard({ fileName: `${version}.json` });
    await gotoDashboardPage({ uid: dashboard.uid });
    await waitForPageSettle(page);

    // Plugin-not-found would show this text
    const pluginError = page.getByText('Panel plugin not found', { exact: false });
    await expect(pluginError).not.toBeVisible();

    // Page URL should remain on the dashboard (no redirect to error page)
    await expect(page).toHaveURL(/\/d\//i);
  });
});

test('default dashboard loads without JS crash', async ({
  readProvisionedDashboard,
  gotoDashboardPage,
  page,
}) => {
  const jsErrors: string[] = [];
  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      jsErrors.push(msg.text());
    }
  });

  const dashboard = await readProvisionedDashboard({ fileName: 'default.json' });
  await gotoDashboardPage({ uid: dashboard.uid });
  await waitForPageSettle(page);

  const reactErrors = jsErrors.filter(
    (e) =>
      (e.includes('Uncaught') || e.includes('TypeError') || e.includes('Cannot read')) &&
      !e.includes('Error loading recent dashboard actions') &&
      !e.includes('Failed to fetch')
  );
  expect(reactErrors, `JS errors on default dashboard: ${reactErrors.join(', ')}`).toHaveLength(0);
});
