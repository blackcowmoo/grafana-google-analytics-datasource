package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseAndTimezoneTime(t *testing.T) {
	localTimezone := time.Now().Local().Location()
	now := time.Now()
	now = now.Truncate(time.Minute)
	dateHourMinFormat := "200601021504"
	dateHourMin := now.Format(dateHourMinFormat)
	type args struct {
		sTime    string
		timezone *time.Location
	}
	tests := []struct {
		name    string
		args    args
		want    *time.Time
		wantErr bool
	}{
		{
			name: "dateHour",
			args: args{
				sTime:    dateHourMin,
				timezone: localTimezone,
			},
			wantErr: false,
			want:    &now,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAndTimezoneTime(tt.args.sTime, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAndTimezoneTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAndTimezoneTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
