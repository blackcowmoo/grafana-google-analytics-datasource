package main

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	analyticsadmin "google.golang.org/api/analyticsadmin/v1beta"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

type GoogleClientv4 struct {
	analyticsdata  *analyticsdata.Service
	analyticsadmin *analyticsadmin.Service
}

func NewGoogleClientv4(ctx context.Context, auth *DatasourceSettings) (*GoogleClientv4, error) {
	analyticsdataService, analyticsdataError := createAnalyticsdataService(ctx, auth)
	if analyticsdataError != nil {
		return nil, analyticsdataError
	}

	analyticsadminService, analyticsadminError := createAnalyticsadminService(ctx, auth)
	if analyticsadminError != nil {
		return nil, analyticsadminError
	}

	return &GoogleClientv4{analyticsdataService, analyticsadminService}, nil
}

func createAnalyticsdataService(ctx context.Context, auth *DatasourceSettings) (*analyticsdata.Service, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), analyticsdata.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT file: %w", err)
	}

	client := jwtConfig.Client(ctx)
	return analyticsdata.NewService(ctx, option.WithHTTPClient(client))
}

func createAnalyticsadminService(ctx context.Context, auth *DatasourceSettings) (*analyticsadmin.Service, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), analyticsadmin.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT file: %w", err)
	}

	client := jwtConfig.Client(ctx)
	return analyticsadmin.NewService(ctx, option.WithHTTPClient(client))
}

func (client *GoogleClientv4) getAccountsList(nextPageToekn string) ([]*analyticsadmin.GoogleAnalyticsAdminV1betaAccount, error) {
	accountsService := analyticsadmin.NewAccountsService(client.analyticsadmin)
	accounts, err := accountsService.List().PageSize(200).PageToken(nextPageToekn).Do()
	if err != nil {
		log.DefaultLogger.Error("getAccountsList Fail", "error", err.Error())
		return nil, err
	}
	nextLink := accounts.NextPageToken

	if nextLink != "" {
		newAccounts, err := client.getAccountsList(nextLink)
		if err != nil {
			return nil, err
		}
		accounts.Accounts = append(accounts.Accounts, newAccounts...)
	}

	return accounts.Accounts, nil
}

func (client *GoogleClientv4) getAllWebpropertiesList() ( []*analyticsadmin.GoogleAnalyticsAdminV1betaProperty, error) {
	accounts, err := client.getAccountsList("")
	if err != nil {
		log.DefaultLogger.Error("getAllWebpropertiesList fail", "error", err.Error())
		return nil, err
	}

	var webpropertiesList = make([]*analyticsadmin.GoogleAnalyticsAdminV1betaProperty, 0)
	for _, account := range accounts {
		webproperties, err := client.getWebpropertiesList(account.Name, "")
		if err != nil {
			log.DefaultLogger.Error("getAllWebpropertiesList", "error", err.Error())
			return nil, err
		}

		webpropertiesList = append(webpropertiesList, webproperties...)
	}

	return webpropertiesList, nil
}

func (client *GoogleClientv4) getWebpropertiesList(accountId string, nextPageToekn string) ( []*analyticsadmin.GoogleAnalyticsAdminV1betaProperty, error) {
	webpropertiesService := analyticsadmin.NewPropertiesService(client.analyticsadmin)
	webproperties, err := webpropertiesService.List().Filter("parent:"+accountId).PageSize(200).PageToken(nextPageToekn).Do()
	if err != nil {
		log.DefaultLogger.Error("getWebpropertiesList fail", "error", err.Error())
		return nil, err
	}

	log.DefaultLogger.Debug("getWebpropertiesList", "WebpropertiesList", webproperties)

	nextLink := webproperties.NextPageToken

	if nextLink != "" {
		nextWebproperties, err := client.getWebpropertiesList(accountId, nextLink)
		if err != nil {
			return nil, err
		}
		webproperties.Properties = append(webproperties.Properties, nextWebproperties...)
	}

	return webproperties.Properties, nil
}

func (client *GoogleClientv4) getReport(query QueryModel) (*analyticsdata.RunReportResponse, error) {
	defer Elapsed("Get report data at GA API")()
	log.DefaultLogger.Debug("getReport", "queries", query)
	Metrics := []*analyticsdata.Metric{}
	Dimensions := []*analyticsdata.Dimension{}
	for _, metric := range query.Metrics {
		Metrics = append(Metrics, &analyticsdata.Metric{Expression: metric})
	}
	for _, dimension := range query.Dimensions {
		Dimensions = append(Dimensions, &analyticsdata.Dimension{Name: dimension})
	}

	req := analyticsdata.RunReportRequest{
		DateRanges: []*analyticsdata.DateRange{
			// Create the DateRange object.
			{StartDate: query.StartDate, EndDate: query.EndDate},
		},
		Metrics:           Metrics,
		Dimensions:        Dimensions,
	}

	log.DefaultLogger.Debug("Doing GET request from analytics reporting", "req", req)
	// Call the BatchGet method and return the response.
	report, err := client.analyticsdata.Properties.RunReport(query.WebPropertyID,&req).Do()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
//  TODO 페이지 네이션
	// log.DefaultLogger.Debug("Do GET report", "report len", report.RowCount, "report", report)

	// if query.UseNextPage && report.Reports[0].NextPageToken != "" {
	// 	query.PageToken = report.Reports[0].NextPageToken
	// 	newReport, err := client.getReport(query)
	// 	if err != nil {
	// 		return nil, fmt.Errorf(err.Error())
	// 	}

	// 	report.Reports[0].Data.Rows = append(report.Reports[0].Data.Rows, newReport.Reports[0].Data.Rows...)
	// 	return report, nil
	// }
	return report, nil
}

// func printResponse(res *reporting.GetReportsResponse) {
// 	log.DefaultLogger.Debug("Printing Response from analytics reporting", "")
// 	for _, report := range res.Reports {
// 		header := report.ColumnHeader
// 		dimHdrs := header.Dimensions
// 		metricHdrs := header.MetricHeader.MetricHeaderEntries
// 		rows := report.Data.Rows

// 		if rows == nil {
// 			log.DefaultLogger.Debug("no data", "")
// 		}
// 		for _, row := range rows {
// 			dims := row.Dimensions
// 			metrics := row.Metrics

// 			for i := 0; i < len(dimHdrs) && i < len(dims); i++ {
// 				log.DefaultLogger.Debug("%s: %s", dimHdrs[i], dims[i])
// 			}

// 			for _, metric := range metrics {
// 				// We have only 1 date range in the example
// 				// So it'll always print "Date Range (0)"
// 				// log.DefaultLogger.Defaultlog.DefaultLogger.Infof("Date Range (%d)", idx)
// 				for j := 0; j < len(metricHdrs) && j < len(metric.Values); j++ {
// 					log.DefaultLogger.Debug("%s: %s", metricHdrs[j].Name, metric.Values[j])
// 				}
// 			}
// 		}
// 	}
// 	log.DefaultLogger.Info("Completed printing response", "", "")
// }
