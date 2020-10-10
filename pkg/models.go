package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type QueryModel struct {
	AccountID     string `json:"accountId"`
	WebPropertyId string `json:"webPropertyId"`
	ProfileID     string `json:"profileId"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	RefID         string `json:"refId"`
	Metric        string `json:"metric"`
	Dimension     string `json:"dimension"`

	// Not from JSON
	// TimeRange     backend.TimeRange `json:"-"`
	// MaxDataPoints int64             `json:"-"`
}

// GetQueryModel returns the well typed query model
func GetQueryModel(query backend.DataQuery) (*QueryModel, error) {
	model := &QueryModel{}

	err := json.Unmarshal(query.JSON, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %s", err.Error())
	}

	// Copy directly from the well typed query
	model.StartDate = query.TimeRange.From.Format("2006-01-02")
	model.EndDate = query.TimeRange.To.Format("2006-01-02")
	// model.TimeRange = query.TimeRange
	// model.MaxDataPoints = query.MaxDataPoints
	return model, nil
}
