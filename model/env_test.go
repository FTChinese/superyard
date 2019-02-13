package model

import (
	"database/sql"
	"strings"
	"testing"
	"time"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"github.com/icrowley/fake"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/ftchinese/backyard-api/oauth"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/subs"
	"gitlab.com/ftchinese/backyard-api/user"
)

func newDevDB() *sql.DB {
	db, err := sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")

	if err != nil {
		panic(err)
	}

	return db
}

func genPassword() string {
	return fake.Password(8, 20, false, true, false)
}

func genUnionID() string {
	id, _ := gorest.RandomBase64(21)
	return id
}

func genOrderID() string {
	id, _ := gorest.RandomHex(8)

	return "FT" + strings.ToUpper(id)
}

func clientApp() gorest.ClientApp {
	return gorest.ClientApp{
		ClientType: enum.PlatformAndroid,
		Version:    "1.1.1",
		UserIP:     fake.IPv4(),
		UserAgent:  fake.UserAgent(),
	}
}

var db = newDevDB()

var defaultPlans = subs.Pricing{
	"standard_year": subs.Plan{
		Tier:        enum.TierStandard,
		Cycle:       enum.CycleYear,
		ListPrice:   258.00,
		NetPrice:    258.00,
		Description: "FT中文网 - 年度标准会员",
	},
	"standard_month": subs.Plan{
		Tier:        enum.TierStandard,
		Cycle:       enum.CycleMonth,
		ListPrice:   28.00,
		NetPrice:    28.00,
		Description: "FT中文网 - 月度标准会员",
	},
	"premium_year": subs.Plan{
		Tier:        enum.TierPremium,
		Cycle:       enum.CycleYear,
		ListPrice:   1998.00,
		NetPrice:    1998.00,
		Description: "FT中文网 - 高端会员",
	},
}

var stdPlan = defaultPlans["standard_year"]

var oauthEnv = OAuthEnv{DB: db}

type mockStaff struct {
	email        string
	userName     string
	password     string
	displayName  string
	department   string
	groupMembers int64
	loop         int
}

func newMockStaff() mockStaff {
	return mockStaff{
		email:        fake.EmailAddress(),
		userName:     fake.UserName(),
		password:     genPassword(),
		displayName:  fake.FullName(),
		department:   "tech",
		groupMembers: 3,
		loop:         5,
	}
}

func (m mockStaff) account() staff.Account {
	a, _ := staff.NewAccount()

	a.Email = m.email
	a.UserName = m.userName
	a.DisplayName = null.StringFrom(m.displayName)
	a.Department = null.StringFrom(m.department)
	a.GroupMembers = m.groupMembers

	return a
}

func (m mockStaff) login() staff.Login {
	return staff.Login{
		UserName: m.userName,
		Password: m.password,
	}
}

func (m mockStaff) newPassword() staff.Password {
	return staff.Password{
		Old: m.password,
		New: genPassword(),
	}
}

func (m mockStaff) schedule() subs.Schedule {
	return subs.Schedule{
		Name:        fake.ProductName(),
		Description: null.StringFrom(fake.Sentence()),
		StartAt:     chrono.TimeNow(),
		EndAt:       chrono.TimeFrom(time.Now().AddDate(0, 0, 1)),
	}
}

func (m mockStaff) banner() subs.Banner {
	return subs.Banner{
		CoverURL:   "http://www." + fake.Word() + ".com/" + fake.Word() + ".png",
		Heading:    fake.Title(),
		SubHeading: fake.Title(),
		Content: []string{
			fake.Paragraph(),
			fake.Paragraph(),
		},
	}
}

