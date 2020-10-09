package main

import (
	"context"
	"fmt"

	"github.com/patrickmn/go-cache"
)

// GoogleAnalyticsDataSource handler for google sheets
type GoogleAnalytics struct {
	Cache *cache.Cache
}

// GetSpreadsheets gets spreadsheets from the Google API.
func (ga *GoogleAnalytics) GetAccounts(ctx context.Context, config *DatasourceSettings) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	accounts, err := client.getAccountsList()
	if err != nil {
		return nil, err
	}

	accountNames := map[string]string{}
	for _, i := range accounts {
		accountNames[i.Id] = i.Name
	}

	return accountNames, nil
}

// GetSpreadsheets gets spreadsheets from the Google API.
func (ga *GoogleAnalytics) GetWebProperties(ctx context.Context, config *DatasourceSettings, accountId string) (map[string]string, error) {
	client, err := NewGoogleClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google API client: %w", err)
	}

	Webproperties, err := client.getWebpropertiesList(accountId)
	if err != nil {
		return nil, err
	}

	WebpropertyNames := map[string]string{}
	for _, i := range Webproperties {
		WebpropertyNames[i.Id] = i.Name
	}

	return WebpropertyNames, nil
}
