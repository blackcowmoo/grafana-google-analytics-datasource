import { expect, test } from '@grafana/plugin-e2e';
import { GADataSourceOptions, GASecureJsonData } from '../../src/types';

test('"Save & test" should be successful when configuration is valid', async ({
  gotoDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({ fileName: 'datasources.yml' });
  const configPage = await gotoDataSourceConfigPage(ds.uid);
  const resetButton =  await page.getByText('Upload another JWT file').isVisible()
  if(resetButton){
    await page.getByText('Upload another JWT file').click()
  }
  await page.locator('input[accept="application/json"]').setInputFiles('./tests/credentials/grafana-success.json')
  await expect(configPage.saveAndTest()).toBeOK();
});

test('"Save & test" should fail when configuration is invalid', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource<GADataSourceOptions, GASecureJsonData>({ fileName: 'datasources.yml' });
  const configPage = await createDataSourceConfigPage({ type: ds.type, deleteDataSourceAfterTest: true });
  await page.locator('input[accept="application/json"]').setInputFiles('./tests/credentials/grafana-fail.json')
  await expect(configPage.saveAndTest()).not.toBeOK();
  await expect(configPage).toHaveAlert('error');
});


