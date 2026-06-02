package model

import (
	"encoding/json"
	"testing"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

func TestGetColumnType(t *testing.T) {
	numberTypes := []string{
		"TYPE_INTEGER", "TYPE_FLOAT", "TYPE_CURRENCY", "TYPE_MILLISECONDS", "TYPE_SECONDS",
	}
	for _, ty := range numberTypes {
		if got := getColumnType(ty); got != ColumTypeNumber {
			t.Errorf("getColumnType(%q) = %q, want %q", ty, got, ColumTypeNumber)
		}
	}

	if got := getColumnType("TIME"); got != ColumTypeTime {
		t.Errorf("getColumnType(TIME) = %q, want %q", got, ColumTypeTime)
	}

	// Anything else falls back to string (including legacy UA numeric types).
	for _, ty := range []string{"STRING", "BOOLEAN", "UNKNOWN", "", "INTEGER", "CURRENCY", "PERCENT"} {
		if got := getColumnType(ty); got != ColumTypeString {
			t.Errorf("getColumnType(%q) = %q, want %q", ty, got, ColumTypeString)
		}
	}
}

func TestNewColumnDefinition(t *testing.T) {
	cd := NewColumnDefinition("activeUsers", 3, "TYPE_INTEGER")
	if cd.Header != "activeUsers" {
		t.Errorf("Header = %q, want %q", cd.Header, "activeUsers")
	}
	if cd.ColumnIndex != 3 {
		t.Errorf("ColumnIndex = %d, want 3", cd.ColumnIndex)
	}
	if cd.GetType() != ColumTypeNumber {
		t.Errorf("GetType() = %q, want %q", cd.GetType(), ColumTypeNumber)
	}
}

// TestQueryModel_FilterRoundTrip verifies that DimensionFilter and MetricFilter
// survive a JSON round-trip without loss, which is the path taken when the
// frontend serialises the query and the backend deserialises it.
func TestQueryModel_FilterRoundTrip(t *testing.T) {
	filter := analyticsdata.FilterExpression{
		Filter: &analyticsdata.Filter{
			FieldName: "country",
			StringFilter: &analyticsdata.StringFilter{
				MatchType: "EXACT",
				Value:     "US",
			},
		},
	}
	qm := QueryModel{
		WebPropertyID:   "properties/123",
		DimensionFilter: &filter,
		MetricFilter:    &filter,
	}

	b, err := json.Marshal(qm)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got QueryModel
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.DimensionFilter == nil || got.DimensionFilter.Filter == nil || got.DimensionFilter.Filter.FieldName != "country" {
		t.Errorf("DimensionFilter round-trip failed: %+v", got.DimensionFilter)
	}
	if got.MetricFilter == nil || got.MetricFilter.Filter == nil || got.MetricFilter.Filter.FieldName != "country" {
		t.Errorf("MetricFilter round-trip failed: %+v", got.MetricFilter)
	}
}

// TestQueryModel_EmptyFilterIsOmitted verifies that zero-value filters are
// omitted from JSON so the backend nil-checks in client.go behave correctly.
func TestQueryModel_EmptyFilterIsOmitted(t *testing.T) {
	qm := QueryModel{WebPropertyID: "properties/456"}

	b, err := json.Marshal(qm)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// Both filter fields carry omitempty — they must not appear in the output
	// when no filter is set, so the nil-checks in getReport / getRealtimeReport
	// keep working correctly.
	if _, ok := raw["dimensionFilter"]; ok {
		t.Error("expected dimensionFilter to be omitted when zero, but it was present")
	}
	if _, ok := raw["metricFilter"]; ok {
		t.Error("expected metricFilter to be omitted when zero, but it was present")
	}
}
