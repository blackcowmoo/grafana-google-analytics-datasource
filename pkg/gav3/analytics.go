package gav3

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

// GoogleAnalyticsv3DataSource handler for google sheets
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

	if len(queryModel.ProfileID) < 1 {
		log.DefaultLogger.Error("Query", "error", "Required ProfileID")
		return nil, fmt.Errorf("Required ProfileID")
	}

	report, err := client.getReport(*queryModel)
	if err != nil {
		log.DefaultLogger.Error("Query", "error", err)
		return nil, err
	}

	return transformReportsResponseToDataFrames(report, queryModel.RefID, queryModel.Timezone)
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

	profiles, err := client.getProfilesList(accountId, webPropertyId, GaDefaultIdx)
	if err != nil {
		return "", err
	}

	var timezone string
	for _, profile := range profiles {
		if profile.Id == profileId {
			timezone = profile.Timezone
			break
		}
	}

	ga.Cache.Set(cacheKey, timezone, 60*time.Second)
	return timezone, nil
}

func (ga *GoogleAnalytics) GetAllProfilesList(ctx context.Context, config *setting.DatasourceSecretSettings) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:account:*:webproperty:*:profiles")
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(map[string]string), nil
	}

	profiles, err := client.getAllProfilesList()
	if err != nil {
		return nil, err
	}

	profileNames := map[string]string{}
	for _, i := range profiles {
		profileNames[i.Id] = i.Name
	}

	ga.Cache.Set(cacheKey, profileNames, 60*time.Second)
	return profileNames, nil
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

	profiles, err := client.getAllProfilesList()
	if err != nil {
		log.DefaultLogger.Error("CheckHealth: Fail getAllProfilesList", "error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Fail getProfileList" + err.Error(),
		}, nil
	}

	testData := model.QueryModel{AccountID: profiles[0].AccountId, WebPropertyID: profiles[0].WebPropertyId, ProfileID: profiles[0].Id, StartDate: "yesterday", EndDate: "today", RefID: "a", Metrics: []string{"ga:sessions"}, TimeDimension: "ga:dateHour", Dimensions: []string{}, PageSize: 1, PageToken: "", UseNextPage: false, Timezone: "UTC", FiltersExpression: "", Offset: GaDefaultIdx}
	res, err := client.getReport(testData)

	if err != nil {
		log.DefaultLogger.Error("CheckHealth: GET request to analyticsreporting/v4 returned error", "error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "CheckHealth: Test Request Fail" + err.Error(),
		}, nil
	}

	if res != nil {
		log.DefaultLogger.Debug("HTTPStatusCode", "status", res.HTTPStatusCode)
		log.DefaultLogger.Debug("res", res)
	}

	printResponse(res)

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

	rawAccountSummaries, err := client.getAccountSummaries(GaDefaultIdx)
	if err != nil {
		return nil, err
	}

	var accounts []*model.AccountSummary
	for _, rawAccountSummary := range rawAccountSummaries {
		var account = &model.AccountSummary{
			Account:     rawAccountSummary.Id,
			DisplayName: rawAccountSummary.Name,
		}
		var propertySummaries = make([]*model.PropertySummary, 0)
		for _, rawPropertySummary := range rawAccountSummary.WebProperties {
			var propertySummary = &model.PropertySummary{
				Property:    rawPropertySummary.Id,
				DisplayName: rawPropertySummary.Name,
				Parent:      rawAccountSummary.Id,
			}
			propertySummaries = append(propertySummaries, propertySummary)

			var profileSummaries = make([]*model.ProfileSummary, 0)

			for _, rawProfileSummaries := range rawPropertySummary.Profiles {
				var profileSummary = &model.ProfileSummary{
					Profile:     rawProfileSummaries.Id,
					DisplayName: rawProfileSummaries.Name,
					Parent:      rawPropertySummary.Id,
					Type:        rawProfileSummaries.Type,
				}
				profileSummaries = append(profileSummaries, profileSummary)
			}
			propertySummary.ProfileSummaries = profileSummaries
		}
		account.PropertySummaries = propertySummaries
		accounts = append(accounts, account)
	}
	ga.Cache.Set(cacheKey, accounts, 60*time.Second)
	return accounts, nil
}
