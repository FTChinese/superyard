package faker

import "github.com/spf13/viper"

func MustConfigViper() {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
