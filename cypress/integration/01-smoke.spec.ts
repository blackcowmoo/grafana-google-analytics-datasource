import { e2e } from '@grafana/e2e';
import pluginJson from '../../src/plugin.json';

e2e.scenario({
  describeName: 'Smoke test',
  itName: 'Smoke test',
  scenario: () => {
    Cypress.config('defaultCommandTimeout', 10000);
    e2e().intercept('api/plugins/blackcowmoo-googleanalytics-datasource/markdown/README').as('pluginPageLoad')
    e2e().intercept(`/api/plugins/${pluginJson.id}/install`).as('pluginInstall')
    e2e().intercept('/api/datasources/uid/*/health').as('pluginHealth')

    // add plugin
    e2e().visit(`http://localhost:3000/plugins/${pluginJson.id}`).wait('@pluginPageLoad');
    // e2e().get('button span').contains('Install').click().wait('@pluginInstall');

    e2e.pages.AddDataSource.visit()
    // e2e().visit('http://localhost:3000/connections/datasources/new');
    e2e.pages.AddDataSource.dataSourcePluginsV2(pluginJson.name).click()
    // e2e().get(`[aria-label="Add new data source ${pluginJson.name}"]`).click();
    e2e().get('[accept="application/json"]').selectFile('/home/runner/work/grafana-google-analytics-datasource/grafana-google-analytics-datasource/test/credencials/auth.json', { action: "drag-drop", force: true });
    e2e().get('[data-testid="data-testid Data source settings page Save and Test button"]').scrollIntoView().click().wait('@pluginHealth', { log: true, timeout: 10000 }).wait(100)
    e2e().get('[data-testid="data-testid Data source settings page Save and Test button"]').scrollIntoView().screenshot('save')

    e2e().get('[data-testid="data-testid Alert success"] div').contains('Success').should('exist')

    e2e.flows.importDashboards('tests/dashboard')
  }
});