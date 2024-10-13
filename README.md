[![CodeQL](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/actions/workflows/codeql-analysis.yml)  ![](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fblackcowmoo-googleanalytics-datasource%3Fversion%3Dlatest&query=%24.downloads&label=downloads) ![](https://img.shields.io/github/v/release/blackcowmoo/Grafana-Google-Analytics-DataSource?style=plastic?label=repo) ![](https://img.shields.io/badge/dynamic/json?url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fblackcowmoo-googleanalytics-datasource%3Fversion%3Dlatest&query=%24.version&label=grafana%20release&prefix=v)
# Google Analytics datasource

Visualize data from Google Analytics UA(Deprecated) And GA4(beta)

## Feature
- AutoComplete AccountID & WebpropertyID & ProfileID
- AutoComplete Metrics & Dimensions
- Query using Metrics & Dimensions
- Setting with json

![query](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/blob/master/src/img/query.png?raw=true)

## Preparations
### Generate a JWT file

1.  if you don't have gcp project, add new gcp project. [link](https://cloud.google.com/resource-manager/docs/creating-managing-projects#console)
2.  Open the [Credentials](https://console.developers.google.com/apis/credentials) page in the Google API Console.
3.  Click **Create Credentials** then click **Service account**.
4.  On the Create service account page, enter the Service account details.
5.  On the `Create service account` page, fill in the `Service account details` and then click `Create`
6.  On the `Service account permissions` page, don't add a role to the service account. Just click `Continue`
7.  In the next step, click `Create Key`. Choose key type `JSON` and click `Create`. A JSON key file will be created and downloaded to your computer
8.  Note your `service account email` ex) *@*.iam.gserviceaccount.com
9.  Open the [Google Analytics Admin API](https://console.cloud.google.com/apis/library/analyticsadmin.googleapis.com)  in API Library and enable access for your account
10. Open the [Google Analytics Data API](https://console.cloud.google.com/apis/library/analyticsdata.googleapis.com)  in API Library and enable access for your GA Data

### Google Analytics Setting

1. Open the [Google Analytics](https://analytics.google.com/)
2. Select Your Analytics Account And Open Admin Page
3. Click **Account User Management** on the **Account Tab**
4. Click plus Button then Add users
5. Enter `service account email` at **Generate a JWT file** 8th step and Permissions add `Read & Analyze`

### Grafana
Go To Add Data source then Drag the file to the dotted zone above. Then click `Save & Test`.   
The file contents will be encrypted and saved in the Grafana database.

## FAQ
[FAQ](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/tree/master/FAQ.md)

## How To Dev
[build directory](https://github.com/blackcowmoo/Grafana-Google-Analytics-DataSource/tree/master/build)
