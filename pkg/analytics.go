package main

import (
	"context"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// newDatasource returns datasource.ServeOpts.
func newDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &AnalyticsDatasource{
		im: im,
	}

	return datasource.ServeOpts{
		// QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

type AnalyticsDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (td *AnalyticsDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {

	var status = backend.HealthStatusOk
	var message = "Success"

	config, err := LoadSettings(req.PluginContext)

	log.DefaultLogger.Info("LoadSetting", config.ViewID)

	if err != nil {
		log.DefaultLogger.Info("Fail LoadSetting", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Setting Configuration Read Fail",
		}, nil
	}

	client, err := NewGoogleClient(ctx, config)

	if err != nil {
		log.DefaultLogger.Info("Fail NewGoogleClient", err.Error())

		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Invalid config",
		}, nil
	}

	testData := QueryData{config.ViewID, "yesterday", "today", "ga:sessions", "ga:country"}
	res, err := getReport(client, testData)

	if err != nil {
		log.DefaultLogger.Info("GET request to analyticsreporting/v4 returned error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Test Request Fail",
		}, nil
	}

	if res != nil {
		log.DefaultLogger.Info("HTTPStatusCode", res.HTTPStatusCode)
		log.DefaultLogger.Info("res", res)
	}

	printResponse(res)

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

type instanceSettings struct {
	httpClient *http.Client
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &instanceSettings{
		httpClient: &http.Client{},
	}, nil
}
