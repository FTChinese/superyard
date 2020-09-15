package test

import (
	"github.com/FTChinese/superyard/faker"
	"testing"
)

func TestCreateFtcMember(t *testing.T) {
	p := NewPersona()

	repo := NewRepo()

	repo.MustCreateReader(p.FtcAccount())

	t.Logf("%s", faker.MustMarshalIndent(p.FtcSubsCreationInput()))
}

func TestUpdateFtcMember(t *testing.T) {
	p := NewPersona()

	repo := NewRepo()

	repo.MustCreateReader(p.FtcAccount())
	repo.MustCreateMembership(p.Membership())

	t.Logf("%s", p.FtcID)
	t.Logf("%s", faker.MustMarshalIndent(p.FtcSubsUpdateInput()))
}

func TestDeleteFtcMember(t *testing.T) {
	p := NewPersona()

	repo := NewRepo()

	repo.MustCreateReader(p.FtcAccount())
	repo.MustCreateMembership(p.Membership())

	t.Logf("%s", p.FtcID)
}
