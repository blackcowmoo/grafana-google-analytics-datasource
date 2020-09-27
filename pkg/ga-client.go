package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	analytics "google.golang.org/api/analytics/v3"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

type QueryData struct {
	ViewID    string `json:"viewId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Metric    string `json:"metric"`
	Dimension string `json:"dimension"`
}

type GoogleClient struct {
	reporting *reporting.Service
	analytics *analytics.Service
}

func NewGoogleClient(ctx context.Context, auth *DatasourceSettings) (*GoogleClient, error) {
	reportingService, reportingError := createReportingService(ctx, auth)
	if reportingError != nil {
		return nil, reportingError
	}

	analyticsService, analyticsError := createAnalticsService(ctx, auth)
	if analyticsError != nil {
		return nil, analyticsError
	}

	return &GoogleClient{reportingService, analyticsService}, nil
}

func createReportingService(ctx context.Context, auth *DatasourceSettings) (*reporting.Service, error) {
	if len(auth.AuthType) == 0 {
		return nil, fmt.Errorf("missing AuthType setting")
	}

	if auth.AuthType == "key" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("missing API Key")
		}
		return reporting.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	if auth.AuthType == "jwt" {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), reporting.AnalyticsReadonlyScope)
		if err != nil {
			return nil, fmt.Errorf("error parsing JWT file: %w", err)
		}

		client := jwtConfig.Client(ctx)
		return reporting.NewService(ctx, option.WithHTTPClient(client))
	}

	return nil, fmt.Errorf("invalid Auth Type: %s", auth.AuthType)
}

func createAnalticsService(ctx context.Context, auth *DatasourceSettings) (*analytics.Service, error) {
	if len(auth.AuthType) == 0 {
		return nil, fmt.Errorf("missing AuthType setting")
	}

	if auth.AuthType == "key" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("missing API Key")
		}
		return analytics.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	if auth.AuthType == "jwt" {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), analytics.AnalyticsReadonlyScope)
		if err != nil {
			return nil, fmt.Errorf("error parsing JWT file: %w", err)
		}

		client := jwtConfig.Client(ctx)
		return analytics.NewService(ctx, option.WithHTTPClient(client))
	}

	return nil, fmt.Errorf("invalid Auth Type: %s", auth.AuthType)

}

func (client *GoogleClient) getAccountsList() ([]*analytics.Account, error) {
	accountsService := analytics.NewManagementAccountsService(client.analytics)
	accounts, err := accountsService.List().Do()
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	return accounts.Items, nil
}

func (client *GoogleClient) getAllWebpropertiesList() ([]*analytics.Webproperty, error) {
	accounts, err := client.getAccountsList()
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	var webpropertiesList = make([]*analytics.Webproperty, 0)
	for _, account := range accounts {
		webproperties, err := client.getWebpropertiesList(account.Id)
		if err != nil {
			log.DefaultLogger.Error(err.Error())
			return nil, err
		}

		webpropertiesList = append(webpropertiesList, webproperties...)
	}

	return webpropertiesList, nil
}

func (client *GoogleClient) getWebpropertiesList(accountId string) ([]*analytics.Webproperty, error) {
	webpropertiesService := analytics.NewManagementWebpropertiesService(client.analytics)
	webproperties, err := webpropertiesService.List(accountId).Do()
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}
	return webproperties.Items, nil
}

func (client *GoogleClient) getAllProfilesList() ([]*analytics.Profile, error) {
	webproperties, err := client.getAllWebpropertiesList()
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	var profilesList = make(chan *analytics.Profile)
	var wait sync.WaitGroup

	var MAX_RETRY_COUNT = 10

	for _, webproperty := range webproperties {
		wait.Add(1)
		go func(accountId string, webpropertyId string) {
			for i := 1; i <= MAX_RETRY_COUNT; i++ {
				log.DefaultLogger.Info("getProfilesList", "accountId", accountId, "webpropertyId", webpropertyId, "retry", i)
				profiles, err := client.getProfilesList(accountId, webpropertyId)
				if err != nil {
					if i < MAX_RETRY_COUNT {
						time.Sleep(time.Second * 1)
					} else {
						log.DefaultLogger.Error(err.Error())
						panic(err)
						wait.Done()
					}
				} else {
					for _, profile := range profiles {
						profilesList <- profile
					}
					i = MAX_RETRY_COUNT
					wait.Done()
				}
			}
			log.DefaultLogger.Info("getProfilesList:End")
		}(webproperty.AccountId, webproperty.Id)
	}
	wait.Wait()
	close(profilesList)

	var profiles = make([]*analytics.Profile, 0)
	for profile := range profilesList {
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func (client *GoogleClient) getProfilesList(accountId string, webpropertyId string) ([]*analytics.Profile, error) {
	profilesService := analytics.NewManagementProfilesService(client.analytics)
	profiles, err := profilesService.List(accountId, webpropertyId).Do()
	if err != nil {

		log.DefaultLogger.Error(err.Error())
		return nil, err
	}

	return profiles.Items, nil
}

func (client *GoogleClient) getReport(query QueryData) (*reporting.GetReportsResponse, error) {
	// A GetReportsRequest instance is a batch request
	// which can have a maximum of 5 requests
	req := &reporting.GetReportsRequest{
		// Our request contains only one request
		// So initialise the slice with one ga.ReportRequest object
		ReportRequests: []*reporting.ReportRequest{
			// Create the ReportRequest object.
			{
				ViewId: query.ViewID,
				DateRanges: []*reporting.DateRange{
					// Create the DateRange object.
					{StartDate: query.StartDate, EndDate: query.EndDate},
				},
				Metrics: []*reporting.Metric{
					// Create the Metrics object.
					{Expression: query.Metric},
				},
				Dimensions: []*reporting.Dimension{
					{Name: query.Dimension},
				},
			},
		},
	}

	log.DefaultLogger.Info("Doing GET request from analytics reporting", req)
	// Call the BatchGet method and return the response.
	return client.reporting.Reports.BatchGet(req).Do()
}

func printResponse(res *reporting.GetReportsResponse) {
	log.DefaultLogger.Info("Printing Response from analytics reporting", "")
	for _, report := range res.Reports {
		header := report.ColumnHeader
		dimHdrs := header.Dimensions
		metricHdrs := header.MetricHeader.MetricHeaderEntries
		rows := report.Data.Rows

		if rows == nil {
			log.DefaultLogger.Info("no data", "")
		}
		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			for i := 0; i < len(dimHdrs) && i < len(dims); i++ {
				log.DefaultLogger.Info("%s: %s", dimHdrs[i], dims[i])
			}

			for _, metric := range metrics {
				// We have only 1 date range in the example
				// So it'll always print "Date Range (0)"
				// log.DefaultLogger.Defaultlog.DefaultLogger.Infof("Date Range (%d)", idx)
				for j := 0; j < len(metricHdrs) && j < len(metric.Values); j++ {
					log.DefaultLogger.Info("%s: %s", metricHdrs[j].Name, metric.Values[j])
				}
			}
		}
	}
	log.DefaultLogger.Info("Completed printing response", "", "")
}
