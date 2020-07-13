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

func TestSignUp_SignUpParcel(t *testing.T) {
	type fields struct {
		Account  Account
		Password string
	}
	type args struct {
		sourceURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Sign up email",
			fields: fields{
				Account:  mockAccount,
				Password: mockPassword,
			},
			args: args{
				sourceURL: "http://localhost:4200",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SignUp{
				Account:  tt.fields.Account,
				Password: tt.fields.Password,
			}
			got, err := s.SignUpParcel(tt.args.sourceURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignUpParcel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}
