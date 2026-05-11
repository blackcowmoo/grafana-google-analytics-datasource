# tests/ — Playwright e2e

End-to-end tests using `@grafana/plugin-e2e` + Playwright. Configured in `playwright.config.ts` at the repo root with three projects: `auth` (provided by `@grafana/plugin-e2e`, writes admin storage state), then `config` and `query` which depend on it. CI runs the suite in `.github/workflows/playwright.yml` against a matrix of Grafana versions.

## Running locally

```bash
docker compose up -d   # Grafana on :3000 with the plugin and provisioning mounted
yarn run e2e           # headless
yarn run e2e:ui        # Playwright UI mode (good for debugging selectors)
```

The tests load datasources via `readProvisionedDataSource({ fileName: 'datasources.yml' })`, which reads `provisioning/datasources/datasources.yml` — this file is mounted into the dev container by `docker-compose.yaml`.

## Credentials

`tests/credentials/grafana-success.json` and `grafana-fail.json` are gitignored service-account JSONs. CI restores them from base64-encoded GitHub secrets (`GOOGLE_AUTH_JSON`, `GOOGLE_AUTH_FAIL_JSON`). Locally you have to drop your own files at those paths before the suite will pass — `grafana-success.json` must point at a service account with read access to a GA4 property; `grafana-fail.json` should be a JSON that fails CheckHealth (e.g. valid format, no GA permissions).

## Two non-obvious quirks (`tests/utils.ts`, `tests/config/configEditor.spec.ts`)

Both are explained at length in inline comments — read them before changing the helpers, because they encode race conditions that have already bitten the suite.

1. **Grafana 13+ "What's new" modal** — mounts a viewport-sized scrim that intercepts pointer events. Any click before dismissal silently times out. Always call `dismissWhatsNewModal(page)` after navigating to a config/query page. The helper returns quickly on older Grafana versions where the dialog isn't rendered.

2. **JWT upload race** — `<ConnectionConfig />`'s FileDropzone fires `setInputFiles`, but the SDK's `FileReader → onOptionsChange → React rerender` chain is async. If you `saveAndTest()` immediately, the save races the credential update and CheckHealth returns 400. The `uploadJWT` helper waits for `dropzone` to become **hidden** (proof the JWT was accepted into `secureJsonData`) before returning. If a previous test left credentials behind, the SDK shows a `Reset` button instead of the dropzone — `uploadJWT` clicks it first.

## Backwards-compat tests

`configEditor.spec.ts` includes a `legacy v3 datasource config loads without crashing` test. The provisioning file has a `Google Analytics (legacy v3)` entry with `jsonData.version = "v3"` left over from the old GA3 datasource — the assertion is that the config page mounts without crashing, since the version field is now ignored and only the GA4 path runs. Don't delete this case without removing the corresponding datasource from `provisioning/datasources/datasources.yml`.

## Adding new specs

Put them in `tests/config/` or `tests/query/`. The project dependency chain (`auth` → `config` → `query`) means `query` specs run after `config` specs and reuse the same admin storage state. If you need a clean datasource per test, use `createDataSourceConfigPage({ ..., deleteDataSourceAfterTest: true })` instead of `gotoDataSourceConfigPage`.
