package stats

import "github.com/FTChinese/go-rest/chrono"

// SignUp calculates how many new users signed up every day
type SignUp struct {
	Count int         `json:"count"`
	Date  chrono.Date `json:"date"`
}
