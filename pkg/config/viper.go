package config

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SetupViper(b []byte) error {
	viper.SetConfigType("toml")

	err := viper.ReadConfig(bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	return nil
}

func MustSetupViper(b []byte) {
	if err := SetupViper(b); err != nil {
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
