package staff

import (
	"github.com/FTChinese/go-rest/render"
	"reflect"
	"testing"
)

func TestInputData_ValidateLogin(t *testing.T) {
	type fields struct {
		SignUp      SignUp
		OldPassword string
		Token       string
		SourceURL   string
	}
	tests := []struct {
		name   string
		fields fields
		want   *render.ValidationError
	}{
		{
			name: "Username + Password",
			fields: fields{
				SignUp: SignUp{
					Account: Account{
						UserName: mockAccount.UserName,
					},
					Password: "12345678",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InputData{
				SignUp:      tt.fields.SignUp,
				OldPassword: tt.fields.OldPassword,
				Token:       tt.fields.Token,
				SourceURL:   tt.fields.SourceURL,
			}
			if got := i.ValidateLogin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputData_ValidateEmail(t *testing.T) {
	type fields struct {
		SignUp      SignUp
		OldPassword string
		Token       string
		SourceURL   string
	}
	tests := []struct {
		name   string
		fields fields
		want   *render.ValidationError
	}{
		{
			name: "Valid email",
			fields: fields{
				SignUp: SignUp{
					Account: Account{
						Email: "valid@ftchinese.com",
					},
				},
			},
			want: nil,
		},
		{
			name: "Invalid email",
			fields: fields{
				SignUp: SignUp{
					Account: Account{
						Email: "georgiannahayes@morissette.biz",
					},
				},
			},
			want: &render.ValidationError{
				Message: "Email must be owned by ftchinese",
				Field:   "email",
				Code:    render.CodeInvalid,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InputData{
				SignUp:      tt.fields.SignUp,
				OldPassword: tt.fields.OldPassword,
				Token:       tt.fields.Token,
				SourceURL:   tt.fields.SourceURL,
			}
			if got := i.ValidateEmail(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputData_ValidatePasswordReset(t *testing.T) {
	type fields struct {
		SignUp      SignUp
		OldPassword string
		Token       string
		SourceURL   string
	}
	tests := []struct {
		name   string
		fields fields
		want   *render.ValidationError
	}{
		{
			name: "Token + Password",
			fields: fields{
				SignUp: SignUp{
					Password: "12345678",
				},
				Token: "af39f4a49c9c5f8e04f80433b66abf2ccfe350d9dc88a77b9f117a9f5f537ebf",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InputData{
				SignUp:      tt.fields.SignUp,
				OldPassword: tt.fields.OldPassword,
				Token:       tt.fields.Token,
				SourceURL:   tt.fields.SourceURL,
			}
			if got := i.ValidatePasswordReset(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidatePasswordReset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputData_ValidatePwUpdater(t *testing.T) {
	type fields struct {
		SignUp      SignUp
		OldPassword string
		Token       string
		SourceURL   string
	}
	tests := []struct {
		name   string
		fields fields
		want   *render.ValidationError
	}{
		{
			name: "OldPassword + Password",
			fields: fields{
				SignUp: SignUp{
					Password: "87654321",
				},
				OldPassword: "12345678",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InputData{
				SignUp:      tt.fields.SignUp,
				OldPassword: tt.fields.OldPassword,
				Token:       tt.fields.Token,
				SourceURL:   tt.fields.SourceURL,
			}
			if got := i.ValidatePwUpdater(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidatePwUpdater() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputData_ValidateAccount(t *testing.T) {
	type fields struct {
		SignUp      SignUp
		OldPassword string
		Token       string
		SourceURL   string
	}
	tests := []struct {
		name   string
		fields fields
		want   *render.ValidationError
	}{
		{
			name: "Validate updating account by admin",
			fields: fields{
				SignUp: SignUp{
					Account: Account{
						UserName:     mockAccount.UserName,
						Email:        "example@ftchinese.com",
						DisplayName:  mockAccount.DisplayName,
						Department:   mockAccount.Department,
						GroupMembers: mockAccount.GroupMembers,
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InputData{
				SignUp:      tt.fields.SignUp,
				OldPassword: tt.fields.OldPassword,
				Token:       tt.fields.Token,
				SourceURL:   tt.fields.SourceURL,
			}
			if got := i.ValidateAccount(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputData_ValidateSignUp(t *testing.T) {
	type fields struct {
		SignUp      SignUp
		OldPassword string
		Token       string
		SourceURL   string
	}
	tests := []struct {
		name   string
		fields fields
		want   *render.ValidationError
	}{
		{
			name: "Validate updating account by admin",
			fields: fields{
				SignUp: SignUp{
					Account: Account{
						UserName:     mockAccount.UserName,
						Email:        "example@ftchinese.com",
						DisplayName:  mockAccount.DisplayName,
						Department:   mockAccount.Department,
						GroupMembers: mockAccount.GroupMembers,
					},
					Password: "12345678",
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InputData{
				SignUp:      tt.fields.SignUp,
				OldPassword: tt.fields.OldPassword,
				Token:       tt.fields.Token,
				SourceURL:   tt.fields.SourceURL,
			}
			if got := i.ValidateSignUp(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateSignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}
