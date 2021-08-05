package config

import "github.com/spf13/viper"

func SetupViper(prod bool) error {
	viper.SetConfigName("api")
	viper.SetConfigType("toml")
	if prod {
		viper.AddConfigPath("/data/opt/server/API/config")
	} else {
		viper.AddConfigPath("$HOME/config")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func MustSetupViper(prod bool) {
	if err := SetupViper(prod); err != nil {
		panic(err)
	}
}
