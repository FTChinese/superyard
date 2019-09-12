package reader

import (
	"errors"
	"github.com/guregu/null"
	"strings"
)

type AccountID struct {
	CompoundID string      `json:"compoundId" db:"compound_id"`
	FtcID      null.String `json:"ftcId" db:"ftc_id"`
	UnionID    null.String `json:"unionId" db:"union_id"`
}

// NewAccountID creates a new User instance and select the correct CompoundID
func NewAccountID(ftcID string, unionID string) AccountID {
	ftcID = strings.TrimSpace(ftcID)
	unionID = strings.TrimSpace(unionID)

	id := AccountID{
		FtcID:   null.NewString(ftcID, ftcID != ""),
		UnionID: null.NewString(unionID, unionID != ""),
	}

	if ftcID != "" {
		id.CompoundID = id.FtcID.String
	} else if unionID != "" {
		id.CompoundID = id.UnionID.String
	}

	return id
}

// QueryArgs turns account id to a slice to be used to
// builder the SQL `in` statement.
func (id AccountID) QueryArgs() []interface{} {
	var args []interface{}

	if id.FtcID.Valid {
		args = append(args, id.FtcID.String)
	}

	if id.UnionID.Valid {
		args = append(args, id.UnionID.String)
	}

	return args
}

// SetCompoundID select a compound id.
func (id *AccountID) SetCompoundID() error {
	if id.FtcID.IsZero() && id.UnionID.IsZero() {
		return errors.New("one of ftc id or union id should be present")
	}

	if id.FtcID.Valid {
		id.CompoundID = id.FtcID.String
		return nil
	}

	id.CompoundID = id.UnionID.String
	return nil
}
