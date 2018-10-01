package controller

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"

	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/validator"
	"gitlab.com/ftchinese/backyard-api/view"
)

type StaffController struct {
	model staff.Env
}

// Auth handles authentication process
// Input {userName: string, password: string, userIp: string}
func (s StaffController) Auth(w http.ResponseWriter, req *http.Request) {
	var login staff.Login

	if err := util.Parse(req.Body, &login); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	login.Sanitize()

	account, err := s.model.Auth(login)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
// Input {email: string}
func (s StaffController) ForgotPassword(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		view.Render(w, view.NewBadRequest(""))
		return
	}

	result := gjson.GetBytes(b, "email")

	if !result.Exists() {
		ue := view.UnprocessableError{
			Field: "email",
			Code:  view.CodeMissingField,
		}

		view.Render(w, view.NewUnprocessable("", ue))
	}

	email := strings.TrimSpace(result.String())

	if err := validator.Email(email); err != nil {
		view.Render(w, view.NewUnprocessable("", err))

		return
	}

}

func (s StaffController) VerifyToken(w http.ResponseWriter, req *http.Request) {

}

func (s StaffController) ResetPassword(w http.ResponseWriter, req *http.Request) {

}
