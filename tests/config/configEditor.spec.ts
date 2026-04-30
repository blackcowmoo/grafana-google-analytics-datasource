import { expect, test } from '@grafana/plugin-e2e';
import type { Page } from '@playwright/test';
import { GADataSourceOptions, GASecureJsonData } from '../../src/types';
import { dismissWhatsNewModal } from '../utils';

// @grafana/google-sdk's <ConnectionConfig /> renders a FileDropzone inside an
// element marked `data-testid="Configuration drop zone"`. The dropzone holds a
// native <input type="file"> we can drive directly. If the datasource already
// has a JWT configured (e.g. left over from a previous test run) the SDK
// shows JWTForm with a Reset button instead — click it first to surface the
// dropzone again.
//
// setInputFiles dispatches the change event but resolves before the SDK's
// FileReader → onOptionsChange → React rerender chain completes. Wait for the
// dropzone to disappear (proves the JWT was accepted and credentials moved
// into secureJsonData) before invoking saveAndTest, otherwise the save races
// the credential update and CheckHealth returns 400.
const uploadJWT = async (page: Page, file: string) => {
  const reset = page.getByRole('button', { name: /^Reset/ });
  if (await reset.isVisible({ timeout: 1500 }).catch(() => false)) {
    await reset.click();
  }
  const dropzone = page.getByTestId('Configuration drop zone');
  await dropzone.locator('input[type="file"]').setInputFiles(file);
  await expect(dropzone).toBeHidden();
};

test('"Save & test" should be successful when configuration is valid', async ({
  gotoDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({ fileName: 'datasources.yml' });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  await dismissWhatsNewModal(page);
  await uploadJWT(page, './tests/credentials/grafana-success.json');
  await expect(configPage.saveAndTest()).toBeOK();
});

test('"Save & test" should fail when configuration is invalid', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({ fileName: 'datasources.yml' });
  const configPage = await createDataSourceConfigPage({ type: ds.type, deleteDataSourceAfterTest: true });
  await dismissWhatsNewModal(page);
  await uploadJWT(page, './tests/credentials/grafana-fail.json');
  await expect(configPage.saveAndTest()).not.toBeOK();
  await expect(configPage).toHaveAlert('error');
});

// Backwards compatibility: a datasource saved before gav3 was removed will
// still have `jsonData.version = "v3"` persisted. The plugin must load that
// config without crashing — the value is ignored and the GA4 path is used.
test('legacy v3 datasource config loads without crashing', async ({
  gotoDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({
    fileName: 'datasources.yml',
    name: ' Google Analytics (legacy v3)',
  });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  await dismissWhatsNewModal(page);
  // ConnectionConfig must mount — either the dropzone (no JWT yet) or the
  // JWTForm (legacy datasources may already have credentials persisted).
  const mounted = page
    .getByTestId('Configuration drop zone')
    .or(page.getByRole('button', { name: /^Reset/ }));
  await expect(mounted.first()).toBeVisible();
});
