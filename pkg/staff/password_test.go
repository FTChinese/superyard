package staff

import (
	"testing"
)

func TestNewPwResetSession(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		//want    PwResetSession
		wantErr bool
	}{
		{
			name:    "Create a new password reset session",
			args:    args{email: mockAccount.Email},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPwResetSession(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPwResetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewPwResetSession() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%+v", got)
		})
	}
}
