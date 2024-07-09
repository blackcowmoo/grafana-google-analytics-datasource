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
	GetAccountSummaries(context.Context, *setting.DatasourceSecretSettings) ([]*model.AccountSummary, error)
	GetTimezone(context.Context, *setting.DatasourceSecretSettings, string, string, string) (string, error)
  GetServiceLevel(context.Context, *setting.DatasourceSecretSettings, string, string) (string, error)
	GetDimensions(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
	GetRealtimeDimensions(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
	GetRealTimeMetrics(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
	GetMetrics(context.Context, *setting.DatasourceSecretSettings, string) ([]model.MetadataItem, error)
	CheckHealth(context.Context, *setting.DatasourceSecretSettings) (*backend.CheckHealthResult, error)
}
