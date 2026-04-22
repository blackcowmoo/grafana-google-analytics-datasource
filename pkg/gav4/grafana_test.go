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

func TestTransformReportToDataFrames_FiltersSubDayRange(t *testing.T) {
	// Regression: Grafana time range smaller than a day (e.g. "Last 6 hours")
	// must drop GA rows whose bucket falls outside the range (issue #108).
	header := &analyticsdata.DimensionHeader{Name: "dateHour"}
	metricHeader := &analyticsdata.MetricHeader{Name: "activeUsers", Type: "TYPE_INTEGER"}
	mkRow := func(t, v string) *analyticsdata.Row {
		return &analyticsdata.Row{
			DimensionValues: []*analyticsdata.DimensionValue{{Value: t}},
			MetricValues:    []*analyticsdata.MetricValue{{Value: v}},
		}
	}
	report := &analyticsdata.RunReportResponse{
		DimensionHeaders: []*analyticsdata.DimensionHeader{header},
		MetricHeaders:    []*analyticsdata.MetricHeader{metricHeader},
		Rows: []*analyticsdata.Row{
			mkRow("2024091202", "1"), // before window (02:00)
			mkRow("2024091204", "2"), // inside window (04:00)
			mkRow("2024091205", "3"), // inside window (05:00)
			mkRow("2024091210", "4"), // after window (10:00)
		},
	}
	from := time.Date(2024, 9, 12, 3, 0, 0, 0, time.UTC)
	to := time.Date(2024, 9, 12, 9, 0, 0, 0, time.UTC)

	frames, err := transformReportToDataFrames(report, "A", "UTC", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) == 0 {
		t.Fatalf("expected at least one frame, got 0")
	}
	// The time-series transform also synthesizes adjacent zero-filled buckets.
	// We only assert that the rows outside the window are dropped.
	frame := frames[0]
	times, ok := frame.Fields[0].At(0).(*time.Time)
	_ = times
	_ = ok
	// Count actual data rows (non-zero activeUsers). With filtering, only the
	// two in-range rows (values 2 and 3) should survive.
	nonZero := 0
	for i := 0; i < frame.Fields[1].Len(); i++ {
		v := frame.Fields[1].At(i)
		if f, ok := v.(*float64); ok && f != nil && *f > 0 {
			nonZero++
		}
	}
	if nonZero != 2 {
		t.Errorf("expected 2 in-range data rows, got %d", nonZero)
	}
}

func TestTransformReportToDataFrames_NoRangeKeepsAllRows(t *testing.T) {
	// When caller passes zero-value from/to, filtering is bypassed and all
	// rows survive (back-compat path).
	header := &analyticsdata.DimensionHeader{Name: "dateHour"}
	metricHeader := &analyticsdata.MetricHeader{Name: "activeUsers", Type: "TYPE_INTEGER"}
	mkRow := func(t, v string) *analyticsdata.Row {
		return &analyticsdata.Row{
			DimensionValues: []*analyticsdata.DimensionValue{{Value: t}},
			MetricValues:    []*analyticsdata.MetricValue{{Value: v}},
		}
	}
	report := &analyticsdata.RunReportResponse{
		DimensionHeaders: []*analyticsdata.DimensionHeader{header},
		MetricHeaders:    []*analyticsdata.MetricHeader{metricHeader},
		Rows: []*analyticsdata.Row{
			mkRow("2024091202", "1"),
			mkRow("2024091210", "4"),
		},
	}

	frames, err := transformReportToDataFrames(report, "A", "UTC", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(frames) == 0 {
		t.Fatalf("expected at least one frame")
	}
	nonZero := 0
	frame := frames[0]
	for i := 0; i < frame.Fields[1].Len(); i++ {
		v := frame.Fields[1].At(i)
		if f, ok := v.(*float64); ok && f != nil && *f > 0 {
			nonZero++
		}
	}
	if nonZero != 2 {
		t.Errorf("expected both rows kept without range filter, got %d non-zero", nonZero)
	}
}
