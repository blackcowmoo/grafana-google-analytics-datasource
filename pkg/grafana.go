package main

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

func transformReportToDataFrame(report *reporting.Report, refId string) (*data.Frame, error) {
	columns := getColumnDefinitions(report.ColumnHeader)
	warnings := []string{}

	converters := make([]data.FieldConverter, len(columns))
	for i, column := range columns {
		fc, ok := converterMap[column.GetType()]
		if !ok {
			return nil, fmt.Errorf("unknown column type: %s", column.GetType())
		}
		converters[i] = fc
	}

	inputConverter, err := data.NewFrameInputConverter(converters, len(sheet.RowData)-start)
	if err != nil {
		return nil, err
	}
	frame := inputConverter.Frame
	frame.RefID = refID
	frame.Name = refID // TODO: should set the name from metadata

	for i, column := range columns {
		field := frame.Fields[i]
		field.Name = column.Header
		field.Config = &data.FieldConfig{
			DisplayName: column.Header,
			Unit:        column.GetUnit(),
		}
		if column.HasMixedTypes() {
			warning := fmt.Sprintf("Multiple data types found in column %q. Using string data type", column.Header)
			warnings = append(warnings, warning)
			backend.Logger.Warn(warning)
		}

		if column.HasMixedUnits() {
			warning := fmt.Sprintf("Multiple units found in column %q. Formatted value will be used", column.Header)
			warnings = append(warnings, warning)
			backend.Logger.Warn(warning)
		}
	}

	for rowIndex := start; rowIndex < len(sheet.RowData); rowIndex++ {
		for columnIndex, cellData := range sheet.RowData[rowIndex].Values {
			if columnIndex >= len(columns) {
				continue
			}

			// Skip any empty values
			if cellData.FormattedValue == "" {
				continue
			}

			err := inputConverter.Set(columnIndex, rowIndex-start, cellData)
			if err != nil {
				warnings = append(warnings, err.Error())
			}
		}
	}

	meta["warnings"] = warnings
	meta["spreadsheetId"] = qm.Spreadsheet
	meta["range"] = qm.Range
	frame.Meta = &data.FrameMeta{Custom: meta}
	backend.Logger.Debug("frame.Meta: %s", spew.Sdump(frame.Meta))
	return frame, nil
}

func transformReportsResponseToDataFrames(reportsResponse *reporting.GetReportsResponse, refId string) (*data.Frames, error) {
	log.DefaultLogger.Info("transformReportToDataFrame", "report", reportsResponse)

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

func getColumnDefinitions(header *reporting.ColumnHeader) []*ColumnDefinition {
	columns := []*ColumnDefinition{}
	columnMap := map[string]bool{}
	headerRow := header.MetricHeader

	if len(rows) > 1 {
		start = 1
		for columnIndex, headerCell := range headerRow {
			name := getUniqueColumnName(strings.TrimSpace(headerCell.FormattedValue), columnIndex, columnMap)
			columnMap[name] = true
			columns = append(columns, NewColumnDefinition(name, columnIndex))
		}
	} else {
		for columnIndex := range headerRow {
			name := getUniqueColumnName("", columnIndex, columnMap)
			columnMap[name] = true
			columns = append(columns, NewColumnDefinition(name, columnIndex))
		}
	}

	// Check the types for each column
	for rowIndex := start; rowIndex < len(rows); rowIndex++ {
		for _, column := range columns {
			if column.ColumnIndex < len(rows[rowIndex].Values) {
				column.CheckCell(rows[rowIndex].Values[column.ColumnIndex])
			}
		}
	}

	return columns, start
}
