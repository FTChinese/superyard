package test

import (
	"testing"
)

func TestRepo_CreateReader(t *testing.T) {

	p := NewPersona()

	t.Logf("Creating user %s", p.FtcID)

	err := NewRepo().CreateReader(p.FtcAccount())
	if err != nil {
		t.Error(err)
	}
}

func TestRepo_CreateVIP(t *testing.T) {
	repo := NewRepo()

	err := repo.CreateVIP(NewPersona().SetVIP().FtcAccount())
	if err != nil {
		t.Error(err)
	}
}
