package android

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/models/util"
	"strings"
)

type Release struct {
	VersionName string      `json:"versionName"`
	VersionCode int64       `json:"versionCode"`
	Body        null.String `json:"body"`
	ApkURL      string      `json:"apkUrl"`
	CreatedAt   chrono.Time `json:"createdAt"`
	UpdatedAt   chrono.Time `json:"updatedAt"`
}

func (r *Release) Sanitize() {
	r.VersionName = strings.TrimSpace(r.VersionName)
	r.Body.String = strings.TrimSpace(r.Body.String)
	r.ApkURL = strings.TrimSpace(r.ApkURL)
}

func (r Release) Validate() *view.Reason {
	if r.VersionCode < 1 {
		r := view.NewReason()
		r.Field = "versionCode"
		r.Code = view.CodeInvalid
		r.SetMessage("version code must be larger than 0")
		return r
	}

	if r := util.RequireNotEmptyWithMax(r.VersionName, 32, "versionName"); r != nil {
		return r
	}

	if r := util.RequireNotEmpty(r.ApkURL, "apkUrl"); r != nil {
		return r
	}

	return nil
}
