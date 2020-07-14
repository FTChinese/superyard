package android

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/validator"
	"strings"
)

type Release struct {
	VersionName string      `json:"versionName" db:"version_name"` // Required.
	VersionCode int64       `json:"versionCode" db:"version_code"` // Required
	Body        null.String `json:"body" db:"body"`                // Optional
	ApkURL      string      `json:"apkUrl" db:"apk_url"`           // Required.
	CreatedAt   chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   chrono.Time `json:"updatedAt" db:"updated_at"`
}

func (r *Release) Validate() *render.ValidationError {
	r.VersionName = strings.TrimSpace(r.VersionName)
	r.Body.String = strings.TrimSpace(r.Body.String)
	r.ApkURL = strings.TrimSpace(r.ApkURL)

	if r.VersionCode < 1 {
		return &render.ValidationError{
			Message: "version code must be larger than 0",
			Field:   "versionCode",
			Code:    render.CodeInvalid,
		}
	}

	ie := validator.New("versionName").
		Required().
		MaxLen(32).
		Validate(r.VersionName)
	if ie != nil {
		return ie
	}

	return validator.New("apkUrl").Required().URL().Validate(r.ApkURL)
}
