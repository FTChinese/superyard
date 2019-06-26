package model

import (
	"database/sql"
	"testing"

	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/android"
	"gitlab.com/ftchinese/backyard-api/test"
)

func TestAndroidEnv_CreateRelease(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		r android.Release
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "New Release",
			fields: fields{DB: test.DB},
			args: args{
				r: test.AndroidMock(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AndroidEnv{
				DB: tt.fields.DB,
			}
			if err := env.CreateRelease(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("AndroidEnv.CreateRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAndroidEnv_ListReleases(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "All Releases",
			fields: fields{
				DB: test.DB,
			},
			args: args{
				p: gorest.NewPagination(1, 20),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AndroidEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListReleases(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("AndroidEnv.ListReleases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("AndroidEnv.ListReleases() = %v, want %v", got, tt.want)
			//}

			t.Logf("Releases: %+v", got)
		})
	}
}

func TestAndroidEnv_SingleRelease(t *testing.T) {

	env := AndroidEnv{
		DB: test.DB,
	}

	release := test.AndroidMock()
	if err := env.CreateRelease(release); err != nil {
		panic(err)
	}

	type args struct {
		versionName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Retrieve a Single Release",
			args: args{
				versionName: release.VersionName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.SingleRelease(tt.args.versionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("AndroidEnv.SingleRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("AndroidEnv.SingleRelease() = %v, want %v", got, tt.want)
			//}

			t.Logf("A Single Release: %+v", got)
		})
	}
}

func TestAndroidEnv_UpdateRelease(t *testing.T) {
	env := AndroidEnv{
		DB: test.DB,
	}

	release := test.AndroidMock()
	if err := env.CreateRelease(release); err != nil {
		panic(err)
	}

	type args struct {
		r           android.Release
		versionName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update Release",
			args: args{
				r:           test.AndroidMock(),
				versionName: release.VersionName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.UpdateRelease(tt.args.r, tt.args.versionName); (err != nil) != tt.wantErr {
				t.Errorf("AndroidEnv.UpdateRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAndroidEnv_DeleteRelease(t *testing.T) {
	env := AndroidEnv{
		DB: test.DB,
	}

	release := test.AndroidMock()
	if err := env.CreateRelease(release); err != nil {
		panic(err)
	}

	type args struct {
		versionName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete a Release",
			args: args{
				versionName: release.VersionName,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DeleteRelease(tt.args.versionName); (err != nil) != tt.wantErr {
				t.Errorf("AndroidEnv.DeleteRelease() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("Deleted release: %s", tt.args.versionName)
		})
	}
}
