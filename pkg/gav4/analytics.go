package gav4

import (
	"context"
	"fmt"
	"time"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/util"
	"google.golang.org/api/analyticsdata/v1beta"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/patrickmn/go-cache"
)

// GoogleAnalyticsv4DataSource handler
type GoogleAnalytics struct {
	Cache *cache.Cache
}

func (ga *GoogleAnalytics) Query(ctx context.Context, config *setting.DatasourceSecretSettings, query backend.DataQuery) (*data.Frames, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		log.DefaultLogger.Error("Query: Fail NewGoogleClient", "error", err.Error())
		return nil, err
	}
	queryModel, err := GetQueryModel(query)
	if err != nil {
		log.DefaultLogger.Error("Failed to read query: %w", "error", err)
		return nil, fmt.Errorf("failed to read query: %w", err)
	}

	if len(queryModel.WebPropertyID) == 0 {
		log.DefaultLogger.Error("Query", "error", "Required WebPropertyID")
		return nil, fmt.Errorf("required webpropertyid")
	}

	if len(queryModel.Dimensions) == 0 && len(queryModel.Metrics) == 0 {
		log.DefaultLogger.Error("Query", "error", "Required Dimensions or Metrics")
		return nil, fmt.Errorf("required dimensions or metrics")
	}

	if queryModel.Mode == "time series" && len(queryModel.TimeDimension) == 0 {
		log.DefaultLogger.Error("Query", "error", "TimeSeries query need TimeDimension")
		return nil, fmt.Errorf("time series query need time dimensions")
	}

	report, err := ga.getReport(ctx, client, queryModel)
	if err != nil {
		log.DefaultLogger.Error("Query", "error", err)
		return nil, err
	}

	return transformReportsResponseToDataFrames(report, queryModel.RefID, queryModel.Timezone, queryModel.Mode)

}

func (ga *GoogleAnalytics) getReport(ctx context.Context, client *GoogleClient, queryModel *model.QueryModel) (*analyticsdata.RunReportResponse, error) {
	var report *analyticsdata.RunReportResponse
	var err error
	switch queryModel.Mode {
	case model.REALTIME:
		log.DefaultLogger.Debug("Query", "realtime")
		r, err := client.getRealtimeReport(*queryModel)
		if err != nil {
			log.DefaultLogger.Error("Query", "error", err)
			return nil, err
		}
		cvt, err := util.TypeConverter[analyticsdata.RunReportResponse](r)
		if err != nil {
			log.DefaultLogger.Error("Query", "error", err)
			return nil, err
		}
		log.DefaultLogger.Debug("Query", "convert end")
		report = cvt
		log.DefaultLogger.Debug("Query", "realtime end")
	case model.TIME_SERIES, model.TABLE:
		report, err = client.getReport(*queryModel)
		if err != nil {
			log.DefaultLogger.Error("Query", "error", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown query mode")
	}
	return report, nil
}

func (ga *GoogleAnalytics) GetTimezone(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string, profileId string) (string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return "", fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:account:%s:webproperty:%s:profile:%s:timezone", accountId, webPropertyId, profileId)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(string), nil
	}

	webproperty, err := client.GetWebProperty(webPropertyId)
	if err != nil {
		return "", err
	}

	timezone := webproperty.TimeZone

	ga.Cache.Set(cacheKey, timezone, 60*time.Second)
	return timezone, nil
}

func (ga *GoogleAnalytics) GetServiceLevel(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string) (string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return "", fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:account:%s:webproperty:%s:service_level", accountId, webPropertyId)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(string), nil
	}

	webproperty, err := client.GetWebProperty(webPropertyId)
	if err != nil {
		return "", err
	}

	serviceLevel := webproperty.ServiceLevel

	ga.Cache.Set(cacheKey, serviceLevel, 60*time.Second)
	return serviceLevel, nil
}

func (ga *GoogleAnalytics) getFilteredMetadata(ctx context.Context, config *setting.DatasourceSecretSettings, propertyId string) ([]model.MetadataItem, []model.MetadataItem, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Google API client: %w", err)
	}
	metadata, err := client.getMetadata(propertyId)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get metadata: %w", err)
	}
	var dimensions = make([]model.MetadataItem, len(metadata.Dimensions))
	var metrics = make([]model.MetadataItem, len(metadata.Metrics))
	for _, metric := range metadata.Metrics {
		var metadataItem = &model.MetadataItem{}
		metadataItem.ID = metric.ApiName
		metadataItem.Attributes.Description = metric.Description
		metadataItem.Attributes.Group = metric.Category
		metadataItem.Attributes.UIName = metric.UiName
		metrics = append(metrics, *metadataItem)
	}
	for _, dimension := range metadata.Dimensions {
		var metadataItem = &model.MetadataItem{}
		metadataItem.ID = dimension.ApiName
		metadataItem.Attributes.Description = dimension.Description
		metadataItem.Attributes.Group = dimension.Category
		metadataItem.Attributes.UIName = dimension.UiName
		dimensions = append(dimensions, *metadataItem)
	}

	return metrics, dimensions, nil
}

