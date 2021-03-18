package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// DatasourceSettings contains Google Sheets API authentication properties.
type DatasourceSettings struct {
	JWT       string `json:"jwt"`
	ProfileId string `json:"profileId"`
}

// LoadSettings gets the relevant settings from the plugin context
func LoadSettings(ctx backend.PluginContext) (*DatasourceSettings, error) {
	model := &DatasourceSettings{}

	settings := ctx.DataSourceInstanceSettings
	err := json.Unmarshal(settings.JSONData, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading settings: %s", err.Error())
	}

	model.JWT = settings.DecryptedSecureJSONData["jwt"]
	model.ProfileId = settings.DecryptedSecureJSONData["profileId"]

	return model, nil
}
