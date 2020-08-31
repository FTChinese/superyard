package reader

import (
	"errors"
	"github.com/guregu/null"
	"strings"
)

type IDs struct {
	FtcID   null.String `json:"ftcId" db:"ftc_id"`
	UnionID null.String `json:"unionId" db:"union_id"`
}

func (i IDs) GetCompoundID() (string, error) {
	if i.FtcID.Valid {
		return i.FtcID.String, nil
	}

	if i.UnionID.Valid {
		return i.UnionID.String, nil
	}

	return "", errors.New("neither ftc id nor union id provided")
}

func (i IDs) MustGetCompoundID() string {
	id, err := i.GetCompoundID()
	if err != nil {
		panic(err)
	}

	return id
}

// BuildFindInSet produces a value that can be using in FIND_IN_SET(col, value).
func (i IDs) BuildFindInSet() string {
	var ids []string

	if i.FtcID.Valid {
		ids = append(ids, i.FtcID.String)
	}

	if i.UnionID.Valid {
		ids = append(ids, i.UnionID.String)
	}

	return strings.Join(ids, ",")
}
