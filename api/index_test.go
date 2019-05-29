package api

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func TestIndexHandler(t *testing.T) {
	type args struct {
		app *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexHandler(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
