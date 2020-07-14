package staff

import "testing"

func TestNewPassportBearer(t *testing.T) {
	key := []byte("signing_key")

	pp, err := NewPassportBearer(mockAccount, key)
	if err != nil {
		t.Error(err)
	}

	t.Logf("PassportBearer: %+v", pp)

	claims, err := ParsePassportClaims(pp.Token, key)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Parsed PassportClaimns: %+v", claims)
}
