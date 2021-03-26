package main

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestLoadSettings(t *testing.T) {
	type args struct {
		ctx backend.PluginContext
	}
	tests := []struct {
		name    string
		args    args
		want    *DatasourceSettings
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadSettings(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}
