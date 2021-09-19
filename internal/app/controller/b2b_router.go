package controller

import "github.com/FTChinese/superyard/internal/app/repository/b2bapi"

type B2BRouter struct {
	apiClient b2bapi.B2BClient
}

func NewB2BRouter(prod bool) B2BRouter {
	return B2BRouter{
		apiClient: b2bapi.NewClient(prod),
	}
}
