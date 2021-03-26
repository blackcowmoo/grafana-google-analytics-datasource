package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func TestGoogleAnalytics_Query(t *testing.T) {
	type args struct {
		client *GoogleClient
		query  backend.DataQuery
	}
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		args    args
		want    *data.Frames
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.Query(tt.args.client, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalytics_GetAccounts(t *testing.T) {
	type args struct {
		ctx    context.Context
		config *DatasourceSettings
	}
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.GetAccounts(tt.args.ctx, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.GetAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.GetAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalytics_GetWebProperties(t *testing.T) {
	type args struct {
		ctx       context.Context
		config    *DatasourceSettings
		accountId string
	}
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.GetWebProperties(tt.args.ctx, tt.args.config, tt.args.accountId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.GetWebProperties() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.GetWebProperties() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalytics_GetProfiles(t *testing.T) {
	type args struct {
		ctx           context.Context
		config        *DatasourceSettings
		accountId     string
		webPropertyId string
	}
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		args    args
		want    map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.GetProfiles(tt.args.ctx, tt.args.config, tt.args.accountId, tt.args.webPropertyId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.GetProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.GetProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalytics_GetProfileTimezone(t *testing.T) {
	type args struct {
		ctx           context.Context
		config        *DatasourceSettings
		accountId     string
		webPropertyId string
		profileId     string
	}
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.GetProfileTimezone(tt.args.ctx, tt.args.config, tt.args.accountId, tt.args.webPropertyId, tt.args.profileId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.GetProfileTimezone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GoogleAnalytics.GetProfileTimezone() = %v, want %v", got, tt.want)
			}
		})
	}
}
