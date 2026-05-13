import { expect, test } from '@grafana/plugin-e2e';
import * as fs from 'fs';
import * as path from 'path';

const waitForPageSettle = async (page: import('@playwright/test').Page) => {
  await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {});
};

const FIXTURES_DIR = path.join(__dirname, '../fixtures');

// ─────────────────────────────────────────────────────────
// PHASE 1  (run with v0.3.1 dist)
// Creates representative dashboards and saves their UIDs
// to tests/fixtures/ for phase2 to consume.
// ─────────────────────────────────────────────────────────

test('phase1: v0.3.1 provisioned dashboards load without JS crash', async ({
  readProvisionedDashboard,
  page,
}) => {
  for (const version of ['v0.2.2', 'v0.2.3', 'v0.3.0']) {
    const dashboard = await readProvisionedDashboard({ fileName: `${version}.json` });
    const errors: string[] = [];
    page.on('console', (m) => {
      if (m.type() === 'error') {
        errors.push(m.text());
      }
    });

    await page.goto(`http://localhost:3000/d/${dashboard.uid}`);
    await waitForPageSettle(page);

    const critical = errors.filter(
      e => (e.includes('TypeError') || e.includes('Cannot read')) &&
           !e.includes('Failed to fetch') && !e.includes('Error loading recent')
    );
    expect(critical, `${version} JS errors on v0.3.1: ${critical.join('\n')}`).toHaveLength(0);
    console.log(`✅ ${version} OK on v0.3.1`);
  }
});

test('phase1: create mutation-artifact dashboard (simulates v0.3.1 filter mutation bug)', async ({
  request,
}) => {
  const ds = await request.get('http://localhost:3000/api/datasources/uid/lcc3108_test');
  const dsData = await ds.json();

  // v0.3.1 bug: after first query run with $country variable,
  // stringFilter.value gets mutated to the resolved value in the saved model.
  const dashboard = {
    dashboard: {
      title: 'v0.3.1 Migration — Filter Mutation Artifact',
      uid: 'v031-filter-artifact',
      tags: ['migration-test'],
      panels: [{
        id: 1,
        type: 'timeseries',
        title: 'GA4 with orGroup filter (post-mutation state)',
        gridPos: { h: 8, w: 12, x: 0, y: 0 },
        datasource: { type: dsData.type, uid: dsData.uid },
        targets: [{
          refId: 'A',
          version: 'v4',
          webPropertyId: 'properties/123456789',
          accountId: 'accounts/123',
          metrics: ['activeUsers'],
          timeDimension: 'date',
          dimensions: ['country'],
          mode: 'time series',
          // Simulates the post-mutation state from v0.3.1:
          // originally "$country" but got permanently overwritten to "United States"
          dimensionFilter: {
            orGroup: {
              expressions: [{
                filter: {
                  fieldName: 'country',
                  filterType: 'STRING',
                  stringFilter: {
                    matchType: 'EXACT',
                    value: 'United States',
                    caseSensitive: false,
                  },
                },
              }],
            },
          },
        }],
      }],
      time: { from: 'now-7d', to: 'now' },
      schemaVersion: 36,
    },
    folderId: 0,
    overwrite: true,
  };

  const resp = await request.post('http://localhost:3000/api/dashboards/db', {
    headers: { 'Content-Type': 'application/json' },
    data: JSON.stringify(dashboard),
  });
  expect(resp.ok()).toBeTruthy();
  const saved = await resp.json();

  fs.mkdirSync(FIXTURES_DIR, { recursive: true });
  fs.writeFileSync(
    path.join(FIXTURES_DIR, 'v031-artifact.json'),
    JSON.stringify({ uid: saved.uid, url: saved.url }, null, 2)
  );
  console.log(`✅ Phase1 — artifact saved: uid=${saved.uid}`);
});

// Also create a dashboard with andGroup / nested filter (new in current branch)
// to verify the reverse direction doesn't break v0.3.1 rendering.
test('phase1: check v0.3.1 handles unknown filter types gracefully', async ({
  request,
  page,
}) => {
  // Create a dashboard with an andGroup filter (not supported in v0.3.1 UI
  // but should still load without crashing since the backend ignores unknown fields).
  const ds = await request.get('http://localhost:3000/api/datasources/uid/lcc3108_test');
  const dsData = await ds.json();

  const dashboard = {
    dashboard: {
      title: 'Forward-compat: andGroup filter on v0.3.1',
      uid: 'v031-forward-compat',
      tags: ['migration-test'],
      panels: [{
        id: 1,
        type: 'timeseries',
        title: 'Panel with andGroup',
        gridPos: { h: 8, w: 12, x: 0, y: 0 },
        datasource: { type: dsData.type, uid: dsData.uid },
        targets: [{
          refId: 'A',
          version: 'v4',
          webPropertyId: 'properties/123456789',
          metrics: ['activeUsers'],
          timeDimension: 'date',
          mode: 'time series',
          dimensionFilter: {
            andGroup: {
              expressions: [
                {
                  filter: {
                    fieldName: 'country',
                    filterType: 'STRING',
                    stringFilter: { matchType: 'EXACT', value: 'United States', caseSensitive: false },
                  },
                },
              ],
            },
          },
        }],
      }],
      time: { from: 'now-7d', to: 'now' },
      schemaVersion: 36,
    },
    folderId: 0,
    overwrite: true,
  };

  const resp = await request.post('http://localhost:3000/api/dashboards/db', {
    headers: { 'Content-Type': 'application/json' },
    data: JSON.stringify(dashboard),
  });
  expect(resp.ok()).toBeTruthy();
  const saved = await resp.json();

  const errors: string[] = [];
  page.on('console', (m) => {
    if (m.type() === 'error') {
      errors.push(m.text());
    }
  });

  await page.goto(`http://localhost:3000/d/${saved.uid}`);
  await waitForPageSettle(page);

  const critical = errors.filter(
    e => (e.includes('TypeError') || e.includes('Cannot read')) &&
         !e.includes('Failed to fetch') && !e.includes('Error loading recent')
  );
  expect(critical, `Forward-compat JS errors: ${critical.join('\n')}`).toHaveLength(0);
  await expect(page.getByText('Panel plugin not found', { exact: false })).not.toBeVisible();

  fs.mkdirSync(FIXTURES_DIR, { recursive: true });
  fs.writeFileSync(
    path.join(FIXTURES_DIR, 'v031-forward-compat.json'),
    JSON.stringify({ uid: saved.uid }, null, 2)
  );
  console.log(`✅ Phase1 — forward-compat dashboard OK on v0.3.1`);
});

