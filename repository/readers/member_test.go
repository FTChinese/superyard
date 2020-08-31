package readers

import (
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/test"
	"testing"
)

func TestEnv_CreateMember(t *testing.T) {

	env := Env{db: test.DBX}

	type args struct {
		m reader.Membership
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create Member",
			args: args{m: test.NewPersona().Membership()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreateMember(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("CreateMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
