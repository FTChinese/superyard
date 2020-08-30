package reader

import (
	"errors"
	"github.com/guregu/null"
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
