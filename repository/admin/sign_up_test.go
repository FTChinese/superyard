package admin

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_Create(t *testing.T) {
	s := test.NewStaff()
	env := Env{DB: test.DBX}

	type args struct {
		a employee.SignUp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create a new staff",
			args:    args{a: s.SignUp()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.Create(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
