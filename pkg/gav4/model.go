package gav4

import (
	"encoding/json"
	"fmt"
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

	log.DefaultLogger.Debug("query timezone", "timezone", timezone.String())

	model.StartDate = query.TimeRange.From.In(timezone).Format("2006-01-02")
	model.EndDate = query.TimeRange.To.In(timezone).Format("2006-01-02")
	model.Dimensions = append([]string{model.TimeDimension}, model.Dimensions...)
	// model.TimeRange = query.TimeRange
	// model.MaxDataPoints = query.MaxDataPoints
	return model, nil
}
