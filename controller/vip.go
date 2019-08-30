package controller

import (
	"github.com/FTChinese/go-rest"
	"net/http"

	"github.com/FTChinese/go-rest/view"
)

// ListVIP lists all ftc account granted vip.
//
//	GET /vip?page=<number>&per_page=<number>
func (router ReaderRouter) ListVIP(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	vips, err := router.env.ListVIP(pagination)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().NoCache().SetBody(vips))
}

// GrantVIP grants vip to an ftc account.
//
//	PUT /vip/{id}
func (router ReaderRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	// 400 Bad Request
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err = router.env.GrantVIP(id); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	_ = view.Render(w, view.NewNoContent())
}

// RevokeVIP removes a ftc account from vip.
//
//	DELETE /vip/{id}
func (router ReaderRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	// 400 Bad Request
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.RevokeVIP(id); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	_ = view.Render(w, view.NewNoContent())
}
