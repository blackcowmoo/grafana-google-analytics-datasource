package model

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
