package readers

import "github.com/FTChinese/superyard/pkg/reader"

// RetrieveWxProfile loads profile of a wx user.
func (env Env) RetrieveWxProfile(unionID string) (reader.WxProfile, error) {
	var p reader.WxProfile

	if err := env.dbs.Read.Get(&p, reader.StmtWxProfile, unionID); err != nil {
		return reader.WxProfile{}, err
	}

	return p, nil
}
