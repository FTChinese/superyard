package readers

import (
	"github.com/FTChinese/superyard/test"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestEnv_AccountByFtcID(t *testing.T) {
	env := NewEnv(test.DBX, zaptest.NewLogger(t))

	a, err := env.AccountByFtcID("8680d6be-9540-4915-ac0d-23acfe636469")

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", a)
}
