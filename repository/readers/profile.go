package readers

import "gitlab.com/ftchinese/superyard/pkg/reader"

// RetrieveFtcProfile loads profile of an email user.
func (env Env) RetrieveFtcProfile(ftcID string) (reader.FtcProfile, error) {
	var p reader.FtcProfileSchema

	if err := env.DB.Get(&p, reader.StmtFtcProfile, ftcID); err != nil {
		return reader.FtcProfile{}, err
	}

	return p.Build(), nil
}

// RetrieveWxProfile loads profile of a wx user.
func (env Env) RetrieveWxProfile(unionID string) (reader.WxProfile, error) {
	var p reader.WxProfile

	if err := env.DB.Get(&p, reader.StmtWxProfile, unionID); err != nil {
		return reader.WxProfile{}, err
	}

	return p, nil
}
