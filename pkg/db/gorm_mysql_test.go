package db

import (
	"testing"

	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/config"
)

func TestMustNewGormDB(t *testing.T) {
	faker.MustSetupViper()

	conn := config.MustMySQLReadConn(false)

	_, err := NewGormDB(conn)
	if err != nil {
		t.Error(err)
	}
}
