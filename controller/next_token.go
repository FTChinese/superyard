package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/oauth"
	"gitlab.com/ftchinese/backyard-api/util"
	"net/http"
)

// NewToken creates an access token for a person or for an app.
//
//	POST /next/apps/{name}/tokens
//
// Input: {description: string}
func (router NextAPIRouter) NewAppToken(w http.ResponseWriter, req *http.Request) {

	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	clientID, err := router.model.FindClientID(slugName)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	acc, err := oauth.NewAccess()
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}
	if err := gorest.ParseJSON(req.Body, &acc); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	_, err = router.model.SaveAppAccess(acc, clientID)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// ListAppTokens show all access tokens used by an app. Header `X-Staff-Name`
//
//	GET /next/apps/{name}/tokens?page=<number>
func (router NextAPIRouter) ListAppTokens(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	page, _ := GetQueryParam(req, "page").ToInt()
	pagination := util.NewPagination(page, 20)

	tokens, err := router.model.ListAppAccess(slugName, pagination)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(tokens))
}

// DeleteAppToken deletes an access token owned by an app
//
//	DELETE /next/apps/{name}/tokens/{id}
func (router NextAPIRouter) RemoveAppToken(w http.ResponseWriter, req *http.Request) {
	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	tokenID, err := GetURLParam(req, "id").ToInt()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if tokenID <= 0 {
		view.Render(w, view.NewBadRequest("Token id must be larger than 0"))
	}

	clientID, err := router.model.FindClientID(slugName)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	err = router.model.RemoveAppAccess(clientID, tokenID)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// CreateKey creates a personal access token. Header `X-User-Name`.
//
//	POST /next/keys
//
// Input: {description: string, myftEmail: string}
func (router NextAPIRouter) CreateKey(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	acc, err := oauth.NewPersonalAccess()
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}
	if err := gorest.ParseJSON(req.Body, &acc); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	token, err := oauth.NewToken()
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	acc.Token = token
	acc.CreatedBy = null.StringFrom(userName)

	myftID := router.model.FindMyftID(acc.MyftEmail)

	_, err = router.model.SavePersonalToken(acc, myftID)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// ListKeys shows all active personal access tokens a user created. Header `X-User-Name`.
//
//	GET /next/keys?page=<number>
func (router NextAPIRouter) ListKeys(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	err := req.ParseForm()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pageNum, _ := GetQueryParam(req, "page").ToInt()
	pagination := util.NewPagination(pageNum, 20)

	acc, err := router.model.ListPersonalTokens(userName, pagination)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(acc))
}

// RemoveKey deactivate a personal access token owned by a user. Header `X-User-Name`.
//
// DELETE /next/keys/{id}
func (router NextAPIRouter) RemoveKey(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	id, err := GetURLParam(req, "tokenId").ToInt()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
	}

	err = router.model.RemovePersonalToken(userName, id)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}
