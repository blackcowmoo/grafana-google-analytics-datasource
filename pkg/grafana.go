package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

func transformReportToDataFrame(reportsResponse *reporting.GetReportsResponse, queryModel *QueryModel) (*data.Frames, error) {
	log.DefaultLogger.Info("transformReportToDataFrame", "report", reportsResponse)

	for _, report := range reportsResponse.Reports {

    for _, dimension := range report.ColumnHeader.Dimensions {

      for _, row := range report.Data.Rows {

      }
    }



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

	return nil, nil
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
