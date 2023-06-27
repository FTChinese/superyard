package readers

import "github.com/FTChinese/superyard/pkg/reader"

// RetrieveWxProfile loads profile of a wx user.
func (env Env) RetrieveWxProfile(unionID string) (reader.WxProfile, error) {
	var p reader.WxProfile

	err := env.gormDBs.Read.
		Where("union_id", unionID).
		First(&p).
		Error

	if err != nil {
		return reader.WxProfile{}, err
	}

	return p, nil
}
