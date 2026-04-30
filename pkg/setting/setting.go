package setting

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type DatasourceSettings struct {
	Version string `json:"version"`
}

// DatasourceSecretSettings contains Google Analytics datasource auth properties.
// JSON-encoded fields come from `jsonData`; secret fields are pulled from
// `DecryptedSecureJSONData` in LoadSettings.
type DatasourceSecretSettings struct {
	// jsonData
	Version            string `json:"version"`
	AuthenticationType string `json:"authenticationType"`
	ClientEmail        string `json:"clientEmail"`
	TokenURI           string `json:"tokenUri"`
	DefaultProject     string `json:"defaultProject"`

	// secureJsonData
	JWT        string `json:"jwt"`        // legacy: full service-account JSON blob
	PrivateKey string `json:"privateKey"` // new: just the PEM private key
	ProfileId  string `json:"profileId"`
}

// LoadSettings gets the relevant settings from the plugin context
func LoadSettings(ctx backend.PluginContext) (*DatasourceSecretSettings, error) {
	model := &DatasourceSecretSettings{}

	settings := ctx.DataSourceInstanceSettings
	err := json.Unmarshal(settings.JSONData, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	model.JWT = settings.DecryptedSecureJSONData["jwt"]
	model.PrivateKey = settings.DecryptedSecureJSONData["privateKey"]

	return model, nil
}
