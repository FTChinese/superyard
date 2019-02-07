package subs

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/util"
	"strings"
)

// Schedule represents the beginning and ending time of
// a promotion event.
type Schedule struct {
	ID          int64       `json:"id"`
	Name        string      `json:"name"`        // Required. Max 256 chars
	Description null.String `json:"description"` // Optional. Max 256 chars
	StartAt     chrono.Time `json:"startAt"`     // Required. ISO 8601 date time string.
	EndAt       chrono.Time `json:"endAt"`       // Required. ISO 8601 date time string.
}

// Sanitize removes leading and trailing spaces of each string fields.
func (s *Schedule) Sanitize() {
	s.Name = strings.TrimSpace(s.Name)
	if s.Description.Valid {
		s.Description.String = strings.TrimSpace(s.Description.String)
	}
}

// Validate validates incoming data for a new schedule.
func (s *Schedule) Validate() *view.Reason {
	if r := util.RequireNotEmptyWithMax(s.Name, 256, "name"); r != nil {
		return r
	}

	if r := util.OptionalMaxLen(s.Description.String, 256, "description"); r != nil {
		return r
	}

	return nil
}

// Promotion contains all data for a promotion campaign.
// A pomotion compaign is divided into three steps:
// Schedule the time when it will start and end;
// Pricing plans for each products;
// Banner content used for pomotion.
// It is created not in one shot, but step by step.
// When retrieving, all data are retrieved together.
// When deleting, everything is deleted for a promotion record.
// When updating, schedule and plans parts are not allowed to edit;
// but banner content is editable.Promotion
type Promotion struct {
	Schedule
	Plans     Pricing     `json:"plans"`
	Banner    Banner     `json:"banner"`
	IsEnabled bool        `json:"isEnabled"`
	CreatedAt chrono.Time `json:"createdAt"`
	UpdatedAt chrono.Time `json:"updatedAt"`
	CreatedBy string      `json:"createdBy"`
}






