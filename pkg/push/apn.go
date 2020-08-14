package push

import (
	"github.com/FTChinese/go-rest/chrono"
)

type MessageTeaser struct {
	ID               int64       `json:"id"`
	PageID           string      `json:"pageId"`
	Action           string      `json:"action"`
	Title            string      `json:"title"`
	ContentAvailable bool        `json:"contentAvailable"`
	CreatedBy        string      `json:"createdBy"`
	CreatedAt        chrono.Time `json:"createdAt"`
	DeviceCount      int64       `json:"deviceCount"`
	InvalidCount     int64       `json:"invalidCount"`
	TimeElapsed      int64       `json:"timeElapsed"`
}

type Message struct {
	MessageTeaser

	PushID string `json:"apnsId"`
	Body   string `json:"body"`
	Sound  bool   `json:"sound"`
	Media  string `json:"media"`
}
