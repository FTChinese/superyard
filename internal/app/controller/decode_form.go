package controller

type LiveRefresh struct {
	Live    bool `schema:"live"`
	Refresh bool `schema:"refresh"`
}
