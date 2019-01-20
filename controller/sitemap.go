package controller

import (
	"net/http"

	"github.com/FTChinese/go-rest/view"
)

// Version show current version of api.
func Version(version, build string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		b := map[string]string{
			"version": version,
			"build":   build,
		}

		view.Render(w, view.NewResponse().NoCache().SetBody(b))
	}
}
