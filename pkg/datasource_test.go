package main

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestNewDataSource(t *testing.T) {
	type args struct {
		mux *http.ServeMux
	}
	tests := []struct {
		name string
		args args
		want *GoogleAnalyticsDataSource
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataSource(tt.args.mux); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalyticsDataSource_CheckHealth(t *testing.T) {
	type args struct {
		ctx context.Context
		req *backend.CheckHealthRequest
	}
	tests := []struct {
		name    string
		ds      *GoogleAnalyticsDataSource
		args    args
		want    *backend.CheckHealthResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ds.CheckHealth(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalyticsDataSource.CheckHealth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalyticsDataSource.CheckHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalyticsDataSource_QueryData(t *testing.T) {
	type args struct {
		ctx context.Context
		req *backend.QueryDataRequest
	}
	tests := []struct {
		name    string
		ds      *GoogleAnalyticsDataSource
		args    args
		want    *backend.QueryDataResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ds.QueryData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalyticsDataSource.QueryData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalyticsDataSource.QueryData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeResult(t *testing.T) {
	type args struct {
		rw   http.ResponseWriter
		path string
		val  interface{}
		err  error
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeResult(tt.args.rw, tt.args.path, tt.args.val, tt.args.err)
		})
	}
}

func TestGoogleAnalyticsDataSource_handleResourceAccounts(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		ds   *GoogleAnalyticsDataSource
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ds.handleResourceAccounts(tt.args.rw, tt.args.req)
		})
	}
}

func TestGoogleAnalyticsDataSource_handleResourceWebProperties(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		ds   *GoogleAnalyticsDataSource
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ds.handleResourceWebProperties(tt.args.rw, tt.args.req)
		})
	}
}

func TestGoogleAnalyticsDataSource_handleResourceProfiles(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		ds   *GoogleAnalyticsDataSource
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ds.handleResourceProfiles(tt.args.rw, tt.args.req)
		})
	}
}

func TestGoogleAnalyticsDataSource_handleResourceDimensions(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		ds   *GoogleAnalyticsDataSource
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ds.handleResourceDimensions(tt.args.rw, tt.args.req)
		})
	}
}

func TestGoogleAnalyticsDataSource_handleResourceMetrics(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		ds   *GoogleAnalyticsDataSource
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ds.handleResourceMetrics(tt.args.rw, tt.args.req)
		})
	}
}

func TestGoogleAnalyticsDataSource_handleResourceProfileTimezone(t *testing.T) {
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		ds   *GoogleAnalyticsDataSource
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ds.handleResourceProfileTimezone(tt.args.rw, tt.args.req)
		})
	}
}
