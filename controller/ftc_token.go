package controller

import (
	"net/http"

	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// NewToken creates an access token for a person or for an app.
//
//	POST /ftc-api/tokens
//
// Input
//	{
//		"description": "string", // optional, max 255 chars
//		"myftId": "string", // optional
//		"ownedByApp": "string" // optional
//	}
// `myftId` and `ownedByApp` should be mutually exclusive.
// If `ownedByApp` is present, it means this access token is created for an app. In such case `myftId` must be empty.
// If both `myftId` and `ownedByApp` are empty, it must be a personal access token.
//
// 400 Bad Request if request body cannot be parsed:
//	{
//		"message": "Problems parsing JSON"
//	}
//
// - `204 No Content` for success.
func (c FTCAPIRouter) NewToken(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var access ftcapi.APIKey

	// 400 Bad Request
	if err := util.Parse(req.Body, &access); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	access.Sanitize()
	if r := access.Validate(); r != nil {
		view.Render(w, util.NewUnprocessable(r))
		return
	}

	// Use userName from request header.
	access.CreatedBy = userName
	if access.MyftID != "" {
		access.OwnedByApp = ""
	} else if access.OwnedByApp != "" {
		access.MyftID = ""
	}

	err := c.apiModel.NewAPIKey(access)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content.
	view.Render(w, util.NewNoContent())
}

// PersonalTokens lists all access tokens created by a user.
//
//	GET /ftc-api/tokens/personal
//
// - 200 OK with body:
// 	[{
// 		"id": 1,
// 		"token": "40 hexdecimal numbers",
// 		"description": "",
// 		"myftId": "",
// 		"createdAt": "",
// 		"updatedAt": "",
// 		"lastUsedAt": ""
// 	}]
func (c FTCAPIRouter) PersonalTokens(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	keys, err := c.apiModel.PersonalAPIKeys(userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(keys))
}

// DeletePersonalToken deletes a personal access token.
//
//	DELETE /ftc-api/token/personal/{tokenId}
//
// - `400 Bad Request` if request URL does not contain `name` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - `204 No Content` for success.
func (c FTCAPIRouter) DeletePersonalToken(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	tokenID, err := getURLParam(req, "tokenID").toInt()
	// NOTE: id == 0 means remove all belong to this user
	if err != nil || tokenID < 1 {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err = c.apiModel.RemovePersonalAccess(tokenID, userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewNoContent())
}

// AppTokens show all access tokens used by an app.
//
//	GET /ftc-api/tokens/app/{name}
//
// - `400 Bad Request` if request URL does not contain `name` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - 200 OK with body:
// 	[{
// 		"id": 1,
// 		"token": "40 hexdecimal numbers",
// 		"description": "",
// 		"myftId": "",
// 		"createdAt": "",
// 		"updatedAt": "",
// 		"lastUsedAt": ""
// 	}]
func (c FTCAPIRouter) AppTokens(w http.ResponseWriter, req *http.Request) {
	// Get app name from url
	slugName := getURLParam(req, "name").toString()

	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	keys, err := c.apiModel.AppAPIKeys(slugName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// 204 No Content
	view.Render(w, util.NewResponse().NoCache().SetBody(keys))
}

// DeleteAppToken deletes an access token owned by an app
//
//	DELETE /ftc-api/tokens/app/{name}/{tokenId}
//
// - `400 Bad Request` if request URL does not contain `name` and `tokenId` part, or tokenId < 1.
//	{
//		"message": "Invalid request URI"
//	}
//
// - `204 No Content` for success.
func (c FTCAPIRouter) DeleteAppToken(w http.ResponseWriter, req *http.Request) {

	// Get app name from url
	slugName := getURLParam(req, "name").toString()
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	// Get token id from url
	tokenID, err := getURLParam(req, "tokenID").toInt()

	// NOTE: id == 0 means remove all belong to this user
	if err != nil || tokenID < 1 {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err = c.apiModel.RemoveAppAccess(tokenID, slugName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewNoContent())
}
