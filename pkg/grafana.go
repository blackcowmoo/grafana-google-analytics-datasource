package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

func transformReportToDataFrame(report *reporting.Report, refId string) (*data.Frame, error) {
	columns, _ := getColumnDefinitions(report.ColumnHeader)
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

	inputConverter, err := data.NewFrameInputConverter(converters, len(report.Data.Rows))
	if err != nil {
		return nil, err
	}

	frame := inputConverter.Frame
	frame.RefID = refId
	frame.Name = refId // TODO: should set the name from metadata

	for i, column := range columns {
		field := frame.Fields[i]
		field.Name = column.Header
		field.Config = &data.FieldConfig{
			DisplayName: column.Header,
			// Unit:        column.GetUnit(),
		}
	}

	for rowIndex, row := range report.Data.Rows {
		for metricIndex, metrics := range row.Metrics {
			// d := row.Dimensions[dateIndex]
			for _, value := range metrics.Values {
				err := inputConverter.Set(metricIndex, rowIndex, value)
				if err != nil {
					warnings = append(warnings, err.Error())
				}
			}
		}
	}

	meta["warnings"] = warnings
	frame.Meta = &data.FrameMeta{Custom: meta}
	return frame, nil
}

func transformReportsResponseToDataFrames(reportsResponse *reporting.GetReportsResponse, refId string) (*data.Frames, error) {
	var frames = make(data.Frames, 0)
	for _, report := range reportsResponse.Reports {
		frame, err := transformReportToDataFrame(report, refId)
		if err != nil {
			return nil, err
		}

		frames = append(frames, frame)
	}

	return &frames, nil
}

// timeConverter handles sheets TIME column types.
var timeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableTime,
	Converter: func(i interface{}) (interface{}, error) {
		return nil, fmt.Errorf("error: %s", i)
		// return i, nil
		// var t *time.Time
		// cellData, ok := i.(*sheets.CellData)
		// if !ok {
		// 	return t, fmt.Errorf("expected type *sheets.CellData, but got %T", i)
		// }
		// parsedTime, err := dateparse.ParseLocal(cellData.FormattedValue)
		// if err != nil {
		// 	return t, fmt.Errorf("Error while parsing date '%v'", cellData.FormattedValue)
		// }
		// return &parsedTime, nil
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
var converterMap = map[ColumnType]data.FieldConverter{
	"TIME":   timeConverter,
	"STRING": stringConverter,
	"NUMBER": numberConverter,
}

func getColumnDefinitions(header *reporting.ColumnHeader) ([]*ColumnDefinition, int) {
	columns := []*ColumnDefinition{}
	headerRow := header.MetricHeader.MetricHeaderEntries

	for columnIndex, headerCell := range headerRow {
		name := strings.TrimSpace(headerCell.Name)
		columns = append(columns, NewColumnDefinition(name, columnIndex, headerCell.Type))
	}

	return columns, -1
}
