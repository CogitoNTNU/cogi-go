package env

import "github.com/sirupsen/logrus"

import (
	"github.com/spf13/viper"
	"os"
)
var (
	PROD string = "PRODUCTION"
	DEV  string = "DEVELOPMENT"
)

type EnvConfig struct {
	env_path string
  logger *logrus.Logger
}

func Configure(path string) *EnvConfig {
  return &EnvConfig{
    env_path: path,
    logger:   logrus.New(),
  }
}

func (p *EnvConfig) Read(key string) (value string) {
    logger := p.logger
    
    if value := os.Getenv(key); value != "" {
        return value
    }

    viper.SetConfigFile(p.env_path)
    err := viper.ReadInConfig()

    if err != nil {
        logger.WithError(err).Fatal("Couldn't intialize and find environmental variable.")
    }

    value, ok := viper.Get(key).(string)

    if !ok {
        logger.Fatal("Invalid type assertion on value from env.")
    }

    return value
}
