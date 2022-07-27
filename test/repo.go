package test

import (
	"github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/staff"
)

type Repo struct {
	db db.ReadWriteMyDBs
}

func NewRepo() Repo {
	return Repo{
		db: DBX,
	}
}

func (repo Repo) CreateWxInfo(info WxInfo) error {
	const query = `
	INSERT INTO user_db.wechat_userinfo
	SET union_id = :union_id,
		nickname = :nickname,
		avatar_url = :avatar,
		gender = :gender,
		country = :country,
		province = :province,
		city = :city,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	_, err := repo.db.Write.NamedExec(query, info)

	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) MustCreateWxInfo(info WxInfo) {
	if err := repo.CreateWxInfo(info); err != nil {
		panic(err)
	}
}

// CreateStaff inserts a new staff account into db.
func (repo Repo) CreateStaff(s staff.SignUp) error {
	_, err := repo.db.Write.NamedExec(staff.StmtCreateAccount, s)

	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) MustCreateStaff(s staff.SignUp) {
	err := repo.CreateStaff(s)

	if err != nil {
		panic(err)
	}
}

func (repo Repo) MustSavePwResetSession(session staff.PwResetSession) {
	_, err := repo.db.Write.NamedExec(staff.StmtInsertPwResetSession, session)

	if err != nil {
		panic(err)
	}
}

func (repo Repo) MustCreateOAuthApp(app oauth.App) {
	_, err := repo.db.Write.NamedExec(oauth.StmtInsertApp, app)

	if err != nil {
		panic(err)
	}
}

func (repo Repo) MustInsertAccessToken(t oauth.Access) {
	_, err := repo.db.Write.NamedExec(oauth.StmtInsertToken, t)

	if err != nil {
		panic(err)
	}
}
