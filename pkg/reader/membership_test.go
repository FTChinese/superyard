package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeUnix(t *testing.T) {
	expireDate := time.Unix(1710864000, 0)

	t.Logf("%s", chrono.DateFrom(expireDate))

	t.Logf("%d", chrono.TimeFrom(expireDate).UTC().Unix())
}

func TestMembership_Normalize(t *testing.T) {
	m1 := Membership{
		CompoundID:   null.StringFrom(uuid.New().String()),
		LegacyTier:   null.IntFrom(10),
		LegacyExpire: null.IntFrom(1710864000),
		Edition:      paywall.Edition{},
		ExpireDate:   chrono.Date{},
	}

	t.Logf("%t", m1.IsZero())

	assert.False(t, m1.IsZero())

	m1 = m1.Normalize()

	assert.Equal(t, time.Unix(m1.LegacyExpire.Int64, 0), m1.ExpireDate.Time)
	assert.Equal(t, m1.Tier, enum.TierStandard)

	m2 := Membership{
		CompoundID:   null.StringFrom(uuid.New().String()),
		LegacyTier:   null.Int{},
		LegacyExpire: null.Int{},
		Edition: paywall.Edition{
			Tier:  enum.TierStandard,
			Cycle: enum.CycleYear,
		},
		ExpireDate: chrono.DateFrom(time.Now().AddDate(1, 0, 0)),
	}

	m2 = m2.Normalize()

	assert.Equal(t, m2.LegacyExpire.Int64, m2.ExpireDate.Unix())
	assert.Equal(t, m2.LegacyTier, null.IntFrom(10))
}
