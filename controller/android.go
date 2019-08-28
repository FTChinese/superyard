package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/android"
	"gitlab.com/ftchinese/backyard-api/repository/apps"
	"net/http"
)

type AndroidRouter struct {
	model apps.AndroidEnv
}

func NewAndroidRouter(db *sqlx.DB) AndroidRouter {
	return AndroidRouter{
		model: apps.AndroidEnv{
			DB: db,
		},
	}
}

func (router AndroidRouter) TagExists(w http.ResponseWriter, req *http.Request) {
	versionName, err := GetURLParam(req, "versionName").ToString()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	ok, err := router.model.Exists(versionName)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	if !ok {
		_ = view.Render(w, view.NewNotFound())
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// CreateRelease inserts the metadata for a new Android release.
//
// POST /android/releases
//
// Body: {versionName: string, versionCode: int, body: string, apkUrl: string}
func (router AndroidRouter) CreateRelease(w http.ResponseWriter, req *http.Request) {
	var r android.Release

	if err := gorest.ParseJSON(req.Body, &r); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	r.Sanitize()

	if reason := r.Validate(); reason != nil {
		_ = view.Render(w, view.NewUnprocessable(reason))
		return
	}

	err := router.model.CreateRelease(r)
	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Field = "versionName"
			reason.Code = view.CodeAlreadyExists
			_ = view.Render(w, view.NewUnprocessable(reason))
			return
		}

		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// Releases retrieves all releases by sorting version code
// in descending order.
//
// GET /android/releases?page=<number>&per_page=<number>
func (router AndroidRouter) Releases(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	releases, err := router.model.ListReleases(pagination)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(releases))
}

// SingleReleases retrieves a release by version name
//
// GET /android/releases/{versionName}
func (router AndroidRouter) SingleRelease(w http.ResponseWriter, req *http.Request) {
	versionName, err := GetURLParam(req, "versionName").ToString()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	release, err := router.model.RetrieveRelease(versionName)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(release))
}

// UpdateRelease updates a single release.
//
// PATCH /android/releases/{versionName}
//
// Body {versionName: string, versionCode: int, body: string, binaryUrl: string}
func (router AndroidRouter) UpdateRelease(w http.ResponseWriter, req *http.Request) {
	versionName, err := GetURLParam(req, "versionName").ToString()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var release android.Release
	if err := gorest.ParseJSON(req.Body, &release); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	release.Sanitize()

	if r := release.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}
	release.VersionName = versionName

	err = router.model.UpdateRelease(release)

	if err != nil {
		if IsAlreadyExists(err) {
			r := view.NewReason()
			r.Field = "versionCode"
			r.Code = view.CodeAlreadyExists
			r.SetMessage("versionCode already exists")
			_ = view.Render(w, view.NewUnprocessable(r))
			return
		}

		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// DeleteRelease deletes a single release
//
// DELETE /android/releases/:versionName
func (router AndroidRouter) DeleteRelease(w http.ResponseWriter, req *http.Request) {
	versionName, err := GetURLParam(req, "versionName").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	err = router.model.DeleteRelease(versionName)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}
