package apps

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/android"
	"github.com/FTChinese/superyard/test"
	"testing"
)

func TestAndroidEnv_CreateRelease(t *testing.T) {
	r := test.NewRelease()

	t.Logf("Release: %+v", r)

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
			args: args{r: r},
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

func TestAndroidEnv_RetrieveRelease(t *testing.T) {
	r := test.NewRelease()
	test.NewRepo().MustCreateAndroid(r)

	env := AndroidEnv{DB: test.DBX}

	type args struct {
		versionName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve a release",
			args:    args{versionName: r.VersionName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveRelease(tt.args.versionName)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveRelease() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Release: %+v", got)
		})
	}
}

func TestAndroidEnv_UpdateRelease(t *testing.T) {
	r := test.NewRelease()
	test.NewRepo().MustCreateAndroid(r)

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
			name: "Update a release",
			args: args{r: r},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateRelease(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAndroidEnv_Exists(t *testing.T) {
	r := test.NewRelease()
	test.NewRepo().MustCreateAndroid(r)

	env := AndroidEnv{DB: test.DBX}

	type args struct {
		tag string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Release exists",
			args:    args{tag: r.VersionName},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.Exists(tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAndroidEnv_ListReleases(t *testing.T) {
	r := test.NewRelease()
	test.NewRepo().MustCreateAndroid(r)

	env := AndroidEnv{DB: test.DBX}

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "List releases",
			args:    args{p: gorest.NewPagination(1, 10)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListReleases(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListReleases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Release: %+v", got)
		})
	}
}

func TestAndroidEnv_DeleteRelease(t *testing.T) {
	r := test.NewRelease()
	test.NewRepo().MustCreateAndroid(r)

	env := AndroidEnv{DB: test.DBX}

	type args struct {
		versionName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Delete a release",
			args:    args{versionName: r.VersionName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DeleteRelease(tt.args.versionName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
