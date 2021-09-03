package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
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

	log.DefaultLogger.Info("query timezone", "timezone", timezone.String())

	model.StartDate = query.TimeRange.From.In(timezone).Format("2006-01-02")
	model.EndDate = query.TimeRange.To.In(timezone).Format("2006-01-02")
	model.Dimensions = append([]string{model.TimeDimension}, model.Dimensions...)
	// model.TimeRange = query.TimeRange
	// model.MaxDataPoints = query.MaxDataPoints
	return model, nil
}

// ColumnType is the set of possible column types
type ColumnType string

const (
	// ColumTypeTime is the TIME type
	ColumTypeTime ColumnType = "TIME"
	// ColumTypeNumber is the NUMBER type
	ColumTypeNumber ColumnType = "NUMBER"
	// ColumTypeString is the STRING type
	ColumTypeString ColumnType = "STRING"
)

// ColumnDefinition represents a spreadsheet column definition.
type ColumnDefinition struct {
	Header      string
	ColumnIndex int
	columnType  ColumnType
}

// GetType gets the type of a ColumnDefinition.
func (cd *ColumnDefinition) GetType() ColumnType {
	return cd.columnType
}

func getColumnType(headerType string) ColumnType {
	switch headerType {
	case "INTEGER", "FLOAT", "CURRENCY", "PERCENT":
		return ColumTypeNumber
	case "TIME":
		return ColumTypeTime
	default:
		return ColumTypeString
	}
}

// NewColumnDefinition creates a new ColumnDefinition.
func NewColumnDefinition(header string, index int, headerType string) *ColumnDefinition {

	return &ColumnDefinition{
		Header:      header,
		ColumnIndex: index,
		columnType:  getColumnType(headerType),
	}
}

// Metadata
const (
	AttributeTypeDimension AttributeType = "DIMENSION"
	AttributeTypeMetric    AttributeType = "METRIC"
)

type Metadata struct {
	Kind           string         `json:"kind"`
	Etag           string         `json:"etag"`
	TotalResults   int64          `json:"totalResults"`
	AttributeNames []string       `json:"attributeNames"`
	Items          []MetadataItem `json:"items"`
}

type MetadataItem struct {
	ID         string                `json:"id"`
	Kind       string                `json:"kind"`
	Attributes MetadataItemAttribute `json:"attributes"`
}

type MetadataItemAttribute struct {
	Type              AttributeType `json:"type,omitempty"`
	DataType          string        `json:"dataType,omitempty"`
	Group             string        `json:"group,omitempty"`
	Status            string        `json:"status,omitempty"`
	UIName            string        `json:"uiName,omitempty"`
	Description       string        `json:"description,omitempty"`
	AllowedInSegments string        `json:"allowedInSegments,omitempty"`
	AddedInAPIVersion string        `json:"addedInApiVersion,omitempty"`
	ReplacedBy        string        `json:"replacedBy,omitempty"`
}

type AttributeType string

func (ga *GoogleAnalytics) getMetadata() (*Metadata, error) {
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

func (ga *GoogleAnalytics) getFilteredMetadata() ([]MetadataItem, []MetadataItem, error) {
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

func (ga *GoogleAnalytics) GetDimensions() ([]MetadataItem, error) {
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

func (ga *GoogleAnalytics) GetMetrics() ([]MetadataItem, error) {
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
