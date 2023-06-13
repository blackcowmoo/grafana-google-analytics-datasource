package gav3

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/setting"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// GetQueryModel returns the well typed query model
func GetQueryModel(query backend.DataQuery) (*model.QueryModel, error) {
	model := &model.QueryModel{
		PageSize:  GaReportMaxResult,
		PageToken: "",
	}
	err := json.Unmarshal(query.JSON, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %s", err.Error())
	}

	// Copy directly from the well typed query
	timezone, err := time.LoadLocation(model.Timezone)
	if err != nil {
		return nil, fmt.Errorf("error get timezone %s", err.Error())
	}

	log.DefaultLogger.Debug("query timezone", "timezone", timezone.String())

	model.StartDate = query.TimeRange.From.In(timezone).Format("2006-01-02")
	model.EndDate = query.TimeRange.To.In(timezone).Format("2006-01-02")
	model.Dimensions = append([]string{model.TimeDimension}, model.Dimensions...)
	// model.TimeRange = query.TimeRange
	// model.MaxDataPoints = query.MaxDataPoints
	return model, nil
}

func (ga *GoogleAnalytics) getMetadata() (*model.Metadata, error) {
	res, err := http.Get(GaMetadataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata api %w", err)
	}
	defer res.Body.Close()

	metadata := model.Metadata{}

	err = json.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		return nil, fmt.Errorf("fail to parsing metadata to json %w", err)
	}
	return &metadata, nil
}

func (ga *GoogleAnalytics) getFilteredMetadata() ([]model.MetadataItem, []model.MetadataItem, error) {
	metadata, err := ga.getMetadata()
	if err != nil {
		return nil, nil, err
	}

	// length := int(metadata.TotalResults)
	var dimensionItems = make([]model.MetadataItem, 0)
	var metricItems = make([]model.MetadataItem, 0)
	for _, item := range metadata.Items {
		if item.Attributes.Status == "DEPRECATED" || item.Attributes.ReplacedBy != "" {
			continue
		}
		if item.Attributes.Type == model.AttributeTypeDimension {
			dimensionItems = append(dimensionItems, item)
		} else if item.Attributes.Type == model.AttributeTypeMetric {
			metricItems = append(metricItems, item)
		}
	}

	return metricItems, dimensionItems, nil
}

func (ga *GoogleAnalytics) GetDimensions(ctx context.Context, config *setting.DatasourceSecretSettings, propertyId string) ([]model.MetadataItem, error) {
	cacheKey := "ga:metadata:dimensions"
	if dimensions, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return dimensions.([]model.MetadataItem), nil
	}

	_, dimensions, err := ga.getFilteredMetadata()
	if err != nil {
		return nil, err
	}

	ga.Cache.Set(cacheKey, dimensions, time.Hour)

	return dimensions, nil
}

func (ga *GoogleAnalytics) GetMetrics(ctx context.Context, config *setting.DatasourceSecretSettings, propertyId string) ([]model.MetadataItem, error) {
	cacheKey := "ga:metadata:metrics"
	if metrics, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return metrics.([]model.MetadataItem), nil
	}
	metrics, _, err := ga.getFilteredMetadata()
	if err != nil {
		return nil, err
	}

	ga.Cache.Set(cacheKey, metrics, time.Hour)

	return metrics, nil
}
