package env

import "github.com/sirupsen/logrus"

import (
	"github.com/spf13/viper"
)

var (
	PROD string = "PRODUCTION"
	DEV  string = "DEVELOPMENT"
)

type EnvConfig struct {
	env_path string
}

func Configure(path string) *EnvConfig {
  return &EnvConfig{env_path: path}
}

func (p *EnvConfig) Read(key string) (value string) {
	logger := logrus.New()
	viper.SetConfigFile(p.env_path)
	err := viper.ReadInConfig()

	if err != nil {
		logger.WithError(err).Fatal("Couldn't intialize and find environmental variable.")
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		logger.Fatal("Invalid type assertion on value from env.")
	}

	return
}
