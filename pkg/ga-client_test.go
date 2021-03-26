package main

import (
	"context"
	"reflect"
	"testing"

	analytics "google.golang.org/api/analytics/v3"
	reporting "google.golang.org/api/analyticsreporting/v4"
)

func TestNewGoogleClient(t *testing.T) {
	type args struct {
		ctx  context.Context
		auth *DatasourceSettings
	}
	tests := []struct {
		name    string
		args    args
		want    *GoogleClient
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGoogleClient(tt.args.ctx, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGoogleClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGoogleClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createReportingService(t *testing.T) {
	type args struct {
		ctx  context.Context
		auth *DatasourceSettings
	}
	tests := []struct {
		name    string
		args    args
		want    *reporting.Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createReportingService(tt.args.ctx, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("createReportingService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createReportingService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createAnalyticsService(t *testing.T) {
	type args struct {
		ctx  context.Context
		auth *DatasourceSettings
	}
	tests := []struct {
		name    string
		args    args
		want    *analytics.Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAnalyticsService(tt.args.ctx, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAnalyticsService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createAnalyticsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleClient_getAccountsList(t *testing.T) {
	type args struct {
		idx int64
	}
	tests := []struct {
		name    string
		client  *GoogleClient
		args    args
		want    []*analytics.Account
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.getAccountsList(tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleClient.getAccountsList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleClient.getAccountsList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleClient_getAllWebpropertiesList(t *testing.T) {
	tests := []struct {
		name    string
		client  *GoogleClient
		want    []*analytics.Webproperty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.getAllWebpropertiesList()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleClient.getAllWebpropertiesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleClient.getAllWebpropertiesList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleClient_getWebpropertiesList(t *testing.T) {
	type args struct {
		accountId string
		idx       int64
	}
	tests := []struct {
		name    string
		client  *GoogleClient
		args    args
		want    []*analytics.Webproperty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.getWebpropertiesList(tt.args.accountId, tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleClient.getWebpropertiesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleClient.getWebpropertiesList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleClient_getAllProfilesList(t *testing.T) {
	tests := []struct {
		name    string
		client  *GoogleClient
		want    []*analytics.Profile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.getAllProfilesList()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleClient.getAllProfilesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleClient.getAllProfilesList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleClient_getProfilesList(t *testing.T) {
	type args struct {
		accountId     string
		webpropertyId string
		idx           int64
	}
	tests := []struct {
		name    string
		client  *GoogleClient
		args    args
		want    []*analytics.Profile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.getProfilesList(tt.args.accountId, tt.args.webpropertyId, tt.args.idx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleClient.getProfilesList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleClient.getProfilesList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleClient_getReport(t *testing.T) {
	type args struct {
		query QueryModel
	}
	tests := []struct {
		name    string
		client  *GoogleClient
		args    args
		want    *reporting.GetReportsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.getReport(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleClient.getReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleClient.getReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printResponse(t *testing.T) {
	type args struct {
		res *reporting.GetReportsResponse
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printResponse(tt.args.res)
		})
	}
}
