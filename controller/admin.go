package controller

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gitlab.com/ftchinese/backyard-api/staff"

	"gitlab.com/ftchinese/backyard-api/admin"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// AdminController performs adaministration tasks
type AdminController struct {
	adminModel admin.Env
	staffModel staff.Env
}

// NewStaff create a new account for a staff
// Input {
//	email: string,
//	userName: string,
//	displayName: string,
//	department: string,
//	groupMembers: int
// }
func (m AdminController) NewStaff(w http.ResponseWriter, req *http.Request) {
	var a staff.Account

	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &a); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	a.Sanitize()

	// { message: "Validation failed" | "The length of email should not exceed 20 chars" | "The length of userName should be within 1 to 20 chars" | "The length of displayName should be within 1 to 20 chars"
	//	 field: "email" | "userName" | "displayName",
	//	 code: "missing_field" | "invalid"
	// }
	if r := a.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := m.adminModel.NewStaff(a)

	// { message: "Validation failed",
	// 	field: "email | userName | displayName",
	//	code: "already_exists"
	// }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}

// StaffRoster lists all staff with pagination support
// TODO: add a middleware to parse form.
func (m AdminController) StaffRoster(w http.ResponseWriter, req *http.Request) {
	// err := req.ParseForm()

	// // 400 Bad Request
	// if err != nil {
	// 	view.Render(w, util.NewBadRequest(err.Error()))
	// 	return
	// }

	page := req.Form.Get("page")

	if page == "" {
		page = "1"
	}

	p, err := strconv.Atoi(page)

	// 400 Bad Request
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))
		return
	}

	accounts, err := m.adminModel.StaffRoster(p, 20)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(accounts))
}

// StaffProfile gets a staff's profile
func (m AdminController) StaffProfile(w http.ResponseWriter, req *http.Request) {
	userName := chi.URLParam(req, "name")

	// { message: "Invalid request URI" }
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	p, err := m.staffModel.Profile(userName)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(p))
}

// ReinstateStaff restore a previously deleted staff
func (m AdminController) ReinstateStaff(w http.ResponseWriter, req *http.Request) {
	userName := chi.URLParam(req, "name")

	// { message: "Invalid request URI" }
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.ActivateStaff(userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewNoContent())
}

// UpdateStaff updates a staff's profile
// Input {
//	email: string,
//	userName: string,
//	displayName: string,
//	department: string,
//	groupMembers: int
// }
func (m AdminController) UpdateStaff(w http.ResponseWriter, req *http.Request) {
	userName := chi.URLParam(req, "name")

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	var a staff.Account

	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &a); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	a.Sanitize()

	// { message: "Validation failed" | "The length of email should not exceed 20 chars" | "The length of userName should be within 1 to 20 chars" | "The length of displayName should be within 1 to 20 chars"
	//	 field: "email" | "userName" | "displayName",
	//	 code: "missing_field" | "invalid"
	// }
	if r := a.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := m.adminModel.UpdateStaff(userName, a)

	// { message: "Validation failed",
	// 	field: "email | userName | displayName",
	//	code: "already_exists"
	// }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
	}

	view.Render(w, util.NewNoContent())
}

// DeleteStaff flags a staff as inactive
// It also deletes all myft account associated with this staff;
// Unset vip of all related myft account;
// Remove all personal access token to access next-api;
// Remove all access token to access backyard-api
func (m AdminController) DeleteStaff(w http.ResponseWriter, req *http.Request) {
	userName := chi.URLParam(req, "name")

	// { message: "Invalid request URI" }
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.RemoveStaff(userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}

// VIPRoster lists all ftc account with vip set to true
func (m AdminController) VIPRoster(w http.ResponseWriter, req *http.Request) {
	myfts, err := m.adminModel.VIPRoster()

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(myfts))
}

// GrantVIP grants vip to a ftc account
func (m AdminController) GrantVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// { message: "Invalid request URI" }
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.GrantVIP(myftID)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}

// RevokeVIP removes a ftc account from vip
func (m AdminController) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// { message: "Invalid request URI" }
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.RevokeVIP(myftID)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}
