package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/repository/customer"
	"net/http"
)

type MemberRouter struct {
	env customer.Env
}

func NewMemberRouter(db *sqlx.DB) MemberRouter {
	return MemberRouter{
		env: customer.Env{DB: db},
	}
}

func (router MemberRouter) ListMembers(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte("Not implemented"))
}

func (router MemberRouter) CreateMember(w http.ResponseWriter, req *http.Request) {
	log := logger.WithField("trace", "MemberRouter.CreateMember")

	var m reader.Membership
	if err := gorest.ParseJSON(req.Body, &m); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	m.GenerateID()

	if r := m.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := router.env.CreateMember(m); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	if err := view.Render(w, view.NewNoContent()); err != nil {
		log.Error(err)
	}
}

func (router MemberRouter) LoadMember(w http.ResponseWriter, req *http.Request) {
	log := logger.WithField("trace", "MemberRouter.CreateMember")

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	m, err := router.env.RetrieveMember(id)
	if err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(m))
}

func (router MemberRouter) UpdateMember(w http.ResponseWriter, req *http.Request) {
	log := logger.WithField("trace", "MemberRouter.UpdateMember")

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var m reader.Membership
	if err := gorest.ParseJSON(req.Body, &m); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	m.ID = null.StringFrom(id)

	if r := m.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := router.env.UpdateMember(m); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	if err := view.Render(w, view.NewNoContent()); err != nil {
		log.Error(err)
	}
}

func (router MemberRouter) DeleteMember(w http.ResponseWriter, req *http.Request) {
	log := logger.WithField("trace", "MemberRouter.UpdateMember")

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.DeleteMember(id); err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}
