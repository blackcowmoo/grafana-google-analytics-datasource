# src/ — Frontend

TypeScript/React frontend for the GA4 datasource. Bundled by webpack (`./.config/webpack/webpack.config.ts`); the entry is `module.ts`. The project uses **path-based imports without `./`** (e.g. `import { DataSource } from 'DataSource'`) — webpack/tsconfig resolve from `src/` as a root, so don't "fix" them to relative paths.

## Wiring

`module.ts` registers a `DataSourcePlugin` with three pieces:

- **Datasource class** (`DataSource.ts`) — extends `DataSourceWithBackend`. Aside from boilerplate it does two things: `applyTemplateVariables` interpolates Grafana variables into `dimensionFilter` (string and inList filters) and `webPropertyId`; and a small wrapper layer (`getMetrics`, `getDimensions`, `getRealtimeMetrics`, `getRealtimeDimensions`, `getAccountSummaries`, `getTimezone`, `getServiceLevel`) calls `this.getResource(...)` against the backend resource paths defined in `pkg/datasource.go`. Any new metadata endpoint needs to be added on **both sides** with matching response keys (`{accountSummaries: ...}`, `{dimensions: ...}`, etc.).
- **ConfigEditor** (`ConfigEditor.tsx`) — wraps `<ConnectionConfig />` from `@grafana/google-sdk` (adopted in PR #163). The SDK component owns the JWT upload UI and persists `clientEmail`/`tokenUri` (jsonData) + `privateKey` (secureJsonData). A small `usingLegacyJWT` check shows an info `<Alert>` when the datasource still has `secureJsonData.jwt` from the pre-SDK era — the backend's `auth.Resolve` falls back to parsing that blob, so existing datasources keep working until the user re-uploads.
- **QueryEditor** (`QueryEditorCommon.tsx` → `QueryEditorGA4.tsx`) — `QueryEditorCommon` is a thin shim that always renders GA4 (the GA3 branch was removed). `QueryEditorGA4` owns the cascader for account/property/profile, the metric/time-dimension/dimension async selects, the dimension filter (`Filter.tsx`), and the `Time Series / Table / Realtime` mode switch.

## Types (`types.ts`)

- `GADataSourceOptions extends GoogleDataSourceOptions` and `GASecureJsonData extends GoogleDataSourceSecureJsonData`. The base types come from `@grafana/google-sdk` so `<ConnectionConfig />` reads/writes the same fields the backend sees. `GASecureJsonData` adds the legacy `jwt?: string` so backwards-compat code can detect it.
- `GAQuery` is the per-panel query shape. Note the `displayName: Map<string, string>` field — set in `QueryEditorGA4`'s constructor; it's an in-memory editor concern, not persisted via JSON.
- `GAFilterExpression` mirrors the GA Data API filter shape one-to-one (see comment linking to the API docs). The mode-switch code in `QueryEditorGA4.onModeChange` only resets `timeDimension` / `selectedTimeDimensions` when switching to `realtime`.

## QueryEditorGA4 specifics

- `willRunQuery` is **debounced 500ms** and only fires `onRunQuery` when `webPropertyId && metrics && (mode === 'table' || mode === 'realtime' || timeDimension)` — i.e. time-series mode requires a time dimension before running.
- The `key={mode + webPropertyId + 'metrics'}` pattern on `<AsyncMultiSelect>` forces the loaders to remount whenever mode or property changes — this is intentional, so the cached `loadOptions` results don't leak between modes.
- `getDimensionsExcludeTimeDimensions` (in `DataSource.ts`) reuses `getDimensions` with `exclude='date'`, which excludes anything containing `date` from the regular dimensions selector — so the time dimension picker and dimension picker can both populate from the same backend endpoint without duplicates.

## Filter component

`Filter.tsx` is recursive (`andGroup` / `orGroup` / `notExpression` nest themselves). Per the tooltip in `QueryEditorGA4`, "Currently, only `or groups` are supported" — this comment refers to backend behavior at the time of writing; check `pkg/gav4/client.go` if you need to confirm what the backend actually forwards (it now forwards `DimensionFilter` whenever any of `OrGroup`/`AndGroup`/`Filter`/`NotExpression` is non-nil).

## Build / tests

- Webpack + jest configs are imported from `./.config/` (Grafana plugin scaffolding) — modify those at the root if you need to change loader behavior. The local `jest.config.js` only re-exports it and forces `TZ=UTC` (`jest-setup.js`).
- There are no `*.test.ts` files in `src/` currently. Frontend correctness is checked with `yarn typecheck` and the Playwright e2e in `tests/`.
