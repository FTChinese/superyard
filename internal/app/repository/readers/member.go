package readers

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) MemberByCompoundID(compoundID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.dbs.Read.Get(&m, reader.StmtSelectMember, compoundID)

	if err != nil {
		if err == sql.ErrNoRows {
			return reader.Membership{}, nil
		}
		return reader.Membership{}, err
	}

	return m.Normalize(), nil
}

type memberAsyncResult struct {
	value reader.Membership
	err   error
}

func (env Env) asyncAccountMember(compoundID string) <-chan memberAsyncResult {
	c := make(chan memberAsyncResult)

	go func() {
		m, err := env.MemberByCompoundID(compoundID)

		c <- memberAsyncResult{
			value: m,
			err:   err,
		}
	}()

	return c
}
