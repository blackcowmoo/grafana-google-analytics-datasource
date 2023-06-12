package gav3

import (
	"context"
	"fmt"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/util"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

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

func (client *GoogleClient) getProfile(accountId, webpropertyId, profileId string) (*analytics.Profile, error) {
	profile, err := client.analytics.Management.Profiles.Get(accountId, webpropertyId, profileId).Do()
	if err != nil {
		log.DefaultLogger.Error("getProfile fail", "error", err.Error(), "accountId", accountId, "webpropertyId", webpropertyId)
		return nil, err
	}

	return profile, nil
}

func (client *GoogleClient) getReport(query model.QueryModel) (*reporting.GetReportsResponse, error) {
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

	if report.Reports[0].NextPageToken != "" {
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

func (client *GoogleClient) getAccountSummaries(start int64) ([]*analytics.AccountSummary, error) {
	accountSummaries, err := client.analytics.Management.AccountSummaries.List().MaxResults(GaManageMaxResult).StartIndex(start).Do()
	if err != nil {
		log.DefaultLogger.Error("getAccountSummary fail", "error", err.Error())
		return nil, err
	}

	if accountSummaries.TotalResults > (start + GaManageMaxResult - 1) {
		start += GaManageMaxResult
		nextAccountSummaries, err := client.getAccountSummaries(start)
		if err != nil {
			return nil, err
		}
		accountSummaries.Items = append(accountSummaries.Items, nextAccountSummaries...)
	}

	return accountSummaries.Items, nil
}
