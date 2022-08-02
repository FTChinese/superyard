package repository

import (
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/FTChinese/superyard/pkg/db"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadFile(t *testing.T) {
	faker.MustSetupViper()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
		return
	}

	objectName := "ftchinese-v6.7.1-stripe-debug.apk"

	fp := filepath.Join(home, "GolandProjects", objectName)

	c := db.MustMinIOClient(config.MustGetMinIOConfig())

	info, err := UploadFile(c, objectName, fp)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", info)
}
