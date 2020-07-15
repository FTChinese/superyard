package staff

import (
	"github.com/guregu/null"
	"testing"
)

var mockAccount = Account{
	ID:           null.StringFrom("stf_NN4sA8TmYDGO"),
	UserName:     "Feeney9284",
	Email:        "ephraimbosco@gibson.org",
	DisplayName:  null.StringFrom("Francisco Crona"),
	Department:   null.StringFrom("tech"),
	GroupMembers: 2,
	IsActive:     true,
}

var mockPassword = "tb2lo13m"

func TestGenStaffID(t *testing.T) {
	t.Logf("Generate a staff id: %s", GenStaffID())
}

func TestNewSignUp(t *testing.T) {
	type args struct {
		input InputData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "New sign up",
			args: args{
				input: InputData{
					SignUp: SignUp{
						Account:  mockAccount,
						Password: "tb2lo13m",
					},
					SourceURL: "http://localhost:4200",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSignUp(tt.args.input)

			if got.ID == mockAccount.ID {
				t.Error("NewSignUp() ID not generated")
			}
		})
	}
}
