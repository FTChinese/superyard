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

func (env Env) ListKeys(owner oauth.KeyOwner, p gorest.Pagination) ([]oauth.Access, error) {
	var keys = make([]oauth.Access, 0)

	var q string
	switch owner.Usage {
	case oauth.KeyUsageApp:
		q = stmtAppKeys
	case oauth.KeyUsagePersonal:
		q = stmtPersonalKeys
	}
	err := env.DB.Select(&keys, q, owner.Value)

	if err != nil {
		logger.WithField("trace", "Env.ListKeys").Error(err)
		return keys, err
	}

	return keys, nil
}

func (env Env) RemoveKey(k oauth.Key) error {

	var err error
	switch k.Usage {
	case oauth.KeyUsageApp:
		_, err = env.DB.NamedExec(stmtRemoveAppKey, k)
	case oauth.KeyUsagePersonal:
		_, err = env.DB.NamedExec(stmtRemovePersonalKey, k)
	}

	if err != nil {
		return err
	}

	return nil
}
