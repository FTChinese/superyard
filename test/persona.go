package test

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"time"
)

const (
	MyFtcID    = "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae"
	MyFtcEmail = "neefrankie@163.com"
	MyUnionID  = "ogfvwjk6bFqv2yQpOrac0J3PqA0o"
	MyEmail    = "neefrankie@gmail.com"
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

	accountKind reader.AccountKind
	linked      bool
	payMethod   enum.PayMethod
	expired     bool
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
	SeedGoFake()

	return &Persona{
		FtcID:       uuid.New().String(),
		UnionID:     genWxID(),
		Email:       gofakeit.Email(),
		Password:    simplePassword(),
		UserName:    gofakeit.Username(),
		Nickname:    gofakeit.Name(),
		Avatar:      gofakeit.ImageURL(20, 20),
		OpenID:      genWxID(),
		IP:          gofakeit.IPv4Address(),
		DeviceToken: mustGenToken(),
		PwToken:     mustGenToken(),
		VrfToken:    mustGenToken(),
		accountKind: reader.AccountKindFtc,
		linked:      false,
		payMethod:   enum.PayMethodAli,
		expired:     false,
	}
}

func (p *Persona) SetAccountKind(k reader.AccountKind) *Persona {
	p.accountKind = k
	return p
}

func (p *Persona) SetLinked(linked bool) *Persona {
	p.linked = linked
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

func (p *Persona) Membership() subs.Membership {
	m := subs.Membership{
		Tier:       enum.TierStandard,
		Cycle:      enum.CycleYear,
		ExpireDate: chrono.DateFrom(time.Now().AddDate(1, 0, 1)),
		PayMethod:  p.payMethod,
	}

	switch p.accountKind {
	case reader.AccountKindFtc:
		m.CompoundID = p.FtcID
		m.FtcID = null.StringFrom(p.FtcID)
		m.UnionID = null.String{}

	case reader.AccountKindWx:
		m.CompoundID = p.UnionID
		m.FtcID = null.String{}
		m.UnionID = null.StringFrom(p.UnionID)

	case reader.AccountKindLinked:
		m.CompoundID = p.FtcID
		m.FtcID = null.StringFrom(p.FtcID)
		m.UnionID = null.StringFrom(p.UnionID)
	}

	if p.expired {
		m.ExpireDate = chrono.DateFrom(time.Now().AddDate(0, -6, 0))
	}

	switch p.payMethod {
	case enum.PayMethodStripe:
		m.StripeSubsID = null.StringFrom(genStripeSubID())
		m.StripePlanID = null.StringFrom(genStripePlanID())
		m.AutoRenewal = true
		m.Status = enum.SubStatusActive

	case enum.PayMethodApple:
		m.AppleSubsID = null.StringFrom(genAppleSubID())

	case enum.PayMethodB2B:
		m.B2BLicenceID = null.StringFrom(genLicenceID())
	}

	return m.Normalize()
}

func (p *Persona) Order(confirmed bool) subs.Order {

	order := subs.Order{
		ID:               genOrderID(),
		Price:            258.00,
		Amount:           258.00,
		Tier:             enum.TierStandard,
		Cycle:            enum.CycleYear,
		Currency:         null.StringFrom("cny"),
		CycleCount:       1,
		ExtraDays:        1,
		Kind:             subs.KindCreate,
		PaymentMethod:    p.payMethod,
		CreatedAt:        chrono.TimeNow(),
		UpgradeID:        null.String{},
		MemberSnapshotID: null.String{},
	}

	switch p.accountKind {
	case reader.AccountKindFtc:
		order.CompoundID = p.FtcID
		order.FtcID = null.StringFrom(p.FtcID)
		order.UnionID = null.String{}

	case reader.AccountKindWx:
		order.CompoundID = p.UnionID
		order.FtcID = null.String{}
		order.UnionID = null.StringFrom(p.UnionID)

	case reader.AccountKindLinked:
		order.CompoundID = p.FtcID
		order.FtcID = null.StringFrom(p.FtcID)
		order.UnionID = null.StringFrom(p.UnionID)
	}

	if confirmed {
		order.ConfirmedAt = chrono.TimeNow()
		order.StartDate = chrono.DateNow()
		order.EndDate = chrono.DateFrom(time.Now().AddDate(1, 0, 1))
	}

	return order
}
