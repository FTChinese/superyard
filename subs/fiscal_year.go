package subs

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"time"
)

type FiscalYear struct {
	StartDate chrono.Date `json:"startDate"`
	LastDate  chrono.Date `json:"endDate"`
	Income    null.Float  `json:"income"`
}

func NewFiscalYear(year int) FiscalYear {
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(year, time.December, 31, 23, 59, 59, 59, time.UTC)

	return FiscalYear{
		StartDate: chrono.DateFrom(firstDay),
		LastDate:  chrono.DateFrom(lastDay),
	}
}
