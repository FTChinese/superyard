//go:build !production

package faker

import (
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MustConfigViper() {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func ReadConfigFile() ([]byte, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filepath.Join(home, "config", "api.toml"))
}

func MustReadConfigFile() []byte {
	b, err := ReadConfigFile()
	if err != nil {
		panic(err)
	}

	return b
}

func MustSetupViper() {
	config.MustSetupViper(MustReadConfigFile())
}
