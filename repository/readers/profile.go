package readers

import "gitlab.com/ftchinese/superyard/models/reader"

// RetrieveFtcProfile loads profile of an email user.
func (env Env) RetrieveFtcProfile(ftcID string) (reader.FtcProfile, error) {
	var p reader.FtcProfile

	if err := env.DB.Get(&p, selectFtcProfile, ftcID); err != nil {
		return reader.FtcProfile{}, err
	}

	return p, nil
}

// RetrieveWxProfile loads profile of a wx user.
func (env Env) RetrieveWxProfile(unionID string) (reader.WxProfile, error) {
	var p reader.WxProfile

	if err := env.DB.Get(&p, selectWxProfile, unionID); err != nil {
		return reader.WxProfile{}, err
	}

	return p, nil
}
