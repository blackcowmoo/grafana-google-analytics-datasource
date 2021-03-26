package main

import (
	"reflect"
	"testing"
	"time"
)

func TestElapsed(t *testing.T) {
	type args struct {
		what string
	}
	tests := []struct {
		name string
		args args
		want func()
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Elapsed(tt.args.what); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Elapsed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAndTimezoneTime(t *testing.T) {
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
		// TODO: Add test cases.
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
