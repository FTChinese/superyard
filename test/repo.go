package test

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/superyard/pkg/android"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo() Repo {
	return Repo{
		db: DBX,
	}
}

func (repo Repo) CreateReader(r *Persona) error {
	query := `
	INSERT INTO cmstmp01.userinfo
	SET user_id = :ftc_id,
		wx_union_id = :wx_union_id,
		email = :email,
		password = MD5(:password),
		user_name = :user_name,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	if _, err := repo.db.NamedExec(query, r); err != nil {
		return err
	}

	return nil
}

func (repo Repo) MustCreateReader(r *Persona) {
	if err := repo.CreateReader(r); err != nil {
		panic(err)
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

	_, err := repo.db.NamedExec(query, info)

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

func (repo Repo) CreateStaff(s Staff) error {
	const query = `
	INSERT INTO backyard.staff
	SET staff_id = :staff_id,
		user_name = :user_name,
		email = :email,
		password = UNHEX(MD5(:password)),
		display_name = :display_name,
		department = :department,
		group_memberships = :group_memberships,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	_, err := repo.db.NamedExec(query, s)

	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) MustCreateStaff(s Staff) {
	err := repo.CreateStaff(s)

	if err != nil {
		panic(err)
	}
}

func (repo Repo) CreateAndroid(r android.Release) error {
	query := `INSERT INTO file_store.android_release
	SET version_name = :version_name,
		version_code = :version_code,
		body = :body,
		apk_url = :apk_url,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`
	_, err := repo.db.NamedExec(
		query,
		r)

	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) MustCreateAndroid(r android.Release) {
	err := repo.CreateAndroid(r)
	if err != nil {
		panic(err)
	}
}
