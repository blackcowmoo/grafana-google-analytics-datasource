package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"time"
)

func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.DefaultLogger.Info("Elapsed", what, time.Since(start).String())
	}
}

func ParseAndTimezoneTime(sTime string, timezone *time.Location) (*time.Time, error) {
	time, err := time.ParseInLocation("200601021504", padRightSide(sTime, "0", 12-len(sTime)), timezone)

	if err != nil {
		log.DefaultLogger.Info("timeConverter", "err", err)
		return nil, err
	}
	return &time, nil
}
