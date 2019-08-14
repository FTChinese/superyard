package customer

import (
	"fmt"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

func (env Env) RetrieveMember(col MemberColumn, val string) (reader.Membership, error) {
	var m reader.Membership

	stmt := fmt.Sprintf(stmtMember, string(col))
	err := env.DB.Get(&m, stmt, val)

	if err != nil {
		logger.WithField("trace", "Env.LoadMember").Error(err)
		return m, err
	}

	return m, nil
}
