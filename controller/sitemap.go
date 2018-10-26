package controller

import (
	"net/http"

	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// Version show current version of api.
func Version(version, build string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		b := map[string]string{
			"version": version,
			"build":   build,
		}

		view.Render(w, util.NewResponse().NoCache().SetBody(b))
	}
}
