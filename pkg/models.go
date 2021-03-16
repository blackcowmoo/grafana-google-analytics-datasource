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
	PageSize      int64  `json:"pageSize,omitempty"`
	PageToken     string `json:"pageToken,omitempty"`
	UseNextPage   bool   `json:"useNextpage,omitempty"`
	// Not from JSON
	// TimeRange     backend.TimeRange `json:"-"`
	// MaxDataPoints int64             `json:"-"`
}

// GetQueryModel returns the well typed query model
func GetQueryModel(query backend.DataQuery) (*QueryModel, error) {
	model := &QueryModel{
		PageSize:    100000,
		PageToken:   "",
		UseNextPage: true,
	}

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
	case "INTEGER":
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
