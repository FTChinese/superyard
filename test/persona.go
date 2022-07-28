package test

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/faker"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/guregu/null"
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
