package test

import "testing"

func TestStaff_Account(t *testing.T) {
	t.Logf("Account: %+v", NewStaff().Account())
}

func TestStaff_Credentials(t *testing.T) {
	t.Logf("Credentials: %+v", NewStaff().Credentials())
}

func TestStaff_SignUp(t *testing.T) {
	t.Logf("SignUp: %+v", NewStaff().SignUp())
}

func TestStaff_PwResetSession(t *testing.T) {
	t.Logf("PwResetSession: %+v", NewStaff().PwResetSession())
}
