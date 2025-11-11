package env

import (
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

func (p *EnvConfig) ReadIntDefault(key string, def int) int {
	if p.logger == nil {
		p.logger = logrus.New()
	}

	viper.SetConfigFile(p.env_path)
	if err := viper.ReadInConfig(); err != nil {
		p.logger.WithError(err).Warn("Couldn't read env, using default for " + key)
		return def
	}

	val := viper.Get(key)
	if val == nil {
		return def
	}

	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		return def
	default:
		return def
	}
}

func Configure(path string) *EnvConfig {
	return &EnvConfig{
		env_path: path,
		logger:   logrus.New(),
	}
}

func Read[T any](p *EnvConfig, key string) T {
	logger := p.logger
	if logger == nil { 
		logger = logrus.New()
		p.logger = logger
	}

	viper.SetConfigFile(p.env_path)
	err := viper.ReadInConfig()

	if err != nil {
		logger.WithError(err).Fatal("Couldn't initialize and find environmental variable.")
	}

	value, ok := viper.Get(key).(T)

	if !ok {
		logger.Fatal("Invalid type assertion on value from env.")
	}

	return value
}

func (p *EnvConfig) Read(key string) string {
	return Read[string](p, key)
}