package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/gav3"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/gav4"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/patrickmn/go-cache"
)

// GoogleAnalyticsDataSource handler for google sheets
type GoogleAnalyticsDataSource struct {
	analytics       GoogleAnalytics
	resourceHandler backend.CallResourceHandler
}

// NewDataSource creates the google analytics datasource and sets up all the routes
func NewDataSource(dis backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	version := &setting.DatasourceSettings{}
	cache := cache.New(300*time.Second, 5*time.Second)
	mux := http.NewServeMux()
	err := json.Unmarshal(dis.JSONData, &version)
	if err != nil {
		return nil, err
	}

	var analytics GoogleAnalytics
	if version.Version == "v3" {
		analytics = &gav3.GoogleAnalytics{
			Cache: cache,
		}
	} else {
		analytics = &gav4.GoogleAnalytics{
			Cache: cache,
		}
	}

	ds := &GoogleAnalyticsDataSource{
		analytics:       analytics,
		resourceHandler: httpadapter.New(mux),
	}
	mux.HandleFunc("/accounts", ds.handleResourceAccounts)
	mux.HandleFunc("/web-properties", ds.handleResourceWebProperties)
	mux.HandleFunc("/profiles", ds.handleResourceProfiles)
	mux.HandleFunc("/profile/timezone", ds.handleResourceProfileTimezone)
	mux.HandleFunc("/dimensions", ds.handleResourceDimensions)
	mux.HandleFunc("/metrics", ds.handleResourceMetrics)

	return ds, nil
}

func (ds *GoogleAnalyticsDataSource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return ds.resourceHandler.CallResource(ctx, req, sender)
}

// CheckHealth checks if the plugin is running properly
func (ds *GoogleAnalyticsDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	config, err := setting.LoadSettings(req.PluginContext)

	if err != nil {
		log.DefaultLogger.Error("CheckHealth: Fail LoadSetting", "error", err.Error())
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Setting Configuration Read Fail",
		}, nil
	}
	return ds.analytics.CheckHealth(ctx, config)
}

// QueryData queries for data.
func (ds *GoogleAnalyticsDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	res := backend.NewQueryDataResponse()
	config, err := setting.LoadSettings(req.PluginContext)
	if err != nil {
		return nil, err
	}

	for _, query := range req.Queries {
		frames, err := ds.analytics.Query(ctx, config, query)
		if err != nil {
			log.DefaultLogger.Error("Fail query", "error", err)
			res.Responses[query.RefID] = backend.DataResponse{Frames: data.Frames{}, Error: err}
			continue
		}
		res.Responses[query.RefID] = backend.DataResponse{Frames: *frames, Error: err}
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
	config, err := setting.LoadSettings(httpadapter.PluginConfigFromContext(ctx))
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
	config, err := setting.LoadSettings(httpadapter.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}
	query := req.URL.Query()
	var (
		accountId = query.Get("accountId")
	)
	res, err := ds.analytics.GetWebProperties(ctx, config, accountId)
	writeResult(rw, "webProperties", res, err)
}

func (ds *GoogleAnalyticsDataSource) handleResourceProfiles(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := setting.LoadSettings(httpadapter.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}
	query := req.URL.Query()
	var (
		accountId     = query.Get("accountId")
		webPropertyId = query.Get("webPropertyId")
	)
	res, err := ds.analytics.GetProfiles(ctx, config, accountId, webPropertyId)
	writeResult(rw, "profiles", res, err)
}

func (ds *GoogleAnalyticsDataSource) handleResourceDimensions(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		return
	}
	ctx := req.Context()
	config, err := setting.LoadSettings(httpadapter.PluginConfigFromContext(ctx))

	query := req.URL.Query()
	var (
		webPropertyId = query.Get("webPropertyId")
	)
	res, err := ds.analytics.GetDimensions(ctx, config, webPropertyId)
	writeResult(rw, "dimensions", res, err)
}

func (ds *GoogleAnalyticsDataSource) handleResourceMetrics(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		return
	}
	ctx := req.Context()
	config, err := setting.LoadSettings(httpadapter.PluginConfigFromContext(ctx))

	query := req.URL.Query()
	var (
		webPropertyId = query.Get("webPropertyId")
	)
	res, err := ds.analytics.GetMetrics(ctx, config, webPropertyId)
	writeResult(rw, "metrics", res, err)
}

func (ds *GoogleAnalyticsDataSource) handleResourceProfileTimezone(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		return
	}

	ctx := req.Context()
	config, err := setting.LoadSettings(httpadapter.PluginConfigFromContext(ctx))
	if err != nil {
		writeResult(rw, "?", nil, err)
		return
	}
	query := req.URL.Query()
	var (
		accountId     = query.Get("accountId")
		webPropertyId = query.Get("webPropertyId")
		profileId     = query.Get("profileId")
	)
	res, err := ds.analytics.GetTimezone(ctx, config, accountId, webPropertyId, profileId)
	writeResult(rw, "timezone", res, err)
}
