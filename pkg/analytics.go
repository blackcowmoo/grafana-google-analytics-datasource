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
	GetWebProperties(context.Context, *setting.DatasourceSecretSettings, string) (map[string]string, error)
	GetProfiles(context.Context, *setting.DatasourceSecretSettings, string, string) (map[string]string, error)
	GetTimezone(context.Context, *setting.DatasourceSecretSettings, string, string, string) (string, error)
	GetAllProfilesList(context.Context, *setting.DatasourceSecretSettings) (map[string]string, error)
	GetDimensions(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
	GetMetrics(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
	CheckHealth(context.Context, *setting.DatasourceSecretSettings) (*backend.CheckHealthResult, error)
}
