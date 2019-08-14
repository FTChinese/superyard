package customer

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env handles FTC user data.
type Env struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository.customer")

func normalizeMemberTier(vipType int64) enum.Tier {
	switch vipType {

	case 10:
		return enum.TierStandard

	case 100:
		return enum.TierPremium

	default:
		return enum.InvalidTier
	}
}
