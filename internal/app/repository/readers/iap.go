package readers

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) IAPMember(origTxID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.dbs.Read.Get(&m, reader.StmtIAPMember, origTxID)

	if err != nil {
		if err == sql.ErrNoRows {
			return reader.Membership{}, nil
		}

		return reader.Membership{}, err
	}

	return m.Normalize(), nil
}
