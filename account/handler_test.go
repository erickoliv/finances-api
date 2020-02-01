package account

import (
	"reflect"
	"testing"

	"github.com/erickoliv/finances-api/repository"
	"github.com/gin-gonic/gin"
)

func TestMakeAccountView(t *testing.T) {
	type args struct {
		repo repository.AccountService
	}
	tests := []struct {
		name string
		args args
		want AccountView
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeAccountView(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeAccountView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_Router(t *testing.T) {
	type fields struct {
		repo repository.AccountService
	}
	type args struct {
		group *gin.RouterGroup
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := handler{
				repo: tt.fields.repo,
			}
			view.Router(tt.args.group)
		})
	}
}

func Test_handler_GetAccounts(t *testing.T) {
	type fields struct {
		repo repository.AccountService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := handler{
				repo: tt.fields.repo,
			}
			view.GetAccounts(tt.args.c)
		})
	}
}

func Test_handler_CreateAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := handler{
				repo: tt.fields.repo,
			}
			view.CreateAccount(tt.args.c)
		})
	}
}

func Test_handler_GetAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := handler{
				repo: tt.fields.repo,
			}
			view.GetAccount(tt.args.c)
		})
	}
}

func Test_handler_UpdateAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := handler{
				repo: tt.fields.repo,
			}
			view.UpdateAccount(tt.args.c)
		})
	}
}

func Test_handler_DeleteAccount(t *testing.T) {
	type fields struct {
		repo repository.AccountService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := handler{
				repo: tt.fields.repo,
			}
			view.DeleteAccount(tt.args.c)
		})
	}
}
