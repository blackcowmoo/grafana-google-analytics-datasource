# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

A Grafana datasource plugin for **Google Analytics 4** (GA4). The plugin is split into a TypeScript/React frontend (`src/`) and a Go backend (`pkg/`); a CLAUDE.md in each subdirectory has the details for that side.

- `src/CLAUDE.md` — frontend (ConfigEditor, QueryEditor, DataSource resource calls)
- `pkg/CLAUDE.md` — Go backend (auth, GA4 client, query/transform pipeline)
- `tests/CLAUDE.md` — Playwright e2e (provisioned datasource, JWT upload race, modal quirks)

The legacy GA3 / Universal Analytics datasource was removed in commit `e5bee9a`; only the GA4 path remains.

## Big picture

```
Grafana ──► Frontend (src/)                 Backend (pkg/)
            DataSource.ts ──HTTP─►          datasource.go (resource handlers
            QueryEditorGA4.tsx                + QueryData → analytics.go)
            ConfigEditor.tsx                       │
                                                   ▼
                                             gav4/ ──► Google APIs
                                             (admin v1beta, data v1beta)
                                                   ▲
                                             auth/Resolve  ──► tokenprovider
                                             (privateKey OR legacy jwt blob)
```

Resource endpoints registered in `pkg/datasource.go` (`/account-summaries`, `/dimensions`, `/metrics`, `/realtime-dimensions`, `/realtime-metrics`, `/profile/timezone`, `/property/service-level`) are how the frontend populates cascaders and async selects — keep names and shapes in sync between `src/DataSource.ts` and `pkg/datasource.go` when adding endpoints.

## Commands

Frontend (yarn 1.22, Node 20 — see `.nvmrc`):

```bash
yarn install                  # postinstall runs patch-package
yarn dev                      # webpack --watch (writes to ./dist)
yarn build                    # production build
yarn typecheck                # tsc --noEmit
yarn lint                     # eslint with cache
yarn test                     # jest --watch --onlyChanged
yarn test:ci                  # jest --passWithNoTests --maxWorkers 4
yarn e2e                      # playwright (requires Grafana on :3000 — see tests/CLAUDE.md)
yarn server                   # docker compose up --build (Grafana with plugin mounted)
```

Backend (Go 1.24+, mage):

```bash
mage -v                       # build for linux/darwin/windows (default = build:All)
mage -l                       # list mage targets
mage build:linux              # single-OS build (used in CI for e2e)
mage coverage                 # backend tests with coverage (CI uses this)
go test ./pkg/...             # run all backend tests
go test -run TestName ./pkg/auth   # single test
```

Local dev loop: `docker compose up` brings up Grafana 10.4 (override with `GRAFANA_VERSION=…`) with `./dist` mounted as the plugin and `./provisioning` as Grafana provisioning. Rebuild frontend with `yarn dev`, backend with `mage`, then `docker restart blackcowmoo-googleanalytics-datasource`.

## Repo-specific things to remember

- **Plugin ID** is `blackcowmoo-googleanalytics-datasource` — appears in `src/plugin.json`, docker-compose, Magefile output, and resource paths. Don't rename casually.
- **Grafana version floor**: `>=9.3.0` per `plugin.json`. Frontend pins `@grafana/data|runtime|ui` at `10.2.2`; bumping these requires checking the Grafana support matrix.
- **Provisioning** under `./provisioning` is mounted into the dev Grafana container — `provisioning/datasources/datasources.yml` is what the e2e tests load via `readProvisionedDataSource`.
- CI runs **typecheck + frontend build + `mage coverage` + `mage buildAll` + plugin-validator** on every PR (`.github/workflows/ci.yaml`); a separate matrix workflow (`playwright.yml`) runs e2e against multiple Grafana versions. Lint and unit tests are currently commented out in CI — local `yarn lint` / `yarn test:ci` is the only gate.
