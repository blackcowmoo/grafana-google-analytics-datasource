package gav4

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/util"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/jinzhu/copier"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
)

func transformReportToDataFrameByDimensions(columns []*model.ColumnDefinition, rows []*analyticsdata.Row, refId string, dimensions string) (*data.Frame, error) {
	warnings := []string{}
	meta := map[string]interface{}{}

	converters := make([]data.FieldConverter, len(columns))
	for i, column := range columns {
		fc, ok := converterMap[column.GetType()]
		if !ok {
			return nil, fmt.Errorf("unknown column type: %s", column.GetType())
		}
		converters[i] = fc
	}

	inputConverter, err := data.NewFrameInputConverter(converters, len(rows))
	if err != nil {
		return nil, err
	}

	frame := inputConverter.Frame
	frame.RefID = refId
	frame.Name = refId // TODO: should set the name from metadata
	if len(dimensions) > 0 {
		frame.Name = dimensions
	}

	for i, column := range columns {
		field := frame.Fields[i]
		field.Name = column.Header
		field.Config = &data.FieldConfig{
			DisplayName: dimensions + column.Header,
			// Unit:        column.GetUnit(),
		}
	}
	i := 0
	for rowIndex, row := range rows {
		var key string
		for _, v := range row.DimensionValues {
			if strings.TrimSpace(v.Value) != "" {
				key += v.Value + "|"
			}
		}
		if dimensions == key {
			for valueIndex, value := range row.MetricValues {
				err := inputConverter.Set(valueIndex, rowIndex, value.Value)
				if err != nil {
					log.DefaultLogger.Error("frame convert", "error", err.Error())
					warnings = append(warnings, err.Error())
					continue
				}
				i++
			}
		}
	}

	meta["warnings"] = warnings
	frame.Meta = &data.FrameMeta{Custom: meta}
	return frame, nil
}

// <--------- primary secondary --------->
var timeDimensions []string = []string{"dateHourMinute", "dateHour", "date", "firstSessionDate"}

func transformReportToDataFramesTableMode(report *analyticsdata.RunReportResponse, refId string, timezone string) ([]*data.Frame, error) {
	otherDimensions := make([]*analyticsdata.MetricHeader, 0)
	for _, dimension := range report.DimensionHeaders {
		otherDimensions = append([]*analyticsdata.MetricHeader{
			{
				Name: dimension.Name,
				Type: "STRING",
			},
		}, otherDimensions...)
	}
	report.MetricHeaders = append(otherDimensions, report.MetricHeaders...)

	for _, row := range report.Rows {
		for _, dimensionValue := range row.DimensionValues {
			row.MetricValues = append([]*analyticsdata.MetricValue{
				{
					Value: dimensionValue.Value,
				},
			}, row.MetricValues...)
			row.DimensionValues = nil
		}
	}

	var frames = make([]*data.Frame, 0)
	columns := getColumnDefinitions(report.MetricHeaders)
	frame, err := transformReportToDataFrameByDimensions(columns, report.Rows, refId, "")
	if err != nil {
		log.DefaultLogger.Error("transformReportToDataFrameByDimensions", "error", err.Error())
		return nil, err
	}

	frames = append(frames, frame)

	return frames, nil
}

func transformReportToDataFrames(report *analyticsdata.RunReportResponse, refId string, timezone string) ([]*data.Frame, error) {

	timeDimension := analyticsdata.MetricHeader{
		Name: report.DimensionHeaders[0].Name,
		Type: "TIME",
	}
	report.MetricHeaders = append([]*analyticsdata.MetricHeader{
		&timeDimension,
	}, report.MetricHeaders...)

	timeAddFunction, timeSubFunction := getTimeFunction(timeDimension.Name)

	tz, err := time.LoadLocation(timezone)
	if err != nil {
		log.DefaultLogger.Error("Load local timezone error", "error", err.Error())
	}

	dimensions := map[string]struct{}{}
	var parsedReportMap = make(map[string]map[int64]*analyticsdata.Row)

	for _, row := range report.Rows {
		parsedRow, parsedTime := parseRow(row, tz)
		var dimension string = ""
		for _, v := range parsedRow.DimensionValues {
			if strings.TrimSpace(v.Value) != "" {
				dimension += v.Value + "|"
			}
		}
		if _, ok := dimensions[dimension]; !ok {
			dimensions[dimension] = struct{}{}
		}
		if _, ok := parsedReportMap[dimension]; !ok {
			inner := make(map[int64]*analyticsdata.Row)
			parsedReportMap[dimension] = inner
		}
		parsedReportMap[dimension][parsedTime.Unix()] = parsedRow

		beforeTime := timeSubFunction(*parsedTime)
		afterTime := timeAddFunction(*parsedTime)
		if _, ok := parsedReportMap[dimension][beforeTime.Unix()]; !ok {
			copyRow := copyRowAndInit(row)
			copyRow.MetricValues[0].Value = beforeTime.Format(time.RFC3339)
			parsedReportMap[dimension][beforeTime.Unix()] = copyRow
		}
		if _, ok := parsedReportMap[dimension][afterTime.Unix()]; !ok {
			copyRow := copyRowAndInit(row)
			copyRow.MetricValues[0].Value = afterTime.Format(time.RFC3339)
			parsedReportMap[dimension][afterTime.Unix()] = copyRow
		}
	}

	var dimensionKeys = make([]string, len(dimensions))
	i := 0
	for value := range dimensions {
		dimensionKeys[i] = value
		i++
	}

	var frames = make([]*data.Frame, 0, len(dimensionKeys))
	columns := getColumnDefinitions(report.MetricHeaders)

	for _, dimension := range dimensionKeys {
		keys := make([]int64, len(parsedReportMap[dimension]))
		i := 0
		for k := range parsedReportMap[dimension] {
			keys[i] = k
			i++
		}

		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		var parsedRows = make([]*analyticsdata.Row, len(parsedReportMap[dimension]))
		i = 0
		for _, t := range keys {
			parsedRows[i] = parsedReportMap[dimension][t]
			i++
		}

		frame, err := transformReportToDataFrameByDimensions(columns, parsedRows, refId, dimension)
		if err != nil {
			log.DefaultLogger.Error("transformReportToDataFrameByDimensions", "error", err.Error())
			return nil, err
		}

		frames = append(frames, frame)
	}

	return frames, nil
}

