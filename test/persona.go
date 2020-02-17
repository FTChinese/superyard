package test

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/brianvoe/gofakeit/v4"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/reader"
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
	gofakeit.Seed(time.Now().UnixNano())

	return &Persona{
		FtcID:       uuid.New().String(),
		UnionID:     GenWxID(),
		Email:       gofakeit.Email(),
		Password:    SimplePassword(),
		UserName:    gofakeit.Username(),
		Nickname:    gofakeit.Name(),
		Avatar:      gofakeit.ImageURL(20, 20),
		OpenID:      GenWxID(),
		IP:          gofakeit.IPv4Address(),
		DeviceToken: GenDeviceToken(),
		PwToken:     GenPwResetToken(),
		VrfToken:    GenVrfToken(),
		accountKind: 0,
		linked:      false,
		payMethod:   0,
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

func (p *Persona) Membership() reader.Membership {
	return reader.Membership{
		ID: null.StringFrom(reader.GenerateMemberID()),
		AccountID: reader.AccountID{
			CompoundID: p.FtcID,
			FtcID:      null.StringFrom(p.FtcID),
			UnionID:    null.StringFrom(p.UnionID),
		},
		LegacyTier:    null.Int{},
		LegacyExpire:  null.Int{},
		Tier:          enum.TierStandard,
		Cycle:         enum.CycleYear,
		ExpireDate:    chrono.DateFrom(time.Now().AddDate(1, 0, 0)),
		PaymentMethod: enum.PayMethodWx,
		StripeSubID:   null.String{},
		StripePlanID:  null.String{},
		AutoRenewal:   false,
		Status:        reader.SubStatusNull,
		AppleSubID:    null.StringFrom(GenAppleSubID()),
		VIP:           false,
	}
}

func (p *Persona) Order(confirmed bool) reader.Order {
	orderID := GenSubID()

	order := reader.Order{
		ID: orderID,
		AccountID: reader.AccountID{
			CompoundID: p.FtcID,
			FtcID:      null.StringFrom(p.FtcID),
			UnionID:    null.StringFrom(p.UnionID),
		},
		Price:            258.00,
		Amount:           258.00,
		Tier:             enum.TierStandard,
		Cycle:            enum.CycleYear,
		Currency:         null.StringFrom("cny"),
		CycleCount:       1,
		ExtraDays:        1,
		Usage:            reader.SubsKindCreate,
		PaymentMethod:    enum.PayMethodAli,
		CreatedAt:        chrono.TimeNow(),
		UpgradeID:        null.String{},
		MemberSnapshotID: null.String{},
	}

	if confirmed {
		order.ConfirmedAt = chrono.TimeNow()
		order.StartDate = chrono.DateNow()
		order.EndDate = chrono.DateFrom(time.Now().AddDate(1, 0, 1))
	}

	return order
}
