package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	reporting "google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/sheets/v4"
)

func transformMetricToDataFrame(rows []*reporting.ReportRow, refId string) (*data.Frame, error) {
	converters := make([]data.FieldConverter, len(rows))
	for i, column := range columns {
		fc, ok := converterMap[column.GetType()]
		if !ok {
			return nil, fmt.Errorf("unknown column type: %s", column.GetType())
		}
		converters[i] = fc
	}

}

func transformReportToDataFrames(report *reporting.Report, refId string) (*data.Frames, error) {
	columns, dateIndex := getColumnDefinitions(report.ColumnHeader)
	// warnings := []string{}

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
			d := row.Dimensions[dateIndex]
			for valueIndex, value := range metrics.Values {
				err := inputConverter.Set(columnIndex, rowIndex-start, cellData)
				if err != nil {
					warnings = append(warnings, err.Error())
				}
			}
		}
	}

	log.DefaultLogger.Info("transformReportToDataFrame", "frame", frame)
	return nil, nil

	// meta["warnings"] = warnings
	// meta["spreadsheetId"] = qm.Spreadsheet
	// meta["range"] = qm.Range
	// frame.Meta = &data.FrameMeta{Custom: meta}
	// backend.Logger.Debug("frame.Meta: %s", spew.Sdump(frame.Meta))
	// return frame, nil
}

func transformReportsResponseToDataFrames(reportsResponse *reporting.GetReportsResponse, refId string) (*data.Frames, error) {
	log.DefaultLogger.Info("transformReportToDataFrame", "report", reportsResponse)

	var frames = make(data.Frames, len(reportsResponse.Reports))
	for _, report := range reportsResponse.Reports {
		frame, err := transformReportToDataFrames(report, refId)
		if err != nil {
			return nil, err
		}

		frames = append(frames, frame)
	}

	// dataFrame data.Frames = frames
	return &frames, nil

	// var frames = make([]*data.Frame, len(reportsResponse.Reports))
	// for reportIndex, report := range reportsResponse.Reports {
	// 	log.DefaultLogger.Info("transformReportToDataFrame", "report", report)
	// 	frames[reportIndex] = &data.Frame{Name: refId, RefID: refId, Meta: &data.FrameMeta{}}
	// 	var fields = make([]*data.Field, len(report.ColumnHeader.Dimensions))
	// 	for _, dimension := range report.ColumnHeader.Dimensions {
	// 		var v data.vector = data.vector{}
	// 		// for _, row := range report.Data.Rows {
	// 		// 	v.Append(row)
	// 		// }
	// 		var field = &data.Field{Name: dimension, Labels: data.Labels{}, Config: &data.FieldConfig{}, vector: v}
	// 		log.DefaultLogger.Info("transformReportToDataFrame:field", "field", field)
	// 	}
	// 	frames[reportIndex].Fields = fields
	// }

	// log.DefaultLogger.Info("transformReportToDataFrame:frame", "frames", frames)
	// var dataFrames data.Frames = frames
	// return &dataFrames, nil

	// columns := report.ColumnHeader.Dimensions
	// converters := make([]data.FieldConverter, len(columns))
	// for i, column := range columns {
	// 	fc, ok := converterMap[column.GetType()]
	// 	if !ok {
	// 		return nil, fmt.Errorf("unknown column type: %s", column.GetType())
	// 	}
	// 	converters[i] = fc
	// }

	// 	for _, row := range report.Data.Rows {
	// 		// row.Dimensions = append(row.Dimensions[:tokenIndex], row.Dimensions[tokenIndex+1:]...)
	// 		for _, dimension := range row.Dimensions {
	// 			log.DefaultLogger.Info("transformReportToDataFrame:report:row:dimension", "dimension", dimension)
	// 		}
	// 		log.DefaultLogger.Info("transformReportToDataFrame:report:row", "row", row, "dimensions", row.Dimensions)
	// 	}
	// 	log.DefaultLogger.Info("transformReportToDataFrame:report", "report", report)
	// }

	// return nil, nil
	// columns, start := getColumnDefinitions(sheet.RowData)
	// warnings := []string{}

	// converters := make([]data.FieldConverter, len(columns))
	// for i, column := range columns {
	// 	fc, ok := converterMap[column.GetType()]
	// 	if !ok {
	// 		return nil, fmt.Errorf("unknown column type: %s", column.GetType())
	// 	}
	// 	converters[i] = fc
	// }

	// inputConverter, err := data.NewFrameInputConverter(converters, len(sheet.RowData)-start)
	// if err != nil {
	// 	return nil, err
	// }
	// frame := inputConverter.Frame
	// frame.RefID = queryModel.RefID
	// frame.Name = queryModel.RefID // TODO: should set the name from metadata

	// for i, column := range columns {
	// 	field := frame.Fields[i]
	// 	field.Name = column.Header
	// 	field.Config = &data.FieldConfig{
	// 		DisplayName: column.Header,
	// 		Unit:        column.GetUnit(),
	// 	}
	// 	if column.HasMixedTypes() {
	// 		warning := fmt.Sprintf("Multiple data types found in column %q. Using string data type", column.Header)
	// 		warnings = append(warnings, warning)
	// 		backend.Logger.Warn(warning)
	// 	}

	// 	if column.HasMixedUnits() {
	// 		warning := fmt.Sprintf("Multiple units found in column %q. Formatted value will be used", column.Header)
	// 		warnings = append(warnings, warning)
	// 		backend.Logger.Warn(warning)
	// 	}
	// }

	// // for rowIndex := start; rowIndex < len(sheet.RowData); rowIndex++ {
	// // 	for columnIndex, cellData := range sheet.RowData[rowIndex].Values {
	// // 		if columnIndex >= len(columns) {
	// // 			continue
	// // 		}

	// // 		// Skip any empty values
	// // 		if cellData.FormattedValue == "" {
	// // 			continue
	// // 		}

	// // 		err := inputConverter.Set(columnIndex, rowIndex-start, cellData)
	// // 		if err != nil {
	// // 			warnings = append(warnings, err.Error())
	// // 		}
	// // 	}
	// // }

	// var meta = make(map[string]interface{})
	// meta["warnings"] = warnings
	// meta["spreadsheetId"] = qm.Spreadsheet
	// meta["range"] = qm.Range
	// frame.Meta = &data.FrameMeta{Custom: meta}
	// return backend.DataResponse{frame, nil}
}

