import { expect, test } from '@grafana/plugin-e2e';
import { GADataSourceOptions, GASecureJsonData } from '../../src/types';
import { dismissWhatsNewModal } from '../utils';

test('"Save & test" should be successful when configuration is valid', async ({
  gotoDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({ fileName: 'datasources.yml' });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  await dismissWhatsNewModal(page);
  const resetButton =  await page.getByText('Upload another JWT file').isVisible()
  if(resetButton){
    await page.getByText('Upload another JWT file').click()
  }
  await page.locator('input[type="file"][accept*="application/json"]').setInputFiles('./tests/credentials/grafana-success.json')
  // setInputFiles dispatches the change event but resolves before the dropzone's
  // FileReader → JWTConfig.onChange → React rerender chain completes. Wait for
  // the dropzone to disappear (proves onChange fired and secureJsonData was
  // populated) before triggering saveAndTest, otherwise the save races the
  // credential update and CheckHealth returns 400.
  await expect(page.getByText('Drop the file here, or click to use the file explorer')).toBeHidden();
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
  await page.locator('input[type="file"][accept*="application/json"]').setInputFiles('./tests/credentials/grafana-fail.json')
  await expect(page.getByText('Drop the file here, or click to use the file explorer')).toBeHidden();
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
  // The JWT upload control must render — proves ConfigEditor mounted despite
  // the legacy `version: v3` value in jsonData.
  await expect(page.locator('input[type="file"][accept*="application/json"]')).toBeAttached();
});


