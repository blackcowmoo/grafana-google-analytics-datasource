package main

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

func Test_transformReportToDataFrameByDimensions(t *testing.T) {
	type args struct {
		columns    []*ColumnDefinition
		report     *reporting.Report
		refId      string
		dimensions string
	}
	tests := []struct {
		name    string
		args    args
		want    *data.Frame
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformReportToDataFrameByDimensions(tt.args.columns, tt.args.report, tt.args.refId, tt.args.dimensions)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformReportToDataFrameByDimensions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformReportToDataFrameByDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformReportToDataFrames(t *testing.T) {
	type args struct {
		report   *reporting.Report
		refId    string
		timezone string
	}
	tests := []struct {
		name    string
		args    args
		want    []*data.Frame
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformReportToDataFrames(tt.args.report, tt.args.refId, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformReportToDataFrames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformReportToDataFrames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformReportsResponseToDataFrames(t *testing.T) {
	type args struct {
		reportsResponse *reporting.GetReportsResponse
		refId           string
		timezone        string
	}
	tests := []struct {
		name    string
		args    args
		want    *data.Frames
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformReportsResponseToDataFrames(tt.args.reportsResponse, tt.args.refId, tt.args.timezone)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformReportsResponseToDataFrames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformReportsResponseToDataFrames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_padRightSide(t *testing.T) {
	type args struct {
		str   string
		item  string
		count int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := padRightSide(tt.args.str, tt.args.item, tt.args.count); got != tt.want {
				t.Errorf("padRightSide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getColumnDefinitions(t *testing.T) {
	type args struct {
		header *reporting.ColumnHeader
	}
	tests := []struct {
		name string
		args args
		want []*ColumnDefinition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getColumnDefinitions(tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getColumnDefinitions() = %v, want %v", got, tt.want)
			}
		})
	}
}
