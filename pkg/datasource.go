package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/patrickmn/go-cache"
)

// GoogleAnalyticsDataSource handler for google sheets
type GoogleAnalyticsDataSource struct {
	analytics *GoogleAnalytics
}

// NewDataSource creates the google analytics datasource and sets up all the routes
func NewDataSource(mux *http.ServeMux) *GoogleAnalyticsDataSource {
	cache := cache.New(300*time.Second, 5*time.Second)
	ds := &GoogleAnalyticsDataSource{
		analytics: &GoogleAnalytics{
			Cache: cache,
		},
	}

	mux.HandleFunc("/accounts", ds.handleResourceAccounts)
	mux.HandleFunc("/web-properties", ds.handleResourceWebProperties)
	mux.HandleFunc("/profiles", ds.handleResourceProfiles)
	return ds
}

// CheckHealth checks if the plugin is running properly
func (ds *GoogleAnalyticsDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Success"

	config, err := LoadSettings(req.PluginContext)

	if err != nil {
		log.DefaultLogger.Error("Fail LoadSetting", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Setting Configuration Read Fail",
		}, nil
	}

	client, err := NewGoogleClient(ctx, config)
	if err != nil {
		log.DefaultLogger.Error("Fail NewGoogleClient", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Invalid config",
		}, nil
	}

	profiles, err := client.getAllProfilesList()
	if err != nil {
		log.DefaultLogger.Error("Fail getAllProfilesList", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Invalid config",
		}, nil
	}

	testData := QueryModel{profiles[0].AccountId, profiles[0].WebPropertyId, profiles[0].Id, "yesterday", "today", "a", "ga:sessions", "ga:country"}
	res, err := client.getReport([]QueryModel{testData})

	if err != nil {
		log.DefaultLogger.Error("GET request to analyticsreporting/v4 returned error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Test Request Fail",
		}, nil
	}

	if res != nil {
		log.DefaultLogger.Info("HTTPStatusCode", "status", res.HTTPStatusCode)
		log.DefaultLogger.Info("res", res)
	}

	printResponse(res)

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

// QueryData queries for data.
func (ds *GoogleAnalyticsDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()
	config, err := LoadSettings(req.PluginContext)
	if err != nil {
		return nil, err
	}

	client, err := NewGoogleClient(ctx, config)
	if err != nil {
		return nil, err
	}

	log.DefaultLogger.Info("QueryData", "QueryData", req.Queries)
	for _, query := range req.Queries {
		frames, err := ds.analytics.Query(client, query)
		if err != nil {
			log.DefaultLogger.Error(err.Error())
			return nil, err
		}
		log.DefaultLogger.Info("QueryData", "frames", frames)
		res.Responses[query.RefID] = backend.DataResponse{*frames, err}
	}

	return res, nil
}

func writeResult(rw http.ResponseWriter, path string, val interface{}, err error) {
	response := make(map[string]interface{})
	code := http.StatusOK
	if err != nil {
		response["error"] = err.Error()
		code = http.StatusBadRequest
	} else {
		response[path] = val
	}

	body, err := json.Marshal(response)
	if err != nil {
		body = []byte(err.Error())
		code = http.StatusInternalServerError
	}
	_, err = rw.Write(body)
	if err != nil {
		code = http.StatusInternalServerError
	}
	rw.WriteHeader(code)
}

func (ds *GoogleAnalyticsDataSource) handleResourceAccounts(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := LoadSettings(httpadapter.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	res, err := ds.analytics.GetAccounts(ctx, config)
	writeResult(rw, "accounts", res, err)
}

func (ds *GoogleAnalyticsDataSource) handleResourceWebProperties(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := LoadSettings(httpadapter.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	res, err := ds.analytics.GetWebProperties(ctx, config, req.URL.Query().Get("accountId"))
	writeResult(rw, "webProperties", res, err)
}

func (ds *GoogleAnalyticsDataSource) handleResourceProfiles(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := LoadSettings(httpadapter.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}

	res, err := ds.analytics.GetProfiles(ctx, config, req.URL.Query().Get("accountId"), req.URL.Query().Get("webPropertyId"))
	writeResult(rw, "profiles", res, err)
}