func transformReportsResponseToDataFrames(reportsResponse *analyticsdata.RunReportResponse, refId string, timezone string, mode model.QueryMode) (*data.Frames, error) {
	var frames = make(data.Frames, 0)
	// for _, report := range reportsResponse.Rows {
	var transformReportToDataFramesFn func(*analyticsdata.RunReportResponse, string, string) ([]*data.Frame, error)
	switch mode {
	case model.TIME_SERIES:
		transformReportToDataFramesFn = transformReportToDataFrames
	case model.TABLE, model.REALTIME:
		transformReportToDataFramesFn = transformReportToDataFramesTableMode
	default:
		transformReportToDataFramesFn = transformReportToDataFramesTableMode
	}
	frame, err := transformReportToDataFramesFn(reportsResponse, refId, timezone)
	if err != nil {
		return nil, err
	}

	frames = append(frames, frame...)
	// }

	return &frames, nil
}

// timeConverter handles sheets TIME column types.
var timeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableTime,
	Converter: func(i interface{}) (interface{}, error) {
		strTime, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("expected type string, but got %T", i)
		}
		time, err := time.Parse(time.RFC3339, strTime)
		if err != nil {
			log.DefaultLogger.Error("timeConverter failed", "error", err)
			return nil, err
		}
		return &time, nil
	},
}

// stringConverter handles sheets STRING column types.
var stringConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableString,
	Converter: func(i interface{}) (interface{}, error) {
		value, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("expected type string, but got %T", i)
		}

		return &value, nil
	},
}

// numberConverter handles sheets STRING column types.
var numberConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(i interface{}) (interface{}, error) {
		value, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("expected type string, but got %T", i)
		}

		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("expected type string, but got %T", value)
		}

		return &num, nil
	},
}

// converterMap is a map sheets.ColumnType to fieldConverter and
// is used to create a data.FrameInputConverter for a returned sheet.
var converterMap = map[model.ColumnType]data.FieldConverter{
	"TIME":   timeConverter,
	"STRING": stringConverter,
	"NUMBER": numberConverter,
}

func getColumnDefinitions(header []*analyticsdata.MetricHeader) []*model.ColumnDefinition {
	columns := []*model.ColumnDefinition{}

	for columnIndex, headerCell := range header {
		name := strings.TrimSpace(headerCell.Name)
		columns = append(columns, model.NewColumnDefinition(name, columnIndex, headerCell.Type))
	}

	return columns
}

func copyRow(row *analyticsdata.Row) *analyticsdata.Row {
	var copyRow analyticsdata.Row
	copier.CopyWithOption(&copyRow.DimensionValues, row.DimensionValues, copier.Option{DeepCopy: true})
	copier.CopyWithOption(&copyRow.MetricValues, row.MetricValues, copier.Option{DeepCopy: true})
	return &copyRow
}

func copyRowAndInit(row *analyticsdata.Row) *analyticsdata.Row {
	copyRow := copyRow(row)
	copyRow.MetricValues = fillRow(make([]*analyticsdata.MetricValue, len(row.MetricValues)), analyticsdata.MetricValue{Value: "0"})
	return copyRow
}

func fillRow(array []*analyticsdata.MetricValue, v analyticsdata.MetricValue) []*analyticsdata.MetricValue {
	for i := range array {
		tmp := v
		array[i] = &tmp
	}
	return array
}

func parseRow(row *analyticsdata.Row, timezone *time.Location) (*analyticsdata.Row, *time.Time) {
	timeDimension := row.DimensionValues[0].Value
	otherDimensions := row.DimensionValues[1:]
	parsedTime, err := util.ParseAndTimezoneTime(timeDimension, timezone)
	if err != nil {
		log.DefaultLogger.Error("parseRow: Failed to parse time dimension", "error", err.Error())
	}
	strTime := parsedTime.Format(time.RFC3339)

	// row.Metrics[0].Values = append(row.Metrics[0].Values, strTime)
	t := &analyticsdata.MetricValue{
		Value: strTime,
	}
	row.MetricValues = append([]*analyticsdata.MetricValue{t}, row.MetricValues...)
	row.DimensionValues = otherDimensions
	return row, parsedTime
}

func getTimeFunction(timeDimension string) (func(time.Time) time.Time, func(time.Time) time.Time) {
	var add, sub func(time.Time) time.Time
	switch timeDimension {
	case timeDimensions[0]:
		add = util.AddOneMinute
		sub = util.SubOneMinute
		break
	case timeDimensions[1]:
		add = util.AddOneHour
		sub = util.SubOneHour
		break
	case timeDimensions[2]:
		add = util.AddOneDay
		sub = util.SubOneDay
		break
	case timeDimensions[3]:
		add = util.AddOneDay
		sub = util.SubOneDay
		break
	default:
		add = util.AddOneHour
		sub = util.SubOneHour
	}
	return add, sub
}