// timeConverter handles sheets TIME column types.
var timeConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableTime,
	Converter: func(i interface{}) (interface{}, error) {
		var t *time.Time
		cellData, ok := i.(*sheets.CellData)
		if !ok {
			return t, fmt.Errorf("expected type *sheets.CellData, but got %T", i)
		}
		parsedTime, err := dateparse.ParseLocal(cellData.FormattedValue)
		if err != nil {
			return t, fmt.Errorf("Error while parsing date '%v'", cellData.FormattedValue)
		}
		return &parsedTime, nil
	},
}

// stringConverter handles sheets STRING column types.
var stringConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableString,
	Converter: func(i interface{}) (interface{}, error) {
		var s *string
		cellData, ok := i.(*sheets.CellData)
		if !ok {
			return s, fmt.Errorf("expected type *sheets.CellData, but got %T", i)
		}
		return &cellData.FormattedValue, nil
	},
}

// numberConverter handles sheets STRING column types.
var numberConverter = data.FieldConverter{
	OutputFieldType: data.FieldTypeNullableFloat64,
	Converter: func(i interface{}) (interface{}, error) {
		cellData, ok := i.(*sheets.CellData)
		if !ok {
			return nil, fmt.Errorf("expected type *sheets.CellData, but got %T", i)
		}
		if &cellData.EffectiveValue.NumberValue != nil {
			return cellData.EffectiveValue.NumberValue, nil
		}
		return nil, nil
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
