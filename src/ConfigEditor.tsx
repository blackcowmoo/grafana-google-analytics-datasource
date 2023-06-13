import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { RadioButtonGroup } from '@grafana/ui';
import React, { PureComponent } from 'react';
import { GADataSourceOptions, GASecureJsonData } from 'types';
import { JWTConfig } from './JWTConfig';

const gaVersion = [
  { label: 'UA(GA3)', value: 'v3' },
  { label: 'GA4(beta)', value: 'v4' },
];

export type Props = DataSourcePluginOptionsEditorProps<GADataSourceOptions, GASecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  constructor(props: Readonly<Props>) {
    super(props)
    if (!this.props.options.jsonData.version) { this.props.options.jsonData.version = 'v3' }
  }
  onResetProfileId = () => {
    const { options } = this.props;
    this.props.onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
      },
      secureJsonFields: {
        ...options.secureJsonFields,
      },
    });
  };

  render() {
    const { options, onOptionsChange } = this.props;
    const { secureJsonFields } = options;
    const secureJsonData = options.secureJsonData;
    const jsonData = options.jsonData
    return (
      <div className="gf-form-group">
        <RadioButtonGroup
          options={gaVersion}
          onChange={(v) => onOptionsChange({
            ...options,
            jsonData: {
              ...jsonData,
              version: v
            }
          })
          }
          value={options.jsonData.version}
        />
        <>
          <JWTConfig
            isConfigured={(secureJsonFields && !!secureJsonFields.jwt) as boolean}
            onChange={(jwt) => {
              onOptionsChange({
                ...options,
                secureJsonData: {
                  ...secureJsonData,
                  jwt,
                },
              });
            }}
          ></JWTConfig>
        </>
        <div className="grafana-info-box" style={{ marginTop: 24 }}>
          <h3 id="generate-a-jwt-file">Generate a JWT file</h3>
          <ol style={{ listStylePosition: 'inside' }}>
            <li>
              if you don&#39;t have gcp project, add new gcp project.
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
              Open the
              {
                jsonData.version === "v3" ?
                  <>
                    <a href="https://console.cloud.google.com/apis/library/analytics.googleapis.com">
                      Google Analytics API(UA)
                    </a>
                  </>
                  :
                  <>
                    <a href="https://console.cloud.google.com/apis/library/analyticsadmin.googleapis.com">
                      Google Analytics Admin API(GA4)
                    </a>
                  </>
              }
              in API Library and enable access for your account
            </li>
            <li>
              Open the
              {
                jsonData.version === "v3" ?
                  <>
                    <a href="https://console.cloud.google.com/marketplace/product/google/analyticsreporting.googleapis.com">
                      Google Analytics Reporting API(UA)
                    </a>
                  </>
                  :
                  <>
                    <a href="https://console.cloud.google.com/apis/library/analyticsdata.googleapis.com">
                      Google Analytics Data API(GA4)
                    </a>
                  </>
              }


              in API Library and enable access for your Analytics Data
            </li>
            <li>
              <a href="https://console.cloud.google.com/apis/dashboard">
                Check your api setting
              </a>

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
        </div>
      </div>

    );
  }
}
