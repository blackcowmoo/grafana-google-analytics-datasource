# pkg/ — Go backend

Backend for the GA4 datasource plugin. Entry point `main.go` calls `datasource.Manage(...)` from the Grafana plugin SDK; the lifecycle (`NewDataSource`, `CheckHealth`, `QueryData`, `CallResource`) is implemented in `datasource.go`.

## Layout

```
pkg/
├── main.go              # plugin entry. blank-imports time/tzdata for Windows tz fix (#101).
├── datasource.go        # implements backend.* interfaces; mounts resource HTTP routes.
├── analytics.go         # GoogleAnalytics interface — the seam between datasource.go and gav4/.
├── auth/                # normalises plugin settings → token provider HTTP client.
├── setting/             # decodes jsonData + DecryptedSecureJSONData into DatasourceSecretSettings.
├── model/               # shared types (QueryModel, MetadataItem, AccountSummary, ColumnDefinition).
├── gav4/                # the only GA implementation (GA3 was removed in e5bee9a).
└── util/                # generics (TypeConverter[R]), time helpers, Elapsed timer.
```

## Request flow

1. **Resource calls** (`CallResource` → `mux.HandleFunc`): the frontend hits `/account-summaries`, `/dimensions`, `/metrics`, `/realtime-{dimensions,metrics}`, `/profile/timezone`, `/property/service-level`. Each handler decodes settings via `setting.LoadSettings`, calls a method on `ds.analytics` (a `gav4.GoogleAnalytics`), and writes the result with `writeResult` (always wraps payload in `{<key>: value}` or `{error: ...}`).
2. **Queries** (`QueryData`): for each query, `gav4.GoogleAnalytics.Query` builds a `model.QueryModel` from the JSON, validates it requires `WebPropertyID` + (`Dimensions` or `Metrics`), and dispatches to `client.getReport` or `client.getRealtimeReport` based on `Mode`.
3. **CheckHealth**: hits `GetAccountSummaries` then runs a small canned report (`active1DayUsers` over yesterday→today on the first property) — if either step fails, the message includes the error string.

## Auth (`pkg/auth`)

`auth.Resolve(*setting.DatasourceSecretSettings) (*Resolved, error)` is the **single normalisation point**. It produces a `Resolved` descriptor consumed by `auth.NewHTTPClient`, which builds an `*http.Client` whose middleware injects an OAuth2 token (`grafana-google-sdk-go/pkg/tokenprovider`).

Key behavior — **dual-read fallback** (`pkg/auth/auth.go:resolveJWT`):

- New datasources persist explicit fields: `clientEmail`, `tokenUri` (jsonData) + `privateKey` (secureJsonData). These come from `<ConnectionConfig />` in the frontend.
- Legacy datasources (pre @grafana/google-sdk adoption — PR #163) only have `secureJsonData.jwt` containing the full service-account JSON blob. `resolveJWT` parses that blob if `privateKey` is empty, so existing datasources keep working without any user action.
- `Resolve` returns an error for `gce` / `workloadIdentityFederation` — the constants exist (`auth.TypeGCE`, `auth.TypeWIF`) but the token-provider switch in `provider.go` only implements `TypeJWT`. Adding support means extending both `Resolve` and `newTokenProvider`.

When touching auth, **don't drop the legacy `jwt` read path** — it's the migration ramp for existing users.

## Caching

`gav4.GoogleAnalytics` holds a `*cache.Cache` (patrickmn/go-cache) created in `NewDataSource` with `300s` default expiration. TTLs vary by resource:

- `analytics:accountsummaries:<jwt>` — 60s. Note: the cache key embeds `config.JWT`, so legacy and new-style settings have different keys; this is intentional and changing it can collide cross-tenant.
- `analytics:account:<a>:webproperty:<wp>:profile:<p>:timezone` — 60s
- `analytics:account:<a>:webproperty:<wp>:service_level` — 60s
- `ga:metadata:<propertyId>:metrics` — 1h. **The `:dimensions` variant is not written back to the cache** (only read) — likely a bug; preserve current behavior unless explicitly fixing.
- `ga:metadata:<propertyId>:realtime-{metrics,dimensions}` — read-through cache around the static `GaRealTimeMetrics` / `GaRealTimeDimensions` slices in `gav4/const.go`.

## GA4 client (`pkg/gav4`)

- `client.go` wraps two Google services (`analyticsdata/v1beta` and `analyticsadmin/v1beta`), each with its own scoped HTTP client.
- `getReport` paginates by recursing while `RowCount > Offset + GaReportMaxResult` (100000). Realtime uses `getRealtimeReport`, which clamps `MinuteRanges` to `[0, 29]` for standard or `[0, 59]` for `ServiceLevelPremium` (GA360).
- `grafana.go` transforms `analyticsdata.RunReportResponse` into `data.Frame`s. Three modes (`pkg/model/models.go`):
  - `TIME_SERIES` — pivots on `TimeDimension`. **Sub-day Grafana time ranges** are respected (PR #156 / commit 6330265).
  - `TABLE` — flat table with one column per dimension/metric.
  - `REALTIME` — uses `RunRealtimeReport` shape; `getReport` is called via `util.TypeConverter[RunReportResponse]` to reuse the same transform path.
- `model.go` (`GetQueryModel`): unmarshals `query.JSON` into `model.QueryModel` and computes `StartDate`/`EndDate` in the **query's timezone** (loaded via `time.LoadLocation` — relies on the tzdata blank import in `main.go`). If `TimeDimension` is set it is **prepended** to `Dimensions` before the API call.

## Tests

`auth_test.go`, `gav4/grafana_test.go`, `model/models_test.go`, `util/util_test.go`. Run with `go test ./pkg/...`. CI runs `mage coverage`. There are **no integration tests against the live Google API** in the Go suite — those live in the e2e Playwright suite under `tests/`.
