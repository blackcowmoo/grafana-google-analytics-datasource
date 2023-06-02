package gav3

import (
	"encoding/json"
	"fmt"
	"time"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	. "github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"

)


type QueryModel struct {
	AccountID         string   `json:"accountId"`
	WebPropertyID     string   `json:"webPropertyId"`
	ProfileID         string   `json:"profileId"`
	StartDate         string   `json:"startDate"`
	EndDate           string   `json:"endDate"`
	RefID             string   `json:"refId"`
	Metrics           []string `json:"metrics"`
	TimeDimension     string   `json:"timeDimension"`
	Dimensions        []string `json:"dimensions"`
	PageSize          int64    `json:"pageSize,omitempty"`
	PageToken         string   `json:"pageToken,omitempty"`
	UseNextPage       bool     `json:"useNextpage,omitempty"`
	Timezone          string   `json:"timezone,omitempty"`
	FiltersExpression string   `json:"filtersExpression,omitempty"`
	// Not from JSON
	// TimeRange     backend.TimeRange `json:"-"`
	// MaxDataPoints int64             `json:"-"`
}

// GetQueryModel returns the well typed query model
func GetQueryModel(query backend.DataQuery) (*QueryModel, error) {
	model := &QueryModel{
		PageSize:    GaReportMaxResult,
		PageToken:   "",
		UseNextPage: true,
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


func (ga *GoogleAnalyticsv3) getMetadata() (*Metadata, error) {
	res, err := http.Get(GaMetadataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata api %w", err)
	}
	defer res.Body.Close()

	metadata := Metadata{}

	err = json.NewDecoder(res.Body).Decode(&metadata)
	if err != nil {
		return nil, fmt.Errorf("fail to parsing metadata to json %w", err)
	}
	return &metadata, nil
}

func (ga *GoogleAnalyticsv3) getFilteredMetadata() ([]MetadataItem, []MetadataItem, error) {
	metadata, err := ga.getMetadata()
	if err != nil {
		return nil, nil, err
	}

	// length := int(metadata.TotalResults)
	var dimensionItems = make([]MetadataItem, 0)
	var metricItems = make([]MetadataItem, 0)
	for _, item := range metadata.Items {
		if item.Attributes.Status == "DEPRECATED" || item.Attributes.ReplacedBy != "" {
			continue
		}
		if item.Attributes.Type == AttributeTypeDimension {
			dimensionItems = append(dimensionItems, item)
		} else if item.Attributes.Type == AttributeTypeMetric {
			metricItems = append(metricItems, item)
		}
	}

	return metricItems, dimensionItems, nil
}

func (ga *GoogleAnalyticsv3) GetDimensions() ([]MetadataItem, error) {
	cacheKey := "ga:metadata:dimensions"
	if dimensions, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return dimensions.([]MetadataItem), nil
	}

	_, dimensions, err := ga.getFilteredMetadata()
	if err != nil {
		return nil, err
	}

	ga.Cache.Set(cacheKey, dimensions, time.Hour)

	return dimensions, nil
}

func (ga *GoogleAnalyticsv3) GetMetrics() ([]MetadataItem, error) {
	cacheKey := "ga:metadata:metrics"
	if metrics, _, found := ga.Cache.GetWithExpiration(cacheKey); found {
		return metrics.([]MetadataItem), nil
	}
	metrics, _, err := ga.getFilteredMetadata()
	if err != nil {
		return nil, err
	}

	ga.Cache.Set(cacheKey, metrics, time.Hour)

	return metrics, nil
}
