package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// NewToken creates a access token for a person or for an app
// Input
// {
//	description: string,
//	myftId: string,
//	ownedByApp?: string
// }
func (c FTCController) NewToken(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var access ftcapi.APIKey

	if err := util.Parse(req.Body, &access); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	access.Sanitize()
	// TODO: validation

	access.CreatedBy = userName

	err := c.ftcModel.NewAPIKey(access)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}

// PersonalTokens lists all access tokens created by a user
func (c FTCController) PersonalTokens(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	keys, err := c.ftcModel.PersonalAPIKeys(userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(keys))
}

// UpdatePersonalToken updates a personal access token
// NOT Impelmented for now
// func (c FTCController) UpdatePersonalToken(w http.ResponseWriter, req *http.Request) {

// }

// RemovePersonalToken deletes a personal access token
func (c FTCController) RemovePersonalToken(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)
	tokenID := chi.URLParam(req, "tokenID")
	if tokenID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	id, err := strconv.Atoi(tokenID)

	// NOTE: id == 0 means remove all belong to this user
	if err != nil || id < 1 {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err = c.ftcModel.RemovePersonalAccess(id, userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewNoContent())
}

// AppTokens show all access tokens used by an app
func (c FTCController) AppTokens(w http.ResponseWriter, req *http.Request) {
	slugName := chi.URLParam(req, "name")

	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	keys, err := c.ftcModel.AppAPIKeys(slugName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(keys))
}

// UpdateAppToken updates an access token owned by an app
// func (c FTCController) UpdateAppToken(w http.ResponseWriter, req *http.Request) {

// }

// RemoveAppToken deletes an access token owned by an app
func (c FTCController) RemoveAppToken(w http.ResponseWriter, req *http.Request) {

	slugName := chi.URLParam(req, "name")
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	tokenID := chi.URLParam(req, "tokenID")
	if tokenID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	id, err := strconv.Atoi(tokenID)

	// NOTE: id == 0 means remove all belong to this user
	if err != nil || id < 1 {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err = c.ftcModel.RemoveAppAccess(id, slugName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewNoContent())
}
