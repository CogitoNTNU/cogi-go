package env

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	PROD string = "PRODUCTION"
	DEV  string = "DEVELOPMENT"
)

type EnvConfig struct {
	env_path string
	logger   *logrus.Logger
}

func Configure(path string) *EnvConfig {
	return &EnvConfig{
		env_path: path,
		logger:   logrus.New(),
	}
}

func (p *EnvConfig) Read(key string) (value string, err error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	viper.SetConfigFile(p.env_path)
	err = viper.ReadInConfig()
	if err != nil {
		return "", fmt.Errorf("couldn't read config file: %w", err)
	}

	if !viper.IsSet(key) {
		return "", fmt.Errorf("key %q not found in config", key)
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		return "", fmt.Errorf("invalid type assertion on value for key %q from env", key)
	}

	return value, nil
}
