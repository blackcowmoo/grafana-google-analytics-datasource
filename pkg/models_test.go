package main

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestGetQueryModel(t *testing.T) {
	type args struct {
		query backend.DataQuery
	}
	tests := []struct {
		name    string
		args    args
		want    *QueryModel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetQueryModel(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueryModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQueryModel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumnDefinition_GetType(t *testing.T) {
	tests := []struct {
		name string
		cd   *ColumnDefinition
		want ColumnType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cd.GetType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ColumnDefinition.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getColumnType(t *testing.T) {
	type args struct {
		headerType string
	}
	tests := []struct {
		name string
		args args
		want ColumnType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getColumnType(tt.args.headerType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getColumnType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewColumnDefinition(t *testing.T) {
	type args struct {
		header     string
		index      int
		headerType string
	}
	tests := []struct {
		name string
		args args
		want *ColumnDefinition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewColumnDefinition(tt.args.header, tt.args.index, tt.args.headerType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewColumnDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalytics_getMetadata(t *testing.T) {
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		want    *Metadata
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.getMetadata()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.getMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.getMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalytics_getFilteredMetadata(t *testing.T) {
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		want    []MetadataItem
		want1   []MetadataItem
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.ga.getFilteredMetadata()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.getFilteredMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.getFilteredMetadata() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GoogleAnalytics.getFilteredMetadata() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGoogleAnalytics_GetDimensions(t *testing.T) {
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		want    []MetadataItem
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.GetDimensions()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.GetDimensions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.GetDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoogleAnalytics_GetMetrics(t *testing.T) {
	tests := []struct {
		name    string
		ga      *GoogleAnalytics
		want    []MetadataItem
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ga.GetMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleAnalytics.GetMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleAnalytics.GetMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
