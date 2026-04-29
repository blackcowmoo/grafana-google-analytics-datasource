[![CodeQL](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/actions/workflows/codeql-analysis.yml)  ![](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fblackcowmoo-googleanalytics-datasource%3Fversion%3Dlatest&query=%24.downloads&label=downloads) ![](https://img.shields.io/github/v/release/blackcowmoo/Grafana-Google-Analytics-DataSource?style=plastic?label=repo) ![](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fblackcowmoo-googleanalytics-datasource%3Fversion%3Dlatest&query=%24.version&label=grafana%20release&prefix=v)
# Google Analytics datasource

Visualize data from Google Analytics UA(Deprecated) And GA4(beta)

## Feature
- AutoComplete AccountID & WebpropertyID & ProfileID
- AutoComplete Metrics & Dimensions
- Query using Metrics & Dimensions
- Dimension filter with AND / OR / NOT groups (GA4)
- Dashboard variable support for `webPropertyId` and dimension-filter values
- Setting with json

![query](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/blob/master/src/img/query.png?raw=true)

## Preparations

Setup has three stages: create a Google Cloud service account, grant it access in Google Analytics, and upload the key file in Grafana. Do them in order.

### 1. Create a Google Cloud service account

1. If you don't already have a GCP project, [create one](https://cloud.google.com/resource-manager/docs/creating-managing-projects#console).
2. Open the [Credentials](https://console.developers.google.com/apis/credentials) page in the Google API Console.
3. Click **Create Credentials** → **Service account**.
4. Fill in the service account details and click **Create**.
5. On the **Service account permissions** page, leave the role empty (no project-level role is required) and click **Continue**.
6. Open the newly-created service account and click **Keys** → **Add Key** → **Create new key** → **JSON**. A JSON key file will download.
7. Note the service account email — it looks like `*@*.iam.gserviceaccount.com`. You'll need it in the next step.

### 2. Enable the required APIs

Enable the following APIs for the GCP project that owns the service account:

- **GA4 (recommended):** [Google Analytics Admin API](https://console.cloud.google.com/apis/library/analyticsadmin.googleapis.com) and [Google Analytics Data API](https://console.cloud.google.com/apis/library/analyticsdata.googleapis.com).
- **UA (deprecated, read-only):** [Google Analytics API](https://console.cloud.google.com/apis/library/analytics.googleapis.com) and [Google Analytics Reporting API](https://console.cloud.google.com/marketplace/product/google/analyticsreporting.googleapis.com).

Confirm they are enabled on the [API dashboard](https://console.cloud.google.com/apis/dashboard).

### 3. Grant the service account access in Google Analytics

1. Open [Google Analytics](https://analytics.google.com/) and select your account.
2. Open **Admin**.
3. **GA4:** open **Property access management** on the property you want to query. **UA:** open **Account User Management** on the account tab.
4. Click **+** → **Add users** and enter the service account email from step 1.7.
5. Grant the **Viewer** role (GA4) or **Read & Analyze** permission (UA).

### 4. Add the data source in Grafana

Go to **Add data source** in Grafana, pick this plugin, drag the JSON key file onto the upload area, and click **Save & Test**. The key is encrypted at rest in the Grafana database.

Common setup problems are covered in the [FAQ](./FAQ.md).

## FAQ
[FAQ](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/tree/master/FAQ.md)

## How To Dev
[build directory](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/tree/master/build)
