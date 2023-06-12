package gav4

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blackcowmoo/grafana-google-analytics-dataSource/pkg/model"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// GetQueryModel returns the well typed query model
func GetQueryModel(query backend.DataQuery) (*model.QueryModel, error) {
	model := &model.QueryModel{
		PageSize:    GaReportMaxResult,
		PageToken:   "",
		Offset:      0,
	}
	err := json.Unmarshal(query.JSON, &model)
	if err != nil {
		return nil, fmt.Errorf("error reading query: %s", err.Error())
	}

	// Copy directly from the well typed query
	timezone, err := time.LoadLocation(model.Timezone)
	if err != nil {
		return nil, fmt.Errorf("error get timezone %s", err.Error())
	}

	log.DefaultLogger.Debug("query timezone", "timezone", timezone.String())

	model.StartDate = query.TimeRange.From.In(timezone).Format("2006-01-02")
	model.EndDate = query.TimeRange.To.In(timezone).Format("2006-01-02")
	model.Dimensions = append([]string{model.TimeDimension}, model.Dimensions...)
	// model.TimeRange = query.TimeRange
	// model.MaxDataPoints = query.MaxDataPoints
	return model, nil
}
