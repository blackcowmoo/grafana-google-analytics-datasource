package main

import (
	"context"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

type GoogleAnalytics interface {
	Query(context.Context, *setting.DatasourceSecretSettings, backend.DataQuery) (*data.Frames, error)
	GetAccounts(context.Context, *setting.DatasourceSecretSettings) (map[string]string, error)
	GetWebProperties(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string) (map[string]string, error)
	GetProfiles(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string) (map[string]string, error)
	GetTimezone(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string, profileId string) (string, error)
	GetAllProfilesList(ctx context.Context, config *setting.DatasourceSecretSettings) (map[string]string, error)
	GetDimensions(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
	GetMetrics(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
}
