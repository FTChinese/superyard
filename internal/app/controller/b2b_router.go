package controller

import "github.com/FTChinese/superyard/internal/app/repository/ftaapi"

type B2BRouter struct {
	apiClient ftaapi.FtaClient
}

func NewB2BRouter(prod bool) B2BRouter {
	return B2BRouter{
		apiClient: ftaapi.NewClient(prod),
	}
}
