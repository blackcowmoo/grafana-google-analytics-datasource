package gav3

import (
	"context"
	"fmt"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/patrickmn/go-cache"
)

// GoogleAnalyticsv3DataSource handler for google sheets
type GoogleAnalyticsv3 struct {
	Cache *cache.Cache
}

func (ga *GoogleAnalyticsv3) Query(ctx context.Context, config *setting.DatasourceSecretSettings, query backend.DataQuery) (*data.Frames, error) {
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

func (ga *GoogleAnalyticsv3) GetAccounts(ctx context.Context, config *setting.DatasourceSecretSettings) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:accounts:%s", config.JWT)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(map[string]string), nil
	}

	accounts, err := client.getAccountsList(GaDefaultIdx)
	if err != nil {
		return nil, err
	}

	accountNames := map[string]string{}
	for _, i := range accounts {
		accountNames[i.Id] = i.Name
	}

	ga.Cache.Set(cacheKey, accountNames, 60*time.Second)
	return accountNames, nil
}

func (ga *GoogleAnalyticsv3) GetWebProperties(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:account:%s:webproperties", accountId)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(map[string]string), nil
	}

	Webproperties, err := client.getWebpropertiesList(accountId, GaDefaultIdx)
	if err != nil {
		return nil, err
	}

	WebpropertyNames := map[string]string{}
	for _, i := range Webproperties {
		WebpropertyNames[i.Id] = i.Name
	}

	ga.Cache.Set(cacheKey, WebpropertyNames, 60*time.Second)
	return WebpropertyNames, nil
}

func (ga *GoogleAnalyticsv3) GetProfiles(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config.JWT)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	cacheKey := fmt.Sprintf("analytics:account:%s:webproperty:%s:profiles", accountId, webPropertyId)
	if item, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return item.(map[string]string), nil
	}

	profiles, err := client.getProfilesList(accountId, webPropertyId, GaDefaultIdx)
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

func (ga *GoogleAnalyticsv3) GetProfileTimezone(ctx context.Context, config *setting.DatasourceSecretSettings, accountId string, webPropertyId string, profileId string) (string, error) {
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

func (ga *GoogleAnalyticsv3) GetAllProfilesList(ctx context.Context, config *setting.DatasourceSecretSettings) (map[string]string, error) {
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
