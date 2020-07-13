package staff

import (
	"github.com/FTChinese/go-rest/postoffice"
	"testing"
)

func TestAccount_PasswordResetParcel(t *testing.T) {

	type args struct {
		session PwResetSession
	}
	tests := []struct {
		name    string
		fields  Account
		args    args
		want    postoffice.Parcel
		wantErr bool
	}{
		{
			name:   "Password reset email",
			fields: mockAccount,
			args: args{
				session: MustNewPwResetSession(mockAccount.Email),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.fields
			got, err := a.PasswordResetParcel(tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswordResetParcel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("PasswordResetParcel() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%+v", got)
		})
	}
}
