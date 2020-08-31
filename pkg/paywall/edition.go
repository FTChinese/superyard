package paywall

import "github.com/FTChinese/go-rest/enum"

type Edition struct {
	Tier  enum.Tier  `json:"tier" db:"tier"`
	Cycle enum.Cycle `json:"cycle" db:"cycle"`
}

func NewStdMonthEdition() Edition {
	return Edition{
		Tier:  enum.TierStandard,
		Cycle: enum.CycleMonth,
	}
}

func NewStdYearEdition() Edition {
	return Edition{
		Tier:  enum.TierStandard,
		Cycle: enum.CycleYear,
	}
}

func NewPremiumEdition() Edition {
	return Edition{
		Tier:  enum.TierPremium,
		Cycle: enum.CycleYear,
	}
}

func (e Edition) NamedKey() string {
	return e.Tier.String() + "_" + e.Cycle.String()
}

// StringCN produces a human readable string of this edition.
// * 标准会员/年
// * 标准会员/月
// * 高端会员/年
func (e Edition) StringCN() string {
	return e.Tier.StringCN() + "/" + e.Cycle.StringCN()
}
