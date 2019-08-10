package stats

import (
	"testing"
)

func TestNewPeriod(t *testing.T) {
	type args struct {
		start string
		end   string
	}
	tests := []struct {
		name    string
		args    args
		want    Period
		wantErr bool
	}{
		{
			name:    "Default Period",
			args:    args{"", ""},
			wantErr: false,
		},
		{
			name:    "Reversed",
			args:    args{"2019-02-11", "2019-02-07"},
			wantErr: false,
		},
		{
			name:    "Start empty",
			args:    args{"", "2019-2-11"},
			wantErr: false,
		},
		{
			name:    "End empty",
			args:    args{"2019-02-10", ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPeriod(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("NewPeriod() = %v, want %v", got, tt.want)
			//}

			t.Logf("%+v", got)
		})
	}
}
