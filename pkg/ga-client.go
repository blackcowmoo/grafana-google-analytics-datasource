package main

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	ga "google.golang.org/api/analyticsreporting/v4"
)

type QueryData struct {
	ViewID    string `json:"viewId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Metric    string `json:"metric"`
	Dimension string `json:"dimension"`
}

func NewGoogleClient(ctx context.Context, auth *DatasourceSettings) (*ga.Service, error) {
	analyticsreportingService, err := CreateGaService(ctx, auth)
	if err != nil {
		return nil, err
	}

	return analyticsreportingService, nil
}

func CreateGaService(ctx context.Context, auth *DatasourceSettings) (*ga.Service, error) {
	if len(auth.AuthType) == 0 {
		return nil, fmt.Errorf("missing AuthType setting")
	}

	if auth.AuthType == "key" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("missing API Key")
		}
		return ga.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	if auth.AuthType == "jwt" {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), ga.AnalyticsReadonlyScope)
		if err != nil {
			return nil, fmt.Errorf("error parsing JWT file: %w", err)
		}

		client := jwtConfig.Client(ctx)
		return ga.NewService(ctx, option.WithHTTPClient(client))
	}

	return nil, fmt.Errorf("invalid Auth Type: %s", auth.AuthType)
}

func getReport(client *ga.Service, query QueryData) (*ga.GetReportsResponse, error) {
	// A GetReportsRequest instance is a batch request
	// which can have a maximum of 5 requests
	req := &ga.GetReportsRequest{
		// Our request contains only one request
		// So initialise the slice with one ga.ReportRequest object
		ReportRequests: []*ga.ReportRequest{
			// Create the ReportRequest object.
			{
				ViewId: query.ViewID,
				DateRanges: []*ga.DateRange{
					// Create the DateRange object.
					{StartDate: query.StartDate, EndDate: query.EndDate},
				},
				Metrics: []*ga.Metric{
					// Create the Metrics object.
					{Expression: query.Metric},
				},
				Dimensions: []*ga.Dimension{
					{Name: query.Dimension},
				},
			},
		},
	}

	log.DefaultLogger.Info("Doing GET request from analytics reporting", req)
	// Call the BatchGet method and return the response.
	return client.Reports.BatchGet(req).Do()
}

func printResponse(res *ga.GetReportsResponse) {
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
