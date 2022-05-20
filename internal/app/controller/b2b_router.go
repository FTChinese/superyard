package controller

import "github.com/FTChinese/superyard/internal/app/repository/fta"

type B2BRouter struct {
	apiClient fta.Client
}

func NewB2BRouter(prod bool) B2BRouter {
	return B2BRouter{
		apiClient: fta.NewClient(prod),
	}
}