func (m mockStaff) app() oauth.App {
	name := fake.ProductName()
	slug := strings.Replace(name, " ", "-", -1)
	app := oauth.App{
		Name:        name,
		Slug:        strings.ToLower(slug),
		RepoURL:     "https://githu.com/FTChinese/" + fake.Word(),
		Description: null.StringFrom(fake.Sentence()),
		HomeURL:     null.StringFrom("http://www.ftchinese.com/" + fake.Word()),
		OwnedBy:     m.userName,
	}

	err := app.GenCredentials()
	if err != nil {
		panic(err)
	}

	return app
}

func (m mockStaff) personalAccess() oauth.PersonalAccess {
	acc := oauth.PersonalAccess{
		Description: null.StringFrom(fake.Sentence()),
		CreatedBy:   null.StringFrom(m.userName),
	}

	token, _ := oauth.NewToken()

	acc.Token = token

	return acc
}

func (m mockStaff) createAccount() staff.Account {
	adminEnv := AdminEnv{DB: db}

	a := m.account()

	err := adminEnv.CreateAccount(a)

	if err != nil {
		panic(err)
	}

	return a
}

func (m mockStaff) createMyft(u user.User) staff.Myft {
	m.createAccount()

	staffEnv := StaffEnv{DB: db}

	myft := staff.Myft{
		StaffName: m.userName,
		MyftID:    u.UserID,
	}

	err := staffEnv.saveMyft(myft)

	if err != nil {
		panic(err)
	}

	return myft
}

func (m mockStaff) createPwToken() staff.TokenHolder {
	a := m.createAccount()
	th, err := a.TokenHolder()
	if err != nil {
		panic(err)
	}

	err = StaffEnv{DB: db}.SavePwResetToken(th)
	if err != nil {
		panic(err)
	}

	return th
}

func (m mockStaff) createSchedule() int64 {
	sch := m.schedule()
	id, err := PromoEnv{DB: db}.NewSchedule(sch, m.userName)
	if err != nil {
		panic(err)
	}

	return id
}

func (m mockStaff) createPromo() int64 {
	sch := m.schedule()
	env := PromoEnv{DB: db}
	id, err := env.NewSchedule(sch, m.userName)
	if err != nil {
		panic(err)
	}

	err = env.SavePlans(id, defaultPlans)
	if err != nil {
		panic(err)
	}

	err = env.SaveBanner(id, m.banner())
	if err != nil {
		panic(err)
	}

	return id
}

func (m mockStaff) createApp() oauth.App {
	m.createAccount()
	app := m.app()

	err := oauthEnv.SaveApp(app)
	if err != nil {
		panic(err)
	}

	return app
}

func (m mockStaff) createAppAccess(app oauth.App) oauth.Access {
	token, err := oauth.NewToken()

	id, err := oauthEnv.SaveAppAccess(
		token,
		app.ClientID)

	if err != nil {
		panic(err)
	}

	return oauth.Access{
		ID:    id,
		Token: token,
	}
}

func (m mockStaff) createPersonalToken(u user.User) oauth.PersonalAccess {
	m.createAccount()

	acc := m.personalAccess()

	id, err := oauthEnv.SavePersonalToken(acc, null.StringFrom(u.UserID))
	if err != nil {
		panic(err)
	}

	acc.ID = id
	return acc
}

type mockUser struct {
	userID   string
	unionID  null.String
	email    string
	password string
	userName string
	nickname string
}

func newMockUser() mockUser {
	return mockUser{
		userID:   uuid.NewV4().String(),
		email:    fake.EmailAddress(),
		password: genPassword(),
		userName: fake.UserName(),
		nickname: fake.UserName(),
	}
}

func (m mockUser) withUnionID() mockUser {
	m.unionID = null.StringFrom(genUnionID())

	return m
}

func (m mockUser) login() user.Login {
	return user.Login{
		Email:    m.email,
		Password: m.password,
	}
}

func (m mockUser) user() user.User {
	return user.User{
		UserID:   m.userID,
		UnionID:  m.unionID,
		Email:    m.email,
		UserName: null.StringFrom(m.userName),
	}
}

