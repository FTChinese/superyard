package apps

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/android"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestAndroidEnv_CreateRelease(t *testing.T) {
	env := AndroidEnv{DB: test.DBX}

	type args struct {
		r android.Release
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create release",
			args: args{
				r: android.Release{
					VersionName: "v4.0.0",
					VersionCode: 30,
					Body:        null.String{},
					ApkURL:      "https://creatives.ftacademy.cn/minio/android/ftchinese-v3.2.4-play-release.apk",
					CreatedAt:   chrono.TimeNow(),
					UpdatedAt:   chrono.TimeNow(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreateRelease(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("CreateRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
