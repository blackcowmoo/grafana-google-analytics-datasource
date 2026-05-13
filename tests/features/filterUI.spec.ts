import { expect, test } from '@grafana/plugin-e2e';
import { dismissWhatsNewModal } from '../utils';

// Tests for the new Filter UI component (AND/OR/NOT/FILTER tree editor)
// These tests verify UI interactions without requiring real GA credentials.

const waitForPageSettle = async (page: import('@playwright/test').Page) => {
  await page.waitForLoadState('networkidle', { timeout: 10000 }).catch(() => {});
};

test.describe('DimensionFilter UI', () => {
  test('filter component renders with "no filter" default', async ({
    readProvisionedDataSource,
    explorePage,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
    await explorePage.datasource.set(ds.name);
    await dismissWhatsNewModal(page);
    await waitForPageSettle(page);

    const queryRow = explorePage.getQueryEditorRow('A');

    // The filter expression container should be present and default to 'No filter'
    const filterExpr = queryRow.locator('[data-testid="filter-expression"]').first();
    await expect(filterExpr).toBeVisible({ timeout: 15000 });

    const noFilterOption = filterExpr.locator('[class*="singleValue"]').first();
    await expect(noFilterOption).toContainText(/no filter/i);
  });

  test('can switch to OR group filter', async ({
    readProvisionedDataSource,
    explorePage,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
    await explorePage.datasource.set(ds.name);
    await dismissWhatsNewModal(page);
    await waitForPageSettle(page);

    const queryRow = explorePage.getQueryEditorRow('A');
    const filterExpr = queryRow.locator('[data-testid="filter-expression"]').first();
    await expect(filterExpr).toBeVisible({ timeout: 15000 });

    // Open the filter type dropdown and select OR group
    await filterExpr.locator('[class*="grafana-select-value-container"]').first().click();
    await expect(page.getByLabel('Select options menu')).toBeVisible();
    await page.getByLabel('Select options menu').getByText('OR group', { exact: true }).click();

    // After selecting OR group, "Add expression" icon button should appear
    await expect(filterExpr.getByRole('button', { name: 'Add expression' })).toBeVisible();
  });

  test('can add a Filter expression inside OR group', async ({
    readProvisionedDataSource,
    explorePage,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
    await explorePage.datasource.set(ds.name);
    await dismissWhatsNewModal(page);
    await waitForPageSettle(page);

    const queryRow = explorePage.getQueryEditorRow('A');
    const filterExpr = queryRow.locator('[data-testid="filter-expression"]').first();
    await expect(filterExpr).toBeVisible({ timeout: 15000 });

    // Switch to OR group (adds one default child expression automatically)
    await filterExpr.locator('[class*="grafana-select-value-container"]').first().click();
    await page.getByLabel('Select options menu').getByText('OR group', { exact: true }).click();

    // Add another child expression
    await filterExpr.getByRole('button', { name: 'Add expression' }).click();

    // Child expressions should appear (default type is Filter with a string filter)
    const childExprs = filterExpr.locator('[data-testid="filter-expression"]');
    await expect(childExprs.first()).toBeVisible();

    // The value input for the string filter should be present
    await expect(childExprs.first().getByPlaceholder('value or $variable')).toBeVisible();
  });

  test('can switch to AND group filter', async ({
    readProvisionedDataSource,
    explorePage,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
    await explorePage.datasource.set(ds.name);
    await dismissWhatsNewModal(page);
    await waitForPageSettle(page);

    const queryRow = explorePage.getQueryEditorRow('A');
    const filterExpr = queryRow.locator('[data-testid="filter-expression"]').first();
    await expect(filterExpr).toBeVisible({ timeout: 15000 });

    // Switch to AND group
    await filterExpr.locator('[class*="grafana-select-value-container"]').first().click();
    await page.getByLabel('Select options menu').getByText('AND group', { exact: true }).click();
    await expect(filterExpr.getByRole('button', { name: 'Add expression' })).toBeVisible();
  });

  test('can switch to NOT expression filter', async ({
    readProvisionedDataSource,
    explorePage,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
    await explorePage.datasource.set(ds.name);
    await dismissWhatsNewModal(page);
    await waitForPageSettle(page);

    const queryRow = explorePage.getQueryEditorRow('A');
    const filterExpr = queryRow.locator('[data-testid="filter-expression"]').first();
    await expect(filterExpr).toBeVisible({ timeout: 15000 });

    // Switch to NOT expression
    await filterExpr.locator('[class*="grafana-select-value-container"]').first().click();
    await page.getByLabel('Select options menu').getByText('NOT', { exact: true }).click();

    // NOT wraps a child expression — inner filter-expression should appear
    await expect(filterExpr.locator('[data-testid="filter-expression"]').first()).toBeVisible();
  });

  test('can add and delete filter expressions', async ({
    readProvisionedDataSource,
    explorePage,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
    await explorePage.datasource.set(ds.name);
    await dismissWhatsNewModal(page);
    await waitForPageSettle(page);

    const queryRow = explorePage.getQueryEditorRow('A');
    const filterExpr = queryRow.locator('[data-testid="filter-expression"]').first();
    await expect(filterExpr).toBeVisible({ timeout: 15000 });

    // Switch to OR group and add two more expressions (one is added automatically)
    await filterExpr.locator('[class*="grafana-select-value-container"]').first().click();
    await page.getByLabel('Select options menu').getByText('OR group', { exact: true }).click();
    await filterExpr.getByRole('button', { name: 'Add expression' }).click();
    await filterExpr.getByRole('button', { name: 'Add expression' }).click();

    // Should have at least 2 child filter-expression elements
    const childCount = await filterExpr.locator('[data-testid="filter-expression"]').count();
    expect(childCount).toBeGreaterThanOrEqual(2);

    // Delete first child via its Remove button
    const deleteBtn = filterExpr.locator('[data-testid="filter-expression"]').first().getByRole('button', { name: 'Remove' });
    await deleteBtn.click();

    // Should now have fewer children
    const newChildCount = await filterExpr.locator('[data-testid="filter-expression"]').count();
    expect(newChildCount).toBeLessThan(childCount);
  });

  test('filter supports string match types', async ({
    readProvisionedDataSource,
    explorePage,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
    await explorePage.datasource.set(ds.name);
    await dismissWhatsNewModal(page);
    await waitForPageSettle(page);

    const queryRow = explorePage.getQueryEditorRow('A');
    const filterExpr = queryRow.locator('[data-testid="filter-expression"]').first();
    await expect(filterExpr).toBeVisible({ timeout: 15000 });

    // Switch to OR group and add a child expression (already a Filter type with string filter)
    await filterExpr.locator('[class*="grafana-select-value-container"]').first().click();
    await page.getByLabel('Select options menu').getByText('OR group', { exact: true }).click();

    const childExpr = filterExpr.locator('[data-testid="filter-expression"]').first();
    await expect(childExpr).toBeVisible();

    // Match type select is the 4th select-value-container in the child:
    // nth(0)=expr type, nth(1)=field name AsyncSelect, nth(2)=filter type, nth(3)=match type
    const matchTypeSelect = childExpr.locator('[class*="grafana-select-value-container"]').nth(3);
    await expect(matchTypeSelect).toContainText(/exact|contains|begins with|ends with/i);

    // Switch to Contains
    await matchTypeSelect.click();
    await page.getByLabel('Select options menu').getByText('Contains', { exact: true }).click();
    await expect(matchTypeSelect).toContainText(/contains/i);
  });
});
