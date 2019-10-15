package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/oauth"
)

func (env Env) CreateToken(acc oauth.Access) (int64, error) {
	result, err := env.DB.NamedExec(stmtInsertToken, acc)
	if err != nil {
		logger.WithField("trace", "Env.CreateKey").Error(err)

		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (env Env) ListKeys(by oauth.KeySelector, p gorest.Pagination) ([]oauth.Access, error) {
	var keys = make([]oauth.Access, 0)

	var q string
	var v string
	switch {
	case by.ClientID.Valid:
		q = stmtAppKeys
		v = by.ClientID.String
	case by.StaffName.Valid:
		q = stmtPersonalKeys
		v = by.StaffName.String
	}
	err := env.DB.Select(&keys, q, v, p.Limit, p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListKeys").Error(err)
		return keys, err
	}

	return keys, nil
}

func (env Env) DeleteKeys(creator string) error {
	_, err := env.DB.Exec(stmtRemovePersonalKeys, creator)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) RemoveKey(by oauth.KeyRemover) error {

	_, err := env.DB.NamedExec(stmtRemoveKey, by)

	if err != nil {
		return err
	}

	return nil
}
