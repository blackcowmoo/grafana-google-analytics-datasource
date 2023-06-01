package main

import (
	"context"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"setting"
)

type GoogleAnalytics interface {
	Query(context.Context, *setting.DatasourceSecretSettings, backend.DataQuery) (*data.Frames, error)
	GetAccounts(context.Context, *setting.DatasourceSecretSettings) (map[string]string, error)
	GetWebProperties(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string) (map[string]string, error)
	GetProfiles(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string) (map[string]string, error)
	GetProfileTimezone(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string, profileId string) (string, error)
	GetAllProfilesList(ctx context.Context, config *setting.DatasourceSecretSettings) (map[string]string, error)
	GetDimensions() ([]interface{}, error)
	GetMetrics() ([]interface{}, error)
}
