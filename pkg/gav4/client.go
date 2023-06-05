package gav4

import (
	"context"
	"fmt"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/util"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	analyticsadmin "google.golang.org/api/analyticsadmin/v1beta"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

type GoogleClient struct {
	analyticsdata  *analyticsdata.Service
	analyticsadmin *analyticsadmin.Service
}

func NewGoogleClient(ctx context.Context, jwt string) (*GoogleClient, error) {
	analyticsdataService, analyticsdataError := createAnalyticsdataService(ctx, jwt)
	if analyticsdataError != nil {
		return nil, analyticsdataError
	}

	analyticsadminService, analyticsadminError := createAnalyticsadminService(ctx, jwt)
	if analyticsadminError != nil {
		return nil, analyticsadminError
	}

	return &GoogleClient{analyticsdataService, analyticsadminService}, nil
}

func createAnalyticsdataService(ctx context.Context, jwt string) (*analyticsdata.Service, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(jwt), analyticsdata.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT file: %w", err)
	}

	client := jwtConfig.Client(ctx)
	return analyticsdata.NewService(ctx, option.WithHTTPClient(client))
}

func createAnalyticsadminService(ctx context.Context, jwt string) (*analyticsadmin.Service, error) {
	jwtConfig, err := google.JWTConfigFromJSON([]byte(jwt), analyticsadmin.AnalyticsReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("error parsing JWT file: %w", err)
	}

	client := jwtConfig.Client(ctx)
	return analyticsadmin.NewService(ctx, option.WithHTTPClient(client))
}

func (client *GoogleClient) getAccountsList(nextPageToekn string) ([]*analyticsadmin.GoogleAnalyticsAdminV1betaAccount, error) {
	accountsService := analyticsadmin.NewAccountsService(client.analyticsadmin)
	accounts, err := accountsService.List().PageSize(GaAdminMaxResult).PageToken(nextPageToekn).Do()
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

func (client *GoogleClient) getAllWebpropertiesList() ([]*analyticsadmin.GoogleAnalyticsAdminV1betaProperty, error) {
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

func (client *GoogleClient) getWebpropertiesList(accountId string, nextPageToekn string) ([]*analyticsadmin.GoogleAnalyticsAdminV1betaProperty, error) {
	webpropertiesService := analyticsadmin.NewPropertiesService(client.analyticsadmin)
	webproperties, err := webpropertiesService.List().Filter("parent:" + accountId).PageSize(GaAdminMaxResult).PageToken(nextPageToekn).Do()
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

func (client *GoogleClient) getReport(query QueryModel) (*analyticsdata.RunReportResponse, error) {
	defer util.Elapsed("Get report data at GA API")()
	log.DefaultLogger.Debug("getReport", "queries", query)
	Metrics := []*analyticsdata.Metric{}
	Dimensions := []*analyticsdata.Dimension{}
	for _, metric := range query.Metrics {
		Metrics = append(Metrics, &analyticsdata.Metric{Name: metric})
	}
	for _, dimension := range query.Dimensions {
		Dimensions = append(Dimensions, &analyticsdata.Dimension{Name: dimension})
	}
	var offset int64 = 0
	req := analyticsdata.RunReportRequest{
		DateRanges: []*analyticsdata.DateRange{
			// Create the DateRange object.
			{StartDate: query.StartDate, EndDate: query.EndDate},
		},
		Metrics:    Metrics,
		Dimensions: Dimensions,
		Offset: offset,
		OrderBys: []*analyticsdata.OrderBy{
			{
				Dimension: &analyticsdata.DimensionOrderBy{
					DimensionName: query.TimeDimension,
				},
			},
		},
	}

	log.DefaultLogger.Debug("Doing GET request from analytics reporting", "req", req)
	// Call the BatchGet method and return the response.
	report, err := client.analyticsdata.Properties.RunReport(query.WebPropertyID, &req).Do()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	//  TODO 페이지 네이션
	log.DefaultLogger.Debug("Do GET report", "report len", report.RowCount, "report", report)

	if query.UseNextPage && report.RowCount > (query.Offset + GaAdminMaxResult) {
		query.Offset = query.Offset + GaAdminMaxResult
		newReport, err := client.getReport(query)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		report.Rows = append(report.Rows, newReport.Rows...)
		return report, nil
	}
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

func (client *GoogleClient) getMetadata(propertyID string) (*analyticsdata.Metadata, error) {
	if propertyID == "" {
		propertyID = "0"
	}
	nameid := "properties/" + propertyID + "/metadata"
	metadata, err := client.analyticsdata.Properties.GetMetadata(nameid).Do()
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
