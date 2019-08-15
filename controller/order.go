package controller

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type OrderRouter struct {
}

func NewOrderRouter(db *sqlx.DB) OrderRouter {
	return OrderRouter{}
}

func (router OrderRouter) ListOrders(w http.ResponseWriter, req *http.Request) {

}

func (router OrderRouter) CreateOrder(w http.ResponseWriter, req *http.Request) {

}

func (router OrderRouter) LoadOrder(w http.ResponseWriter, req *http.Request) {

}

func (router OrderRouter) UpdateOrder(w http.ResponseWriter, req *http.Request) {

}
