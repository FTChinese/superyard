package controller

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type SubRouter struct {
}

func NewSubRouter(db *sqlx.DB) SubRouter {
	return SubRouter{}
}

func (router SubRouter) ListSubscriptions(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte("Not implemented"))
}

func (router SubRouter) CreateSubscription(w http.ResponseWriter, req *http.Request) {

}

func (router SubRouter) LoadSubscription(w http.ResponseWriter, req *http.Request) {

}

func (router SubRouter) UpdateSubscription(w http.ResponseWriter, req *http.Request) {

}

func (router SubRouter) DeleteSubscription(w http.ResponseWriter, req *http.Request) {

}
