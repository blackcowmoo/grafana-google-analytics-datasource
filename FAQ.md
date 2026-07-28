## Configuration error

1. [google api check](https://console.cloud.google.com/apis/dashboard)
  ![img](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/blob/master/img/api-check.png?raw=true)
2. [service account check](https://console.cloud.google.com/apis/credentials)
  ![img](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/blob/master/img/serviceaccount.png?raw=true)
3. [analytics account access check](https://analytics.google.com/analytics)

## Query problems

### `Web property ID` — what format?

- **GA4:** `properties/XXXXXXXXX` (numeric property id prefixed with `properties/`). You can copy this from GA4 Admin → Property Settings.
- **UA (deprecated):** numeric web property id, e.g. `UA-XXXXXXX-1`.

### Why does my `Last 6 hours` panel show data for the whole day?

Google Analytics reporting APIs only accept day-resolution date ranges, so the plugin pulls full days from GA and filters the response down to the Grafana time range. To make this work correctly, pick a time dimension with the right resolution:

- `dateHour` (or `dateHourMinute` for GA4) for sub-day ranges.
- `date` for multi-day ranges.

The plugin drops buckets that fall outside the requested window; choosing a coarser dimension (e.g. `date` for a 6-hour range) may therefore produce empty results.

### Why does my query fail with `Cohort metrics can only be used in requests with a CohortSpec`?

Cohort metrics (any metric under the **Cohort** group) require a `cohortSpec` in the request payload, which the plugin does not yet build automatically. Until cohort support lands, select metrics from other groups.

### Dashboard variables

The plugin interpolates Grafana variables in:

- **Web property ID** — allows one panel repeated across multiple GA4 properties. Use a multi-value variable and Grafana's *Repeat panel* feature.
- **Dimension filter** — both `STRING` and `IN_LIST` filter values. For multi-select variables with `IN_LIST`, the plugin expands the variable into one filter value per selected item; just write the plain variable reference, e.g. `$campaigns`.
- **`filtersExpression`** (UA) — plain string substitution.

For string filters that need to match any of several values, use the `FULL_REGEXP` match type with Grafana's regex format specifier, e.g. `${myVar:regex}`.

## Service account & permissions

### Which APIs must be enabled on the GCP project?

- **GA4:** `analyticsadmin.googleapis.com` and `analyticsdata.googleapis.com`.
- **UA:** `analytics.googleapis.com` and `analyticsreporting.googleapis.com`.

### What GA role does the service account need?

Add the service account email to your Google Analytics property with the **Viewer** (GA4) or **Read & Analyze** (UA) role. No role at the GCP project level is required.

### The data source test says "Not Exist Valid Profile"

The service account was successfully authenticated but has no visibility into any GA4 property / UA view. Re-check:

1. The service account email (ending in `iam.gserviceaccount.com`) is added as a user on the GA property.
2. For GA4, the property is reachable via the Admin API (not just the Data API) — the admin API lists properties and is used by the account picker.

## Dev & troubleshooting

### Enable plugin debug logs

Set `log.level=debug` in `grafana.ini` to see the plugin's request/response logs:

```
[log]
level = debug
```

See the [Grafana logging reference](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/#log) for more detail.