func (m mockUser) order(p subs.Plan, login enum.LoginMethod) user.Order {
	app := clientApp()

	o := user.Order{
		ID:            genOrderID(),
		LoginMethod:   enum.LoginMethodEmail,
		Tier:          p.Tier,
		Cycle:         p.Cycle,
		ListPrice:     p.ListPrice,
		NetPrice:      p.NetPrice,
		PaymentMethod: enum.PayMethodWx,
		CreatedAt:     chrono.TimeNow(),
		ConfirmedAt:   chrono.TimeNow(),
		ClientType:    app.ClientType,
		ClientVersion: null.StringFrom(app.Version),
		UserIP:        null.StringFrom(app.UserIP),
		UserAgent:     null.StringFrom(app.UserAgent),
	}

	if login == enum.LoginMethodWx {
		o.UserID = m.unionID.String
	} else {
		o.UserID = m.userID
	}

	o.StartDate = chrono.DateNow()
	endTime, _ := p.Cycle.TimeAfterACycle(time.Now())
	o.EndDate = chrono.DateFrom(endTime)

	return o
}

func (m mockUser) createUser() user.User {
	u := m.user()
	app := clientApp()

	query := `
	INSERT INTO cmstmp01.userinfo
	SET user_id = ?,
	    wx_union_id = ?,
		email = ?,
		password = MD5(?),
		user_name = ?,
		client_type = ?,
		client_version = ?,
		user_ip = INET6_ATON(?),
		user_agent = ?,
		created_utc = UTC_TIMESTAMP()
	ON DUPLICATE KEY UPDATE
		user_id = ?,
		email = ?,
		password = MD5(?),
		user_name = ?`

	_, err := db.Exec(query,
		u.UserID,
		u.UnionID,
		u.Email,
		m.password,
		u.UserName,
		app.ClientType,
		app.Version,
		app.UserIP,
		app.UserAgent,
		u.UserID,
		u.Email,
		m.password,
		u.UserName)

	if err != nil {
		panic(err)
	}
	return u
}

func (m mockUser) createOrder(o user.Order) {
	query := `
	INSERT INTO premium.ftc_trade
	SET trade_no = ?,
	    user_id = ?,
		login_method = ?,
		tier_to_buy = ?,
		billing_cycle = ?,
		trade_price = ?,
		trade_amount = ?,
		payment_method = ?,
		created_utc = ?,
		confirmed_utc = ?,
		start_date = ?,
		end_date = ?,
		client_type = ?,
		client_version = ?,
		user_ip_bin = INET6_ATON(?),
		user_agent = ?`

	_, err := db.Exec(query,
		o.ID,
		o.UserID,
		o.LoginMethod,
		o.Tier,
		o.Cycle,
		o.ListPrice,
		o.NetPrice,
		o.PaymentMethod,
		o.CreatedAt,
		o.ConfirmedAt,
		o.StartDate,
		o.EndDate,
		o.ClientType,
		o.ClientVersion,
		o.UserIP,
		o.UserAgent)

	if err != nil {
		logger.WithField("trace", "createOrder").Error(err)
		panic(err)
	}
}

func (m mockUser) createWxUser() {
	query := `
	INSERT INTO user_db.wechat_userinfo
	SET union_id = ?,
		nickname = ?`

	_, err := db.Exec(query,
		m.unionID.String,
		m.nickname,
	)

	if err != nil {
		panic(err)
	}
}

func (m mockUser) createMember(order user.Order) {
	query := `
	INSERT INTO premium.ftc_vip
	SET vip_id = ?,
		member_tier = ?,
		billing_cycle = ?,
		expire_date = ?`

	_, err := db.Exec(query,
		m.userID,
		order.Tier,
		order.Cycle,
		order.EndDate,
	)

	if err != nil {
		panic(err)
	}
}

func TestMockStaff(t *testing.T) {
	m := newMockStaff()

	t.Logf("%+v", m.account())
}
