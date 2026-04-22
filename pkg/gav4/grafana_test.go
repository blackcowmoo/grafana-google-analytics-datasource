package gav4

import (
	"testing"
	"time"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

func TestParseRow_ValidTimeDimension(t *testing.T) {
	tz, _ := time.LoadLocation("UTC")
	row := &analyticsdata.Row{
		DimensionValues: []*analyticsdata.DimensionValue{
			{Value: "2024091215"},
			{Value: "dim-a"},
		},
		MetricValues: []*analyticsdata.MetricValue{
			{Value: "42"},
		},
	}

	parsedRow, parsedTime := parseRow(row, tz)

	if parsedRow == nil || parsedTime == nil {
		t.Fatalf("expected row to parse, got nil")
	}
	if parsedTime.Year() != 2024 || parsedTime.Month() != time.September || parsedTime.Day() != 12 {
		t.Errorf("unexpected parsed time: %s", parsedTime)
	}
	if len(parsedRow.DimensionValues) != 1 || parsedRow.DimensionValues[0].Value != "dim-a" {
		t.Errorf("time dimension should be stripped; got %#v", parsedRow.DimensionValues)
	}
}

func TestParseRow_UnparseableTimeDimensionIsSkipped(t *testing.T) {
	// Regression: Google Analytics emits an aggregate "(other)" row when the
	// response exceeds the cardinality limit. It should be skipped rather
	// than crashing on a nil time dereference (issue #121).
	tz, _ := time.LoadLocation("UTC")
	row := &analyticsdata.Row{
		DimensionValues: []*analyticsdata.DimensionValue{
			{Value: "(other)"},
			{Value: "dim-a"},
		},
		MetricValues: []*analyticsdata.MetricValue{
			{Value: "42"},
		},
	}

	parsedRow, parsedTime := parseRow(row, tz)

	if parsedRow != nil || parsedTime != nil {
		t.Fatalf("expected nil row and time for unparseable value, got row=%v time=%v", parsedRow, parsedTime)
	}
}
