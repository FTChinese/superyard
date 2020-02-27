package android

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/validator"
	"strings"
)

type Release struct {
	VersionName string      `json:"versionName" db:"version_name"`
	VersionCode int64       `json:"versionCode" db:"version_code"`
	Body        null.String `json:"body" db:"body"`
	ApkURL      string      `json:"apkUrl" db:"apk_url"`
	CreatedAt   chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   chrono.Time `json:"updatedAt" db:"updated_at"`
}

func (r *Release) Sanitize() {
	r.VersionName = strings.TrimSpace(r.VersionName)
	r.Body.String = strings.TrimSpace(r.Body.String)
	r.ApkURL = strings.TrimSpace(r.ApkURL)
}

func (r Release) Validate() *render.ValidationError {
	if r.VersionCode < 1 {
		return &render.ValidationError{
			Message: "version code must be larger than 0",
			Field:   "versionCode",
			Code:    render.CodeInvalid,
		}
	}

	ie := validator.New("versionName").Required().Max(32).Validate(r.VersionName)
	if ie != nil {
		return ie
	}

	return validator.New("apkUrl").Required().URL().Validate(r.ApkURL)
}
