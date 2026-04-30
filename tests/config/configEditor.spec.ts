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
const uploadJWT = async (page: Page, file: string) => {
  const reset = page.getByRole('button', { name: /^Reset/ });
  if (await reset.isVisible({ timeout: 1500 }).catch(() => false)) {
    await reset.click();
  }
  await page
    .getByTestId('Configuration drop zone')
    .locator('input[type="file"]')
    .setInputFiles(file);
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


