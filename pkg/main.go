package main

import (
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	// window timezone issue #101
	_ "time/tzdata"
)

func main() {
	if err := datasource.Manage("blackcowmoo-googleanalytics-datasource", NewDataSource, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
