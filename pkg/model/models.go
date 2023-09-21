package model

import (
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

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
	case /*gav4*/ "TYPE_INTEGER", "TYPE_FLOAT", "TYPE_CURRENCY" /*gav3*/, "CURRENCY", "INTEGER", "FLOAT", "PERCENT":
		return ColumTypeNumber
	case "TYPE_MILLISECONDS", "TYPE_SECONDS", "TIME":
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

type AccountSummary struct {
	Account, DisplayName string
	PropertySummaries    []*PropertySummary
}

type PropertySummary struct {
	Property, DisplayName, Parent string
	ProfileSummaries              []*ProfileSummary
}

type ProfileSummary struct {
	Profile, DisplayName, Parent, Type string
}

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
	Offset            int64    `json:"offset,omitempty"`
	Mode              string   `json:"mode,omitempty"`
	// TODO type convert
	DimensionFilter  analyticsdata.FilterExpression `json:"dimensionFilter,omitempty"`
	// Not from JSON
	// TimeRange     backend.TimeRange `json:"-"`
	// MaxDataPoints int64             `json:"-"`
}
