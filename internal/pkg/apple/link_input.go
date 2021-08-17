package apple

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"strings"
)

// LinkInput defines the request body to link IAP to ftc account.
type LinkInput struct {
	FtcID        string `json:"ftcId"`
	OriginalTxID string `json:"originalTxId"` // Retrieved from URL path param.
}

func (i *LinkInput) Validate() *render.ValidationError {
	i.FtcID = strings.TrimSpace(i.FtcID)
	i.OriginalTxID = strings.TrimSpace(i.OriginalTxID)

	ve := validator.New("ftcId").Required().Validate(i.FtcID)
	if ve != nil {
		return ve
	}

	return validator.New("originalTxId").Required().Validate(i.OriginalTxID)
}
