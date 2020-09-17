package android

import (
	"errors"
	"fmt"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/gh"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strconv"
	"strings"
)

const apkURL = "https://creatives.ftacademy.cn/minio/android/ftchinese-%s-ftc-release.apk"

// ReleaseInput contains the fields required to create or update a release.
// For creation, VersionName + VersionCode + ApkURL are required.
// For update, only ApkURL is required.
type ReleaseInput struct {
	// Required only when creating a new release. Ignore it when updating since the version name is
	// acquired from path parameter.
	VersionName string `json:"versionName" db:"version_name"`
	// Required only when creating a new release. This cannot be changed upon updating.
	VersionCode int64       `json:"versionCode" db:"version_code"`
	Body        null.String `json:"body" db:"body"`      // Optional
	ApkURL      string      `json:"apkUrl" db:"apk_url"` // Required.
}

func (r *ReleaseInput) ValidateUpdate() *render.ValidationError {
	r.Body.String = strings.TrimSpace(r.Body.String)
	r.ApkURL = strings.TrimSpace(r.ApkURL)

	return validator.New("apkUrl").
		Required().
		URL().
		Validate(r.ApkURL)
}

func (r *ReleaseInput) ValidateCreation() *render.ValidationError {
	r.VersionName = strings.TrimSpace(r.VersionName)

	ie := validator.New("versionName").
		Required().
		MaxLen(32).
		Validate(r.VersionName)
	if ie != nil {
		return ie
	}

	if r.VersionCode < 1 {
		return &render.ValidationError{
			Message: "version code must be larger than 0",
			Field:   "versionCode",
			Code:    render.CodeInvalid,
		}
	}

	return r.ValidateUpdate()
}

type Release struct {
	ReleaseInput
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt chrono.Time `json:"updatedAt" db:"updated_at"`
}

func NewRelease(input ReleaseInput) Release {
	return Release{
		ReleaseInput: input,
		CreatedAt:    chrono.TimeNow(),
		UpdatedAt:    chrono.TimeNow(),
	}
}

func FromGHRelease(r gh.Release, versionCode int64) Release {
	return Release{
		ReleaseInput: ReleaseInput{
			VersionName: r.TagName,
			VersionCode: versionCode,
			Body:        r.Body,
			ApkURL:      fmt.Sprintf(apkURL, r.TagName),
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.PublishedAt,
	}
}

type ReleaseList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []Release `json:"data"`
	Err  error     `json:"-"`
}

// ParseVersionCode gets the value of versionCode field
// from a gradle file.
func ParseVersionCode(content string) (int64, error) {
	codeStr := extractVersionCode(content)

	if codeStr == "" {
		return 0, errors.New("versionCode missing in gradle file")
	}

	versionCode, err := strconv.ParseInt(codeStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return versionCode, nil
}

// extractVersionCode get the versionName field from gradle file.
func extractVersionCode(str string) string {
	lines := strings.Split(str, "\n")

	var target string
	for _, v := range lines {
		if strings.Contains(v, "versionCode") {
			target = v
			break
		}
	}

	if target == "" {
		return ""
	}
	parts := strings.Split(strings.TrimSpace(target), " ")
	return strings.TrimSpace(parts[len(parts)-1])
}
