package api

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestExtractFilters(t *testing.T) {
	type args struct {
		f url.Values
	}
	tests := []struct {
		name string
		args args
		want QueryData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractFilters(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryData_Build(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		q    *QueryData
		args args
		want *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Build(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryData.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
