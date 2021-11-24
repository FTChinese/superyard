package test

import (
	"github.com/FTChinese/superyard/internal/pkg/android"
	"github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/staff"
	"github.com/FTChinese/superyard/pkg/subs"
)

type Repo struct {
	db db.ReadWriteMyDBs
}

func NewRepo() Repo {
	return Repo{
		db: DBX,
	}
}

func (repo Repo) CreateReader(a reader.FtcAccount) error {

	if _, err := repo.db.Write.NamedExec(reader.StmtCreateReader, a); err != nil {
		return err
	}

	return nil
}

func (repo Repo) MustCreateReader(a reader.FtcAccount) {
	if err := repo.CreateReader(a); err != nil {
		panic(err)
	}
}

func (repo Repo) CreateVIP(a reader.FtcAccount) error {
	_, err := repo.db.Write.NamedExec(reader.StmtCreateReader+", is_vip = :is_vip", a)
	if err != nil {
		return err
	}

	return nil
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

func (repo Repo) CreateOrder(order subs.Order) error {
	query := `INSERT INTO premium.ftc_trade
	SET trade_no = :order_id,
		user_id = :compound_id,
		ftc_user_id = :ftc_id,
		wx_union_id = :union_id,
		trade_price = :price,
		trade_amount = :amount,
		tier_to_buy = :tier,
		billing_cycle = :cycle,
		cycle_count = :cycle_count,
		extra_days = :extra_days,
		category = :kind,
		payment_method = :payment_method,
		created_utc = UTC_TIMESTAMP()`

	_, err := repo.db.Write.NamedExec(query, order)

	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) MustCreateOrder(order subs.Order) {
	err := repo.CreateOrder(order)

	if err != nil {
		panic(err)
	}
}
func (repo Repo) MustCreateMembership(m reader.Membership) {
	_, err := repo.db.Write.NamedExec(reader.StmtCreateMember, m)

	if err != nil {
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

// CreateAndroid inserts a new android release into db.
func (repo Repo) CreateAndroid(r android.Release) error {

	_, err := repo.db.Write.NamedExec(
		android.StmtInsertRelease,
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
