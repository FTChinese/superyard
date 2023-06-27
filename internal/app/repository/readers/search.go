package readers

import (
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) SearchReader(kw string, by reader.SearchBy) (reader.SearchResult, error) {
	stmt, err := reader.GetSearchStmt(by)
	if err != nil {
		return reader.SearchResult{}, err
	}

	var sr reader.SearchResult
	err = env.gormDBs.Read.Raw(stmt, kw).Scan(&sr).Error

	if err != nil {
		return reader.SearchResult{}, nil
	}

	return sr, nil
}

func (env Env) RetrieveAccount(id string, by reader.SearchBy) (reader.BaseAccount, error) {
	stmt, err := reader.GetAccountStmt(by)
	if err != nil {
		return reader.BaseAccount{}, err
	}
	var a reader.BaseAccount

	err = env.gormDBs.Read.
		Raw(stmt, id).
		Scan(&a).
		Error

	if err != nil {
		return reader.BaseAccount{}, err
	}

	return a, nil
}
