package gav3

import (
	"context"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"sync"
	"time"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/util"

	analytics "google.golang.org/api/analytics/v3"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

type GoogleClient struct {
	reporting *reporting.Service
	analytics *analytics.Service
}

func NewGoogleClient(ctx context.Context, jwt string) (*GoogleClient, error) {
	reportingService, reportingError := createReportingService(ctx, jwt)
	if reportingError != nil {
		return nil, reportingError
	}

	analyticsService, analyticsError := createAnalyticsService(ctx, jwt)
	if analyticsError != nil {
		return nil, analyticsError
	}

	return &GoogleClient{reportingService, analyticsService}, nil
}

func createReportingService(ctx context.Context, jwt string) (*reporting.Service, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(jwt), reporting.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT file: %w", err)
	}

	client := jwtConfig.Client(ctx)
	return reporting.NewService(ctx, option.WithHTTPClient(client))
}

func createAnalyticsService(ctx context.Context, jwt string) (*analytics.Service, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(jwt), analytics.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT file: %w", err)
	}

	client := jwtConfig.Client(ctx)
	return analytics.NewService(ctx, option.WithHTTPClient(client))
}

func (client *GoogleClient) getAccountsList(idx int64) ([]*analytics.Account, error) {
	accountsService := analytics.NewManagementAccountsService(client.analytics)
	accounts, err := accountsService.List().StartIndex(idx).MaxResults(GaManageMaxResult).Do()
	if err != nil {
		log.DefaultLogger.Error("getAccountsList Fail", "error", err.Error())
		return nil, err
	}

	itemPerPage := accounts.ItemsPerPage
	nextLink := accounts.NextLink
	startIdx := accounts.StartIndex

	if nextLink != "" {
		newAccounts, err := client.getAccountsList(startIdx + itemPerPage)
		if err != nil {
			return nil, err
		}
		accounts.Items = append(accounts.Items, newAccounts...)
	}

	return accounts.Items, nil
}

func (client *GoogleClient) getAllWebpropertiesList() ([]*analytics.Webproperty, error) {
	accounts, err := client.getAccountsList(GaDefaultIdx)
	if err != nil {
		log.DefaultLogger.Error("getAllWebpropertiesList fail", "error", err.Error())
		return nil, err
	}

	var webpropertiesList = make([]*analytics.Webproperty, 0)
	for _, account := range accounts {
		webproperties, err := client.getWebpropertiesList(account.Id, GaDefaultIdx)
		if err != nil {
			log.DefaultLogger.Error("getAllWebpropertiesList", "error", err.Error())
			return nil, err
		}

		webpropertiesList = append(webpropertiesList, webproperties...)
	}

	return webpropertiesList, nil
}

func (client *GoogleClient) getWebpropertiesList(accountId string, idx int64) ([]*analytics.Webproperty, error) {
	webpropertiesService := analytics.NewManagementWebpropertiesService(client.analytics)
	webproperties, err := webpropertiesService.List(accountId).StartIndex(idx).MaxResults(GaManageMaxResult).Do()
	if err != nil {
		log.DefaultLogger.Error("getWebpropertiesList fail", "error", err.Error())
		return nil, err
	}

	log.DefaultLogger.Debug("getWebpropertiesList", "WebpropertiesList", webproperties)

	nextLink := webproperties.NextLink
	itemPerPage := webproperties.ItemsPerPage
	startIdx := webproperties.StartIndex

	if nextLink != "" {
		nextWebproperties, err := client.getWebpropertiesList(accountId, startIdx+itemPerPage)
		if err != nil {
			return nil, err
		}
		webproperties.Items = append(webproperties.Items, nextWebproperties...)
	}

	return webproperties.Items, nil
}

