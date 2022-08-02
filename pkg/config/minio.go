package config

import (
	"github.com/spf13/viper"
	"log"
)

type MinIOConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

func GetMinIOConfig() (MinIOConfig, error) {
	var c MinIOConfig
	err := viper.UnmarshalKey("minio", &c)
	if err != nil {
		return MinIOConfig{}, err
	}

	return c, nil
}

func MustGetMinIOConfig() MinIOConfig {
	c, err := GetMinIOConfig()
	if err != nil {
		log.Fatalf("cannot get minio config: %s", err.Error())
	}

	return c
}
