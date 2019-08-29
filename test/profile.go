package test

import (
	"github.com/google/uuid"
	"github.com/icrowley/fake"
)

const (
	MyFtcID    = "e1a1f5c0-0e23-11e8-aa75-977ba2bcc6ae"
	MyFtcEmail = "neefrankie@163.com"
	MyUnionID  = "ogfvwjk6bFqv2yQpOrac0J3PqA0o"
	MyEmail    = "neefrankie@gmail.com"
)

type Profile struct {
	FtcID    string
	UnionID  string
	StripeID string
	Email    string
	Password string
	UserName string
	Nickname string
	Avatar   string
	OpenID   string
	IP       string
}

func NewProfile() Profile {
	return Profile{
		FtcID:    uuid.New().String(),
		UnionID:  GenWxID(),
		StripeID: GetCusID(),
		Email:    fake.EmailAddress(),
		Password: fake.SimplePassword(),
		UserName: fake.UserName(),
		Nickname: fake.UserName(),
		Avatar:   GenAvatar(),
		OpenID:   GenWxID(),
		IP:       fake.IPv4(),
	}
}

var MyProfile = Profile{
	FtcID:    MyFtcID,
	UnionID:  MyUnionID,
	StripeID: "cus_FOgRRgj9aMzpAv",
	Email:    MyEmail,
	Password: "12345678",
	UserName: "weiguo.ni",
	Nickname: fake.UserName(),
	Avatar:   "http://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIibCfVIicoNXZ15Af6nWkXwq5QgFcrNdkEKMHT7P1oJVI6McLT2qFia2ialF4FSMnm33yS0eAq7MK1cA/132",
	IP:       fake.IPv4(),
}