func (client *GoogleClient) getAllProfilesList() ([]*analytics.Profile, error) {
	webproperties, err := client.getAllWebpropertiesList()
	if err != nil {
		log.DefaultLogger.Error("getAllProfilesList fail", "error", err.Error())
		return nil, err
	}

	var profilesList = make(chan *analytics.Profile, len(webproperties))
	var wait sync.WaitGroup
	var MAX_RETRY_COUNT = 10

	for _, webproperty := range webproperties {
		wait.Add(1)
		go func(accountId string, webpropertyId string) {
			defer wait.Done()
			for i := 1; i <= MAX_RETRY_COUNT; i++ {
				profiles, err := client.getProfilesList(accountId, webpropertyId, GaDefaultIdx)
				if err != nil {
					if i < MAX_RETRY_COUNT {
						time.Sleep(time.Millisecond * 500)
						continue
					}

					log.DefaultLogger.Error("getProfilesList 10 retries fail", "error", err.Error())
					panic(err)
				}

				for _, profile := range profiles {
					profilesList <- profile
				}

				return
			}
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

func (client *GoogleClient) getProfilesList(accountId string, webpropertyId string, idx int64) ([]*analytics.Profile, error) {
	profilesService := analytics.NewManagementProfilesService(client.analytics)
	profiles, err := profilesService.List(accountId, webpropertyId).MaxResults(GaManageMaxResult).StartIndex(idx).Do()
	if err != nil {
		log.DefaultLogger.Error("getProfilesList fail", "error", err.Error(), "accountId", accountId, "webpropertyId", webpropertyId)
		return nil, err
	}

	nextLink := profiles.NextLink
	itemPerPage := profiles.ItemsPerPage
	startIdx := profiles.StartIndex

	if nextLink != "" {
		nextProfiles, err := client.getProfilesList(accountId, webpropertyId, startIdx+itemPerPage)
		if err != nil {
			return nil, err
		}
		profiles.Items = append(profiles.Items, nextProfiles...)
	}

	return profiles.Items, nil
}

func (client *GoogleClient) getReport(query QueryModel) (*reporting.GetReportsResponse, error) {
	defer util.Elapsed("Get report data at GA API")()
	log.DefaultLogger.Debug("getReport", "queries", query)
	Metrics := []*reporting.Metric{}
	Dimensions := []*reporting.Dimension{}
	for _, metric := range query.Metrics {
		Metrics = append(Metrics, &reporting.Metric{Expression: metric})
	}
	for _, dimension := range query.Dimensions {
		Dimensions = append(Dimensions, &reporting.Dimension{Name: dimension})
	}

	reportRequest := reporting.ReportRequest{
		ViewId: query.ProfileID,
		DateRanges: []*reporting.DateRange{
			// Create the DateRange object.
			{StartDate: query.StartDate, EndDate: query.EndDate},
		},
		Metrics:           Metrics,
		Dimensions:        Dimensions,
		PageSize:          query.PageSize,
		PageToken:         query.PageToken,
		IncludeEmptyRows:  true,
		FiltersExpression: query.FiltersExpression,
	}

	log.DefaultLogger.Debug("getReport", "reportRequests", reportRequest)

	// A GetReportsRequest instance is a batch request
	// which can have a maximum of 5 requests
	req := &reporting.GetReportsRequest{
		// Our request contains only one request
		// So initialise the slice with one ga.ReportRequest object
		ReportRequests: []*reporting.ReportRequest{&reportRequest},
	}

	log.DefaultLogger.Debug("Doing GET request from analytics reporting", "req", req)
	// Call the BatchGet method and return the response.
	report, err := client.reporting.Reports.BatchGet(req).Do()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	log.DefaultLogger.Debug("Do GET report", "report len", len(report.Reports), "report", report)

	if query.UseNextPage && report.Reports[0].NextPageToken != "" {
		query.PageToken = report.Reports[0].NextPageToken
		newReport, err := client.getReport(query)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		report.Reports[0].Data.Rows = append(report.Reports[0].Data.Rows, newReport.Reports[0].Data.Rows...)
		return report, nil
	}
	return report, nil
}

func printResponse(res *reporting.GetReportsResponse) {
	log.DefaultLogger.Debug("Printing Response from analytics reporting", "")
	for _, report := range res.Reports {
		header := report.ColumnHeader
		dimHdrs := header.Dimensions
		metricHdrs := header.MetricHeader.MetricHeaderEntries
		rows := report.Data.Rows

		if rows == nil {
			log.DefaultLogger.Debug("no data", "")
		}
		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			for i := 0; i < len(dimHdrs) && i < len(dims); i++ {
				log.DefaultLogger.Debug("%s: %s", dimHdrs[i], dims[i])
			}

			for _, metric := range metrics {
				// We have only 1 date range in the example
				// So it'll always print "Date Range (0)"
				// log.DefaultLogger.Defaultlog.DefaultLogger.Infof("Date Range (%d)", idx)
				for j := 0; j < len(metricHdrs) && j < len(metric.Values); j++ {
					log.DefaultLogger.Debug("%s: %s", metricHdrs[j].Name, metric.Values[j])
				}
			}
		}
	}
	log.DefaultLogger.Info("Completed printing response", "", "")
}
