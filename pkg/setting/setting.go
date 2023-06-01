package setting

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)
type DatasourceSettings struct {
	Version string `json:"version"`
}


// DatasourceSecretSettings contains Google Sheets API authentication properties.
type DatasourceSecretSettings struct {
	JWT       string `json:"jwt"`
	ProfileId string `json:"profileId"`
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

	return model, nil
}
