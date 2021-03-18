package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

func transformReportToDataFrameByDimensions(columns []*ColumnDefinition, report *reporting.Report, refId string, dimensions string) (*data.Frame, error) {
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
	if len(dimensions) > 0 {
		frame.Name = dimensions
	}

	for i, column := range columns {
		field := frame.Fields[i]
		field.Name = column.Header
		displayName := dimensions
		if len(dimensions) > 0 {
			displayName = displayName + ":"
		}
		field.Config = &data.FieldConfig{
			DisplayName: displayName + column.Header,
			// Unit:        column.GetUnit(),
		}
	}

	for rowIndex, row := range report.Data.Rows {
		if dimensions == strings.Join(row.Dimensions, "|") {
			for _, metrics := range row.Metrics {
				// d := row.Dimensions[dateIndex]
				for valueIndex, value := range metrics.Values {
					err := inputConverter.Set(valueIndex, rowIndex, value)
					if err != nil {
						warnings = append(warnings, err.Error())
					}
				}
			}
		}
	}

	meta["warnings"] = warnings
	frame.Meta = &data.FrameMeta{Custom: meta}
	return frame, nil
}

var timeDimensions []string = []string{"ga:dateHourMinute", "ga:dateHour", "ga:date"}

func transformReportToDataFrames(report *reporting.Report, refId string, timezone string) ([]*data.Frame, error) {
	var metricDateDimensionIndex int = -1
Exit:
	for _, tDimension := range timeDimensions {
		for index, dimension := range report.ColumnHeader.Dimensions {
			if tDimension == dimension {
				metricDateDimensionIndex = index
				break Exit
			}
		}
	}

	if metricDateDimensionIndex >= 0 {
		report.ColumnHeader.MetricHeader.MetricHeaderEntries = append(report.ColumnHeader.MetricHeader.MetricHeaderEntries, &reporting.MetricHeaderEntry{
			Name: report.ColumnHeader.Dimensions[metricDateDimensionIndex],
			Type: "TIME",
		})
	}

	var dateDimensionsIndex []int = []int{}
	var newDimensions []string = []string{}
	for index, dimension := range report.ColumnHeader.Dimensions {
		for _, tDimension := range timeDimensions {
			if dimension == tDimension {
				dateDimensionsIndex = append(dateDimensionsIndex, index)
			} else {
				newDimensions = append(newDimensions, dimension)
			}
		}
	}

	report.ColumnHeader.Dimensions = newDimensions
	var dimensions []string = []string{}
	for _, row := range report.Data.Rows {
		var rowDimensions []string = []string{}
		for index, dimension := range row.Dimensions {
			var find bool = false
			for _, dateDimensionIndex := range dateDimensionsIndex {
				if index == dateDimensionIndex {
					find = true
					if metricDateDimensionIndex == index {
						timezone, err := time.LoadLocation(timezone)
						if err != nil {
							log.DefaultLogger.Warn("LoadTimeZone", "err", err.Error())
						}
						parsedTime, err := parseAndTimezoneTime(dimension, timezone)
						if err != nil {
							log.DefaultLogger.Warn("paresdTime", "err", err.Error())
						}
						sTime := parsedTime.Format(time.RFC3339)
						row.Metrics[0].Values = append(row.Metrics[0].Values, sTime)
					}
				}
			}
			if !find {
				rowDimensions = append(rowDimensions, dimension)
			}
		}
		row.Dimensions = rowDimensions
		find := false
		for _, dimension := range dimensions {
			if strings.Join(rowDimensions, "|") == dimension {
				find = true
				break
			}
		}

		if !find {
			dimensions = append(dimensions, strings.Join(rowDimensions, "|"))
		}
	}

	var frames = make([]*data.Frame, 0)
	columns := getColumnDefinitions(report.ColumnHeader)

	for _, dimension := range dimensions {
		frame, err := transformReportToDataFrameByDimensions(columns, report, refId, dimension)
		if err != nil {
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

func padRightSide(str string, item string, count int) string {
	return str + strings.Repeat(item, count)
}

func parseAndTimezoneTime(sTime string, timezone *time.Location) (*time.Time, error) {
	time, err := time.ParseInLocation("200601021504", padRightSide(sTime, "0", 12-len(sTime)), timezone)

	if err != nil {
		log.DefaultLogger.Info("timeConverter", "err", err)
		return nil, err
	}
	return &time, nil
}

// timeConverter handles sheets TIME column types.
var timeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableTime,
	Converter: func(i interface{}) (interface{}, error) {
		sTime, ok := i.(string)
		if !ok {
			return nil, fmt.Errorf("expected type string, but got %T", i)
		}
		time, err := time.Parse(time.RFC3339, sTime)
		if err != nil {
			log.DefaultLogger.Info("timeConverter", "err", err)
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
var converterMap = map[ColumnType]data.FieldConverter{
	"TIME":   timeConverter,
	"STRING": stringConverter,
	"NUMBER": numberConverter,
}

func getColumnDefinitions(header *reporting.ColumnHeader) []*ColumnDefinition {
	columns := []*ColumnDefinition{}
	headerRow := header.MetricHeader.MetricHeaderEntries

	for columnIndex, headerCell := range headerRow {
		name := strings.TrimSpace(headerCell.Name)
		columns = append(columns, NewColumnDefinition(name, columnIndex, headerCell.Type))
	}

	return columns
}
