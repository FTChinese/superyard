package admin

import (
	"testing"

	"github.com/icrowley/fake"
	"gitlab.com/ftchinese/backyard-api/staff"
)

func TestCreateStaff(t *testing.T) {
	pass, err := devEnv.createStaff(mockStaff)

	if err != nil {
		t.Error(err)
	}

	t.Log(pass)
}

func TestNewStaff(t *testing.T) {
	parcel, err := devEnv.NewStaff(staff.Account{
		Email:        fake.EmailAddress(),
		UserName:     fake.UserName(),
		DisplayName:  fake.FullName(),
		Department:   "tech",
		GroupMembers: 3,
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(parcel)
}

func TestStaffRoster(t *testing.T) {
	accounts, err := devEnv.StaffRoster(1, 20)

	if err != nil {
		t.Error(err)
	}

	t.Log(accounts)
}

func TestUpdateStaff(t *testing.T) {
	err := devEnv.UpdateStaff(mockStaff.UserName, mockStaff)

	if err != nil {
		t.Error(err)
	}
}

func TestDeactivateStaff(t *testing.T) {
	err := devEnv.deactivateStaff(mockStaff.UserName)

	if err != nil {
		t.Error(err)
	}
}
func TestActivateStaff(t *testing.T) {
	err := devEnv.ActivateStaff(mockStaff.UserName)

	if err != nil {
		t.Error(err)
	}
}

func TestRevokeStaffVIP(t *testing.T) {
	err := devEnv.revokeStaffVIP(mockStaff.UserName)

	if err != nil {
		t.Error(err)
	}
}
