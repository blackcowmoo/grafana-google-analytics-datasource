package gav3

import (
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/jinzhu/copier"
	reporting "google.golang.org/api/analyticsreporting/v4"
	"sort"
	"strconv"
	"strings"
	"time"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/util"
	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
)

func transformReportToDataFrameByDimensions(columns []*model.ColumnDefinition, rows []*reporting.ReportRow, refId string, dimensions string) (*data.Frame, error) {
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
		displayName := dimensions
		if len(dimensions) > 0 {
			displayName = displayName + "|"
		}
		field.Config = &data.FieldConfig{
			DisplayName: displayName + column.Header,
			// Unit:        column.GetUnit(),
		}
	}
	i := 0
	for rowIndex, row := range rows {
		if dimensions == strings.Join(row.Dimensions, "|") {
			for _, metrics := range row.Metrics {
				// d := row.Dimensions[dateIndex]
				for valueIndex, value := range metrics.Values {
					err := inputConverter.Set(valueIndex, rowIndex, value)
					if err != nil {
						log.DefaultLogger.Error("frame convert", "error", err.Error())
						warnings = append(warnings, err.Error())
						continue
					}
					i++
				}
			}
		}
	}

	meta["warnings"] = warnings
	frame.Meta = &data.FrameMeta{Custom: meta}
	return frame, nil
}

//                                      <--------- primary secondary --------->
var timeDimensions []string = []string{"ga:dateHourMinute", "ga:dateHour", "ga:date"}

func transformReportToDataFrames(report *reporting.Report, refId string, timezone string) ([]*data.Frame, error) {
	timeDimension := reporting.MetricHeaderEntry{
		Name: report.ColumnHeader.Dimensions[0],
		Type: "TIME",
	}
	report.ColumnHeader.MetricHeader.MetricHeaderEntries = append([]*reporting.MetricHeaderEntry{
		&timeDimension,
	}, report.ColumnHeader.MetricHeader.MetricHeaderEntries...)

	timeAddFunction, timeSubFunction := getTimeFunction(timeDimension.Name)

	tz, err := time.LoadLocation(timezone)
	if err != nil {
		log.DefaultLogger.Error("Load local timezone error", "error", err.Error())
	}

	dimensions := map[string]struct{}{}
	var parsedReportMap = make(map[string]map[int64]*reporting.ReportRow)

	for _, row := range report.Data.Rows {
		parsedRow, parsedTime := parseRow(row, tz)
		dimension := strings.Join(parsedRow.Dimensions, "|")
		if _, ok := dimensions[dimension]; !ok {
			dimensions[dimension] = struct{}{}
		}
		if inner, ok := parsedReportMap[dimension]; !ok {
			inner = make(map[int64]*reporting.ReportRow)
			parsedReportMap[dimension] = inner
		}
		parsedReportMap[dimension][parsedTime.Unix()] = parsedRow

		beforeTime := timeSubFunction(*parsedTime)
		afterTime := timeAddFunction(*parsedTime)
		if _, ok := parsedReportMap[dimension][beforeTime.Unix()]; !ok {
			copyRow := copyRowAndInit(row)
			copyRow.Metrics[0].Values[0] = beforeTime.Format(time.RFC3339)
			parsedReportMap[dimension][beforeTime.Unix()] = copyRow
		}
		if _, ok := parsedReportMap[dimension][afterTime.Unix()]; !ok {
			copyRow := copyRowAndInit(row)
			copyRow.Metrics[0].Values[0] = afterTime.Format(time.RFC3339)
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
	columns := getColumnDefinitions(report.ColumnHeader)

	for _, dimension := range dimensionKeys {
		keys := make([]int64, len(parsedReportMap[dimension]))
		i := 0
		for k := range parsedReportMap[dimension] {
			keys[i] = k
			i++
		}

		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		var parsedRows = make([]*reporting.ReportRow, len(parsedReportMap[dimension]))
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

func transformReportsResponseToDataFrames(reportsResponse *reporting.GetReportsResponse, refId string, timezone string) (*data.Frames, error) {
	var frames = make(data.Frames, 0)
	for _, report := range reportsResponse.Reports {
		frame, err := transformReportToDataFrames(report, refId, timezone)
		if err != nil {
			return nil, err
		}

		frames = append(frames, frame...)
	}

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

func getColumnDefinitions(header *reporting.ColumnHeader) []*model.ColumnDefinition {
	columns := []*model.ColumnDefinition{}
	headerRow := header.MetricHeader.MetricHeaderEntries

	for columnIndex, headerCell := range headerRow {
		name := strings.TrimSpace(headerCell.Name)
		columns = append(columns, model.NewColumnDefinition(name, columnIndex, headerCell.Type))
	}

	return columns
}

func copyRow(row *reporting.ReportRow) *reporting.ReportRow {
	var copyRow reporting.ReportRow
	copier.CopyWithOption(&copyRow.Dimensions, row.Dimensions, copier.Option{DeepCopy: true})
	copier.CopyWithOption(&copyRow.Metrics, row.Metrics, copier.Option{DeepCopy: true})
	return &copyRow
}

func copyRowAndInit(row *reporting.ReportRow) *reporting.ReportRow {
	copyRow := copyRow(row)
	copyRow.Metrics[0].Values = util.FillArray(make([]string, len(row.Metrics[0].Values)), "0")
	return copyRow
}

func parseRow(row *reporting.ReportRow, timezone *time.Location) (*reporting.ReportRow, *time.Time) {
	timeDimension := row.Dimensions[0]
	otherDimensions := row.Dimensions[1:]
	parsedTime, err := util.ParseAndTimezoneTime(timeDimension, timezone)
	if err != nil {
		log.DefaultLogger.Error("parsedTime err", "err", err.Error())
	}
	strTime := parsedTime.Format(time.RFC3339)

	// row.Metrics[0].Values = append(row.Metrics[0].Values, strTime)
	row.Metrics[0].Values = append([]string{strTime}, row.Metrics[0].Values...)
	row.Dimensions = otherDimensions
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
	default:
		add = util.AddOneHour
		sub = util.SubOneHour
	}
	return add, sub
}
