package util

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func Elapsed(what string) func() {
	start := time.Now()
	return func() {
		log.DefaultLogger.Debug("Elapsed", what, time.Since(start).String())
	}
}

func padRightSide(str string, item string, count int) string {
	return str + strings.Repeat(item, count)
}

func ParseAndTimezoneTime(sTime string, timezone *time.Location) (*time.Time, error) {
	time, err := time.ParseInLocation("200601021504", padRightSide(sTime, "0", 12-len(sTime)), timezone)

	if err != nil {
		log.DefaultLogger.Error("timeConverter", "error", err)
		return nil, err
	}
	return &time, nil
}

func addTime(t1 time.Time, t2 time.Duration) time.Time {
	tmp := time.Time(t1)
	return tmp.Add(t2)
}

func subTime(t1 time.Time, t2 time.Duration) time.Time {
	return addTime(t1, t2*-1)
}

func SubOneHour(t1 time.Time) time.Time {
	return subTime(t1, time.Hour)
}

func SubOneDay(t1 time.Time) time.Time {
	return subTime(t1, time.Hour*24)
}

func SubOneMinute(t1 time.Time) time.Time {
	return subTime(t1, time.Minute)
}

func AddOneHour(t1 time.Time) time.Time {
	return addTime(t1, time.Hour)
}

func AddOneDay(t1 time.Time) time.Time {
	return addTime(t1, time.Hour*24)
}

func AddOneMinute(t1 time.Time) time.Time {
	return addTime(t1, time.Minute)
}

func FillArray(array []string, value string) []string {
	for i := range array {
		array[i] = value
	}
	return array
}

func TypeConverter[R any](data any) (*R, error) {
	var result R
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
