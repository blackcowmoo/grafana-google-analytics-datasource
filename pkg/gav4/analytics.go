package gav4

import (
	"context"
	"fmt"
	"time"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"

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

	if len(queryModel.AccountID) < 1 {
		log.DefaultLogger.Error("Query", "error", "Required AccountID")
		return nil, fmt.Errorf("Required AccountID")
	}

	if len(queryModel.WebPropertyID) < 1 {
		log.DefaultLogger.Error("Query", "error", "Required WebPropertyID")
		return nil, fmt.Errorf("Required WebPropertyID")
	}

	report, err := client.getReport(*queryModel)
	if err != nil {
		log.DefaultLogger.Error("Query", "error", err)
		return nil, err
	}

	return transformReportsResponseToDataFrames(report, queryModel.RefID, queryModel.Timezone)

}

func (ga *GoogleAnalytics) GetAccounts(ctx context.Context, config *setting.DatasourceSecretSettings) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:accounts:%s", config.JWT)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(map[string]string), nil
	}

	accounts, err := client.getAccountsList("")
	if err != nil {
		return nil, err
	}

	accountNames := map[string]string{}
	for _, i := range accounts {
		accountNames[i.Name] = i.DisplayName
	}

	ga.Cache.Set(cacheKey, accountNames, 60*time.Second)
	return accountNames, nil
}

func (ga *GoogleAnalytics) GetWebProperties(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:account:%s:webproperties", accountId)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(map[string]string), nil
	}

	Webproperties, err := client.getWebpropertiesList(accountId, "")
	if err != nil {
		return nil, err
	}

	WebpropertyNames := map[string]string{}
	for _, i := range Webproperties {
		WebpropertyNames[i.Name] = i.DisplayName
	}

	ga.Cache.Set(cacheKey, WebpropertyNames, 60*time.Second)
	return WebpropertyNames, nil
}

func (ga *GoogleAnalytics) GetProfiles(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string) (map[string]string, error) {
	// gav4 no profle
	return nil, nil
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

func (ga *GoogleAnalytics) GetAllProfilesList(ctx context.Context, config *setting.DatasourceSecretSettings) (map[string]string, error) {
	return nil, fmt.Errorf("ga4 have not profile")
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
	var dimensions = make([]model.MetadataItem, 0)
	var metrics = make([]model.MetadataItem, 0)
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

func (ga *GoogleAnalytics) CheckHealth(ctx context.Context, config *setting.DatasourceSecretSettings) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Success"

	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		log.DefaultLogger.Error("CheckHealth: Fail NewGoogleClient", "error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Fail NewGoogleClient" + err.Error(),
		}, nil
	}

	webProperties, err := client.getAllWebpropertiesList()
	if err != nil {
		log.DefaultLogger.Error("CheckHealth: Fail getPropetyList", "error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Fail getPropetyList" + err.Error(),
		}, nil
	}

	testData := QueryModel{webProperties[0].Account, webProperties[0].Name, "", "2daysAgo", "today", "a", []string{"active1DayUsers"}, "date", []string{}, 1, "", false, "UTC", "", 0}
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

	var accounts []*model.AccountSummary
	for _, rawAccountSummary := range rawAccountSummaries {
		var account = &model.AccountSummary{
			Account:     rawAccountSummary.Account,
			DisplayName: rawAccountSummary.DisplayName,
		}
		var propertySummaries = make([]*model.PropertySummary, 0)
		for _, rawPpropertySummary := range rawAccountSummary.PropertySummaries {
			var propertySummary = &model.PropertySummary{
				Property:     rawPpropertySummary.Property,
				DisplayName: rawPpropertySummary.DisplayName,
				Parent:      rawPpropertySummary.DisplayName,
			}
			propertySummaries = append(propertySummaries, propertySummary)
		}
		account.PropertySummaries = propertySummaries
		accounts = append(accounts, account)
	}
	ga.Cache.Set(cacheKey, accounts, 60*time.Second)
	return accounts, nil
}