func (ga *GoogleAnalytics) GetDimensions(ctx context.Context, config *setting.DatasourceSecretSettings, propertyId string) ([]model.MetadataItem, error) {
	cacheKey := "ga:metadata:" + propertyId + ":dimensions"
	if dimensions, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return dimensions.([]model.MetadataItem), nil
	}
	_, dimensions, err := ga.getFilteredMetadata(ctx, config, propertyId)
	if err != nil {
		return nil, err
	}

	return dimensions, nil
}

func (ga *GoogleAnalytics) GetMetrics(ctx context.Context, config *setting.DatasourceSecretSettings, propertyId string) ([]model.MetadataItem, error) {
	cacheKey := "ga:metadata:" + propertyId + ":metrics"
	if metrics, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return metrics.([]model.MetadataItem), nil
	}
	metrics, _, err := ga.getFilteredMetadata(ctx, config, propertyId)
	if err != nil {
		return nil, err
	}

	ga.Cache.Set(cacheKey, metrics, time.Hour)

	return metrics, nil
}

func (ga *GoogleAnalytics) GetRealtimeDimensions(ctx context.Context, config *setting.DatasourceSecretSettings, propertyId string) ([]model.MetadataItem, error) {
	cacheKey := "ga:metadata:" + propertyId + ":realtime-dimensions"
	if dimensions, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return dimensions.([]model.MetadataItem), nil
	}

	return GaRealTimeDimensions, nil
}

func (ga *GoogleAnalytics) GetRealTimeMetrics(ctx context.Context, config *setting.DatasourceSecretSettings, propertyId string) ([]model.MetadataItem, error) {
	cacheKey := "ga:metadata:" + propertyId + ":realtime-metrics"
	if metrics, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return metrics.([]model.MetadataItem), nil
	}

	return GaRealTimeMetrics, nil
}

func (ga *GoogleAnalytics) CheckHealth(ctx context.Context, config *setting.DatasourceSecretSettings) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Success"

	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		log.DefaultLogger.Error("CheckHealth: Fail NewGoogleClient", "error", config.JWT)
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Fail NewGoogleClient" + err.Error() + "json:" + config.JWT,
		}, nil
	}

	accountSummaries, err := ga.GetAccountSummaries(ctx, config)
	if err != nil {
		log.DefaultLogger.Error("CheckHealth: Fail getPropetyList", "error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Fail getPropetyList" + err.Error(),
		}, nil
	}
	if len(accountSummaries) == 0 {
		log.DefaultLogger.Error("CheckHealth: Not Exist Valid Proerty")
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Not Exist Valid Proerty",
		}, nil
	}

	testData := model.QueryModel{AccountID: accountSummaries[0].Account, WebPropertyID: accountSummaries[0].PropertySummaries[0].Property, ProfileID: "", StartDate: "yesterday", EndDate: "today", RefID: "a", Metrics: []string{"active1DayUsers"}, TimeDimension: "date", Dimensions: []string{"date"}, PageSize: GaReportMaxResult, PageToken: "", UseNextPage: false, Timezone: "UTC", FiltersExpression: "", Offset: 0}
	res, err := client.getReport(testData)

	if err != nil {
		log.DefaultLogger.Error("CheckHealth: GET request to analyticsdata beta returned error", "error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Test Request Fail" + err.Error(),
		}, nil
	}

	if res != nil {
		log.DefaultLogger.Debug("HTTPStatusCode", "status", res.HTTPStatusCode)
		log.DefaultLogger.Debug("res", res)
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (ga *GoogleAnalytics) GetAccountSummaries(ctx context.Context, config *setting.DatasourceSecretSettings) ([]*model.AccountSummary, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:accountsummaries:%s", config.JWT)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.([]*model.AccountSummary), nil
	}

	rawAccountSummaries, err := client.getAccountSummaries("")
	if err != nil {
		return nil, err
	}
	log.DefaultLogger.Debug("GA4 GetAccountSummaries raw accounts", "debug", rawAccountSummaries)

	var accounts []*model.AccountSummary
	for _, rawAccountSummary := range rawAccountSummaries {
		if len(rawAccountSummary.PropertySummaries) == 0 {
			continue
		}
		var account = &model.AccountSummary{
			Account:     rawAccountSummary.Account,
			DisplayName: rawAccountSummary.DisplayName,
		}
		var propertySummaries = make([]*model.PropertySummary, 0)
		for _, rawPpropertySummary := range rawAccountSummary.PropertySummaries {
			var propertySummary = &model.PropertySummary{
				Property:    rawPpropertySummary.Property,
				DisplayName: rawPpropertySummary.DisplayName,
				Parent:      rawPpropertySummary.DisplayName,
			}
			propertySummaries = append(propertySummaries, propertySummary)
		}
		if len(propertySummaries) > 0 {
			account.PropertySummaries = propertySummaries
			accounts = append(accounts, account)
		}
	}
	log.DefaultLogger.Debug("GA4 GetAccountSummaries parsed accounts", "debug", accounts)
	ga.Cache.Set(cacheKey, accounts, 60*time.Second)
	return accounts, nil
}
