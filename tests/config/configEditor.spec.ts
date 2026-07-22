import { expect, test } from '@grafana/plugin-e2e';
import type { Page } from '@playwright/test';
import { GADataSourceOptions, GASecureJsonData } from '../../src/types';
import { dismissWhatsNewModal } from '../utils';

// @grafana/google-sdk's <ConnectionConfig /> renders a FileDropzone inside
// `data-testid="Configuration drop zone"` only when isUploading=true (the
// SDK's default for a brand-new datasource). If a previous test run left JWT
// credentials behind, the SDK shows JWTForm instead and the dropzone is not
// rendered until the user clicks "Upload JWT Token".
//
// setInputFiles dispatches the change event but the SDK's
// FileReader → onOptionsChange → React rerender chain is async. Wait for the
// dropzone to disappear (proves the JWT was accepted into secureJsonData)
// before invoking saveAndTest, otherwise the save races the credential update
// and CheckHealth returns 400.
const uploadJWT = async (page: Page, file: string) => {
  const dropzone = page.getByTestId('Configuration drop zone');
  const uploadBtn = page.getByTestId('Upload JWT button');

  // ConnectionConfig sets authenticationType=JWT on mount via a useEffect,
  // which triggers a re-render before the JWT section appears. Wait for either
  // the dropzone (fresh/cleared datasource — isUploading defaults to true) or
  // the "Upload JWT Token" button (JWTForm shown when credentials already set).
  await expect(dropzone.or(uploadBtn).first()).toBeVisible({ timeout: 15000 });

  if (!await dropzone.isVisible().catch(() => false)) {
    await uploadBtn.click();
    await expect(dropzone).toBeVisible({ timeout: 10000 });
  }

  await dropzone.locator('input[type="file"]').setInputFiles(file);
  await expect(dropzone).toBeHidden({ timeout: 10000 });
};

test('"Save & test" should be successful when configuration is valid', async ({
  gotoDataSourceConfigPage,
  readProvisionedDataSource,
  page,
  request,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({ fileName: 'datasources.yml' });
  // Clear any previously saved JWT so ConnectionConfig renders the dropzone
  // immediately (isUploading defaults to true for a fresh datasource config).
  await request.put(`http://localhost:3000/api/datasources/uid/${ds.uid}`, {
    headers: { 'Content-Type': 'application/json' },
    data: JSON.stringify({
      name: ds.name,
      type: ds.type,
      access: 'proxy',
      jsonData: {},
      secureJsonData: {},
    }),
  });
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

// GCE and Workload Identity Federation are recognized by @grafana/google-sdk
// but pkg/auth/auth.go returns a "not yet supported" error for both, so
// ConfigEditor.tsx renders <AuthConfig /> with a JWT-only authOptions list
// instead of the SDK's default <ConnectionConfig />. This asserts the GCE
// radio option stays hidden so users can't pick an auth type that always
// fails Save & test.
test('auth type selector only offers JWT (GCE is not implemented on the backend)', async ({
  gotoDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({ fileName: 'datasources.yml' });
  await gotoDataSourceConfigPage(ds.uid);
  await dismissWhatsNewModal(page);
  await expect(page.getByText('Google JWT File')).toBeVisible();
  await expect(page.getByText('GCE Default Service Account')).toHaveCount(0);
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
    .or(page.getByTestId('Upload JWT button'));
  await expect(mounted.first()).toBeVisible({ timeout: 15000 });
});