// ─────────────────────────────────────────────────────────
// PHASE 2  (run after swapping to current dist)
// Loads all phase1 dashboards and verifies they work.
// ─────────────────────────────────────────────────────────

const loadFixture = (name: string) => {
  const p = path.join(FIXTURES_DIR, name);
  if (!fs.existsSync(p)) { return null; }
  return JSON.parse(fs.readFileSync(p, 'utf-8'));
};

test('phase2: v0.3.1 mutation-artifact loads after upgrade', async ({ page }) => {
  const fixture = loadFixture('v031-artifact.json');
  if (!fixture) { test.skip(); return; }

  const errors: string[] = [];
  page.on('console', (m) => {
    if (m.type() === 'error') {
      errors.push(m.text());
    }
  });

  await page.goto(`http://localhost:3000/d/${fixture.uid}`);
  await waitForPageSettle(page);

  const critical = errors.filter(
    e => (e.includes('TypeError') || e.includes('Cannot read')) &&
         !e.includes('Failed to fetch') && !e.includes('Error loading recent')
  );
  expect(critical, `Upgrade JS errors:\n${critical.join('\n')}`).toHaveLength(0);
  await expect(page.getByText('Panel plugin not found', { exact: false })).not.toBeVisible();
  await expect(page).toHaveURL(/\/d\//);

  // Verify the filter UI shows the orGroup value (not crashed/blank)
  await page.goto(`http://localhost:3000/d/${fixture.uid}?orgId=1&editPanel=1`);
  await waitForPageSettle(page);
  // The filter expression should still render (current version must handle orGroup from v0.3.1)
  const filterFieldset = page.locator('[data-testid="filter-expression"]').first();
  await expect(filterFieldset).toBeVisible({ timeout: 15000 });

  console.log('✅ Phase2 — v0.3.1 mutation artifact loads correctly after upgrade');
});

test('phase2: forward-compat andGroup filter loads after upgrade', async ({ page }) => {
  const fixture = loadFixture('v031-forward-compat.json');
  if (!fixture) { test.skip(); return; }

  const errors: string[] = [];
  page.on('console', (m) => {
    if (m.type() === 'error') {
      errors.push(m.text());
    }
  });

  await page.goto(`http://localhost:3000/d/${fixture.uid}`);
  await waitForPageSettle(page);

  const critical = errors.filter(
    e => (e.includes('TypeError') || e.includes('Cannot read')) &&
         !e.includes('Failed to fetch') && !e.includes('Error loading recent')
  );
  expect(critical, `Forward-compat upgrade errors:\n${critical.join('\n')}`).toHaveLength(0);
  await expect(page.getByText('Panel plugin not found', { exact: false })).not.toBeVisible();
  console.log('✅ Phase2 — andGroup filter dashboard loads after upgrade');
});

test('phase2: provisioned dashboards all pass after upgrade', async ({
  readProvisionedDashboard,
  page,
}) => {
  for (const version of ['v0.2.2', 'v0.2.3', 'v0.3.0']) {
    const dashboard = await readProvisionedDashboard({ fileName: `${version}.json` });
    const errors: string[] = [];
    page.on('console', (m) => {
      if (m.type() === 'error') {
        errors.push(m.text());
      }
    });

    await page.goto(`http://localhost:3000/d/${dashboard.uid}`);
    await waitForPageSettle(page);

    const critical = errors.filter(
      e => (e.includes('TypeError') || e.includes('Cannot read')) &&
           !e.includes('Failed to fetch') && !e.includes('Error loading recent')
    );
    expect(critical, `${version} post-upgrade errors:\n${critical.join('\n')}`).toHaveLength(0);
    await expect(page.getByText('Panel plugin not found', { exact: false })).not.toBeVisible();
    console.log(`✅ ${version} OK after upgrade`);
  }
});
