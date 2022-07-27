package test

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/ids"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

const (
	MyFtcID   = "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae"
	MyUnionID = "ogfvwjk6bFqv2yQpOrac0J3PqA0o"
	MyEmail   = "neefrankie@gmail.com"
)

type WxInfo struct {
	UnionID  string      `db:"union_id"`
	Nickname null.String `db:"nickname"`
	Avatar   null.String `db:"avatar"`
	Gender   enum.Gender `db:"gender"`
	Country  null.String `db:"country"`
	Province null.String `db:"province"`
	City     null.String `db:"city"`
}

type Persona struct {
	FtcID       string `db:"ftc_id"`
	UnionID     string `db:"wx_union_id"`
	StripeID    string `db:"stripe_customer_id"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	UserName    string `db:"user_name"`
	Nickname    string
	Avatar      string
	OpenID      string
	IP          string
	DeviceToken string
	PwToken     string
	VrfToken    string

	accountKind enum.AccountKind
	payMethod   enum.PayMethod
	expired     bool
	vip         bool

	orders map[string]subs.Order
	member reader.Membership
}

var MyProfile = Persona{
	FtcID:    MyFtcID,
	UnionID:  MyUnionID,
	StripeID: "cus_FOgRRgj9aMzpAv",
	Email:    MyEmail,
	Password: "12345678",
	UserName: "weiguo.ni",
	Nickname: "🐆测试",
	Avatar:   "http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIibCfVIicoNXZ15Af6nWkXwq5QgFcrNdkEKMHT7P1oJVI6McLT2qFia2ialF4FSMnm33yS0eAq7MK1cA/132",
	IP:       gofakeit.IPv4Address(),
}

func NewPersona() *Persona {
	faker.SeedGoFake()

	return &Persona{
		FtcID:       uuid.New().String(),
		UnionID:     faker.GenWxID(),
		Email:       gofakeit.Email(),
		Password:    faker.SimplePassword(),
		UserName:    gofakeit.Username(),
		Nickname:    gofakeit.Name(),
		Avatar:      gofakeit.ImageURL(20, 20),
		OpenID:      faker.GenWxID(),
		IP:          gofakeit.IPv4Address(),
		DeviceToken: faker.GenToken32Bytes(),
		PwToken:     faker.GenToken32Bytes(),
		VrfToken:    faker.GenToken32Bytes(),

		accountKind: enum.AccountKindFtc,
		payMethod:   enum.PayMethodAli,
		expired:     false,

		orders: make(map[string]subs.Order),
		member: reader.Membership{},
	}
}

func (p *Persona) SetAccountKind(k enum.AccountKind) *Persona {
	p.accountKind = k
	return p
}

func (p *Persona) SetPayMethod(m enum.PayMethod) *Persona {
	p.payMethod = m
	return p
}

func (p *Persona) SetExpired(expired bool) *Persona {
	p.expired = expired
	return p
}

func (p *Persona) SetVIP() *Persona {
	faker.SeedGoFake()
	p.vip = true

	return p
}

func (p *Persona) ReaderIDs() ids.UserIDs {

	var uid ids.UserIDs
	switch p.accountKind {
	case enum.AccountKindFtc:
		uid = ids.UserIDs{
			FtcID:   null.StringFrom(p.FtcID),
			UnionID: null.String{},
		}

	case enum.AccountKindWx:
		uid = ids.UserIDs{
			FtcID:   null.String{},
			UnionID: null.StringFrom(p.UnionID),
		}

	case enum.AccountKindLinked:
		uid = ids.UserIDs{
			FtcID:   null.StringFrom(p.FtcID),
			UnionID: null.StringFrom(p.UnionID),
		}
	}

	return uid
}

func (p *Persona) FtcAccount() reader.FtcAccount {
	return reader.FtcAccount{
		UserIDs:    p.ReaderIDs(),
		StripeID:   null.StringFrom(p.StripeID),
		Email:      null.StringFrom(p.Email),
		UserName:   null.StringFrom(p.UserName),
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
		Password:   p.Password,
		CreatedBy:  "weiguo.ni",
		VIP:        p.vip,
	}
}

func (p *Persona) WxInfo() WxInfo {
	return WxInfo{
		UnionID:  p.UnionID,
		Nickname: null.StringFrom(p.Nickname),
		Avatar:   null.StringFrom(p.Avatar),
		Gender:   enum.Gender(Rand.Intn(3)),
		Country:  null.StringFrom(gofakeit.Country()),
		Province: null.StringFrom(gofakeit.State()),
		City:     null.StringFrom(gofakeit.City()),
	}
}

func (p *Persona) Membership() reader.Membership {
	m := reader.Membership{
		Edition: paywall.Edition{
			Tier:  enum.TierStandard,
			Cycle: enum.CycleYear,
		},
		ExpireDate: chrono.DateFrom(time.Now().AddDate(1, 0, 1)),
		PayMethod:  p.payMethod,
	}

	switch p.accountKind {
	case enum.AccountKindFtc:
		m.CompoundID = null.StringFrom(p.FtcID)
		m.FtcID = null.StringFrom(p.FtcID)
		m.UnionID = null.String{}

	case enum.AccountKindWx:
		m.CompoundID = null.StringFrom(p.UnionID)
		m.FtcID = null.String{}
		m.UnionID = null.StringFrom(p.UnionID)

	case enum.AccountKindLinked:
		m.CompoundID = null.StringFrom(p.FtcID)
		m.FtcID = null.StringFrom(p.FtcID)
		m.UnionID = null.StringFrom(p.UnionID)
	}

	if p.expired {
		m.ExpireDate = chrono.DateFrom(time.Now().AddDate(0, -6, 0))
	}

	switch p.payMethod {
	case enum.PayMethodStripe:
		m.StripeSubsID = null.StringFrom(faker.GenStripeSubID())
		m.StripePlanID = null.StringFrom(faker.GenStripePlanID())
		m.AutoRenewal = true
		m.Status = enum.SubsStatusActive

	case enum.PayMethodApple:
		m.AppleSubsID = null.StringFrom(faker.GenAppleSubID())

	case enum.PayMethodB2B:
		m.B2BLicenceID = null.StringFrom(faker.GenLicenceID())
	}

	return m.Normalize()
}
