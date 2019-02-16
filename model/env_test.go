package model

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
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

type OAuthAccess struct {
	SessionID string
	// Example: 16_Ix0E3WfWs9u5Rh9f-lB7_LgsQJ4zm1eodolFJpSzoQibTAuhIlp682vDmkZSaYIjD9gekOa1zQl-6c6S_CrN_cN9vx9mybwXNVgFbwPMMwM
	AccessToken string `json:"access_token"`
	// Example: 7200
	ExpiresIn int64 `json:"expires_in"`
	// Exmaple: 16_IlmA9eLGjJw7gBKBT48wff1V1hAYAdpmIqUAypspepm6DsQ6kkcLeZmP932s9PcKp1WM5P_1YwUNQqF-29B_0CqGTqMpWkaaiNSYp26MmB4
	RefreshToken string `json:"refresh_token"`
	// Example: ob7fA0h69OO0sTLyQQpYc55iF_P0
	OpenID string `json:"openid"`
	// Example: snsapi_userinfo
	Scope string `json:"scope"`
	// Example: String:ogfvwjk6bFqv2yQpOrac0J3PqA0o Valid:true
	UnionID null.String `json:"unionid"`
}

func (a *OAuthAccess) GenerateSessionID() {
	data := fmt.Sprintf("%s:%s:%s", a.AccessToken, a.RefreshToken, a.OpenID)
	h := md5.Sum([]byte(data))
	a.SessionID = hex.EncodeToString(h[:])
}

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

func generateAvatarURL() string {
	return fmt.Sprintf("http://thirdwx.qlogo.cn/mmopen/vi_32/%s/132", fake.CharactersN(90))
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

func generateToken() string {
	token, _ := gorest.RandomBase64(82)
	return token
}

func generateWxID() string {
	id, _ := gorest.RandomBase64(21)
	return id
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

func (m mockUser) withEmail(email string) mockUser {
	m.email = email
	return m
}

func (m mockUser) withPassword(pw string) mockUser {
	m.password = pw
	return m
}

func (m mockUser) login() user.Login {
	return user.Login{
		Email:    m.email,
		Password: m.password,
	}
}

func (m mockUser) loginHistory() user.LoginHistory {
	app := clientApp()
	h := user.LoginHistory{
		UserID:     m.userID,
		AuthMethod: enum.LoginMethodEmail,
	}

	h.ClientType = app.ClientType
	h.Version = null.StringFrom(app.Version)
	h.UserIP = null.StringFrom(app.UserIP)
	h.UserAgent = null.StringFrom(app.UserAgent)

	return h
}

func (m mockUser) user() user.User {
	return user.User{
		UserID:   m.userID,
		UnionID:  m.unionID,
		Email:    m.email,
		UserName: null.StringFrom(m.userName),
	}
}

func (m mockUser) wxUser() user.WxInfo {
	return user.WxInfo{
		UnionID:    m.unionID.String,
		Nickname:   fake.UserName(),
		AvatarURL:  generateAvatarURL(),
		Gender:     enum.GenderFemale,
		Country:    fake.Country(),
		Province:   fake.State(),
		City:       fake.City(),
		Privileges: []string{},
		CreatedAt:  chrono.TimeNow(),
		UpdatedAt:  chrono.TimeNow(),
	}
}
func (m mockUser) wxAccess() OAuthAccess {
	acc := OAuthAccess{
		AccessToken:  generateToken(),
		ExpiresIn:    7200,
		RefreshToken: generateToken(),
		OpenID:       generateWxID(),
		Scope:        "snsapi_userinfo",
		UnionID:      m.unionID,
	}

	acc.GenerateSessionID()

	return acc
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
	  	wx_union_id = NULL,
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

func (m mockUser) createLoginHistory() user.LoginHistory {
	h := m.loginHistory()

	query := `
	INSERT INTO user_db.login_history
	SET user_id = ?,
		auth_method = ?,
		client_type = ?,
		client_version = ?,
		user_ip = INET6_ATON(?),
		user_agent = ?`

	_, err := db.Exec(query,
		h.UserID,
		h.AuthMethod,
		h.ClientType,
		h.Version,
		h.UserIP,
		h.UserAgent)

	if err != nil {
		panic(err)
	}

	return h
}
func (m mockUser) createOrder(loginMethod enum.LoginMethod) user.Order {

	o := m.order(stdPlan, loginMethod)

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
		panic(err)
	}

	return o
}

func (m mockUser) createWxUser() user.WxInfo {

	u := m.wxUser()

	query := `
	INSERT INTO user_db.wechat_userinfo
	SET union_id = ?,
		nickname = ?,
		avatar_url = ?,
		gender = ?,
		country = ?,
		province = ?,
		city = ?,
		privilege = NULLIF(?, ''),
		created_utc = ?,
	    updated_utc = ?`

	_, err := db.Exec(query,
		u.UnionID,
		u.Nickname,
		u.AvatarURL,
		1,
		u.Country,
		u.Province,
		u.City,
		"",
		u.CreatedAt,
		u.UpdatedAt,
	)

	if err != nil {
		panic(err)
	}

	return u
}

func (m mockUser) createWxAccess() OAuthAccess {

	acc := m.wxAccess()

	c := clientApp()

	query := `
	INSERT INTO user_db.wechat_access
	SET session_id = UNHEX(?),
		app_id = ?,
		access_token = ?,
		expires_in = ?,
		refresh_token = ?,
		open_id = ?,
		scope = ?,
		union_id = ?,
		client_type = ?,
		client_version = ?,
		user_ip = INET6_ATON(?),
		user_agent = ?`

	appID := "wxacddf1c20516eb69"

	_, err := db.Exec(query,
		acc.SessionID,
		appID,
		acc.AccessToken,
		acc.ExpiresIn,
		acc.RefreshToken,
		acc.OpenID,
		acc.Scope,
		acc.UnionID,
		c.ClientType,
		c.Version,
		c.UserIP,
		c.UserAgent,
	)

	if err != nil {
		panic(err)
	}

	return acc
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
