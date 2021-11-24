package test

import (
	"testing"
)

func TestDeleteFtcMember(t *testing.T) {
	p := NewPersona()

	repo := NewRepo()

	repo.MustCreateReader(p.FtcAccount())
	repo.MustCreateMembership(p.Membership())

	t.Logf("%s", p.FtcID)
}
