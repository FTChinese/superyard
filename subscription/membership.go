package subscription

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	standard    = "standard"
	premium     = "premium"
	year        = "year"
	month       = "month"
	keyStdYear  = "standard_year"
	keyStdMonth = "standard_month"
	keyPrmYear  = "premium_year"
)

// Tier represents member tier enum
type Tier string

// UnmarshalJSON parses a tier raw value
func (t *Tier) UnmarshalJSON(b []byte) error {
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	raw = strings.TrimSpace(raw)

	switch raw {
	case standard:
		*t = TierStandard
	case premium:
		*t = TierPremium
	default:
		return errors.New("raw value for enum Tier not found")
	}

	return nil
}

// MarshalJSON stringifies a Tier value.
func (t Tier) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}

// Cycle represents billing cycle enum
type Cycle string

// UnmarshalJSON parses a tier raw value
func (c *Cycle) UnmarshalJSON(b []byte) error {
	var raw string
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	raw = strings.TrimSpace(raw)

	switch raw {
	case year:
		*c = CycleYear
	case month:
		*c = CycleMonth
	default:
		return errors.New("raw value for enum Cycle not found")
	}

	return nil
}

// MarshalJSON stringifies a Tier value.
func (c Cycle) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

// Enum instance
const (
	TierStandard Tier  = standard
	TierPremium  Tier  = premium
	CycleYear    Cycle = year
	CycleMonth   Cycle = month
)
