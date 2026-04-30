import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { ConnectionConfig } from '@grafana/google-sdk';
import { Alert, RadioButtonGroup } from '@grafana/ui';
import React, { useEffect } from 'react';
import { GADataSourceOptions, GASecureJsonData } from 'types';

const gaVersion = [
  { label: 'UA(GA3)', value: 'v3' },
  { label: 'GA4', value: 'v4' },
];

export type Props = DataSourcePluginOptionsEditorProps<GADataSourceOptions, GASecureJsonData>;

export const ConfigEditor: React.FC<Props> = (props) => {
  const { options, onOptionsChange } = props;
  const { jsonData } = options;

  // Default GA version to v4 (GA4) on first mount, immutably.
  useEffect(() => {
    if (!jsonData.version) {
      onOptionsChange({
        ...options,
        jsonData: { ...options.jsonData, version: 'v4' },
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Backwards compat: a datasource created before @grafana/google-sdk
  // adoption stores the full service-account JSON in `secureJsonData.jwt`.
  // The backend's auth.Resolve still parses that blob, so existing
  // datasources keep authenticating without any user action. New
  // datasources go through <ConnectionConfig /> below and persist the
  // explicit (clientEmail / tokenUri / privateKey) fields the SDK reads.
  const usingLegacyJWT =
    options.secureJsonFields?.jwt &&
    !options.jsonData.clientEmail &&
    !options.secureJsonFields?.privateKey;

  return (
    <div className="gf-form-group">
      <RadioButtonGroup
        options={gaVersion}
        onChange={(v) =>
          onOptionsChange({
            ...options,
            jsonData: { ...jsonData, version: v },
          })
        }
        value={jsonData.version}
      />

      {usingLegacyJWT && (
        <Alert title="Legacy JWT credentials detected" severity="info">
          This datasource was set up before the unified Google authentication
          UI. Existing queries continue to work, but to take advantage of new
          auth options (GCE / Workload Identity / Service Account
          Impersonation) re-upload your service-account JSON below.
        </Alert>
      )}

      <ConnectionConfig {...props} />

      <Alert title="Generate a JWT file" severity="info">
        <ol style={{ listStylePosition: 'inside' }}>
          <li>
            if you don&#39;t have gcp project, add new gcp project.{' '}
            <a href="https://cloud.google.com/resource-manager/docs/creating-managing-projects#console">link</a>
          </li>
          <li>
            Open the <a href="https://console.developers.google.com/apis/credentials">Credentials</a> page in the
            Google API Console.
          </li>
          <li>
            Click <strong>Create Credentials</strong> then click <strong>Service account</strong>.
          </li>
          <li>On the Create service account page, enter the Service account details.</li>
          <li>
            On the <code>Create service account</code> page, fill in the <code>Service account details</code> and then
            click <code>Create</code>
          </li>
          <li>
            On the <code>Service account permissions</code> page, don&#39;t add a role to the service account. Just
            click <code>Continue</code>
          </li>
          <li>
            In the next step, click <code>Create Key</code>. Choose key type <code>JSON</code> and click
            <code>Create</code>. A JSON key file will be created and downloaded to your computer
          </li>
          <li>
            Note your <code>service account email</code> ex) *<em>@</em>*.iam.gserviceaccount.com
          </li>
          <li>
            Open the{' '}
            {jsonData.version === 'v3' ? (
              <a href="https://console.cloud.google.com/apis/library/analytics.googleapis.com">
                Google Analytics API(UA)
              </a>
            ) : (
              <a href="https://console.cloud.google.com/apis/library/analyticsadmin.googleapis.com">
                Google Analytics Admin API(GA4)
              </a>
            )}{' '}
            in API Library and enable access for your account
          </li>
          <li>
            Open the{' '}
            {jsonData.version === 'v3' ? (
              <a href="https://console.cloud.google.com/marketplace/product/google/analyticsreporting.googleapis.com">
                Google Analytics Reporting API(UA)
              </a>
            ) : (
              <a href="https://console.cloud.google.com/apis/library/analyticsdata.googleapis.com">
                Google Analytics Data API(GA4)
              </a>
            )}{' '}
            in API Library and enable access for your Analytics Data
          </li>
          <li>
            <a href="https://console.cloud.google.com/apis/dashboard">Check your api setting</a>
          </li>
        </ol>

        <h3 id="google-analytics-setting">Google Analytics Setting</h3>
        <ol style={{ listStylePosition: 'inside' }}>
          <li>
            Open the <a href="https://analytics.google.com/">Google Analytics</a>
          </li>
          <li>Select Your Analytics Account And Open Admin Page</li>
          <li>
            Click <strong>Account User Management</strong> on the <strong>Account Tab</strong>
          </li>
          <li>Click plus Button then Add users</li>
          <li>
            Enter <code>service account email</code> at <strong>Generate a JWT file</strong> 8th step and Permissions
            add <code>Read &amp; Analyze</code>
          </li>
        </ol>
      </Alert>
    </div>
  );
};
