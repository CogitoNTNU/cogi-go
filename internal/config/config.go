package config

import (
	"github.com/CogitoNTNU/cogi-go/internal/util/env"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type Api struct {
	AppName    string `split_words:"true" default:"cogi-go api v1"`
	LogLevel   string `split_words:"true" default:"DEBUG"`
	MainHost   string `split_words:"true" default:"localhost"`
	MainPort   int    `split_words:"true" default:"8080"`
	HealthPort int    `split_words:"true" default:"8083"`
}

func LoadApiConfig() *Api {
	var cfg Api = Api{
		AppName: "cogi-go api v1",
	}
	return &cfg
}

type LocalDbConfig struct {
	SQLHost     string `split_words:"true" default:"localhost"`
	SQLPort     int    `split_words:"true" default:"5432"`
	SQLUser     string `split_words:"true" default:"root"`
	SQLPassword string `split_words:"true" default:"postgres"`
	SQLDatabase string `split_words:"true" default:"postgres"`
}

func NewLocalDbConfig(options ...LocalDbOption) *LocalDbConfig {
	var localDbCfg LocalDbConfig

	localDbCfg.loadDefault()

	for _, option := range options {
		option(&localDbCfg)
	}

	return &localDbCfg
}

func (c *LocalDbConfig) loadDefault() (*LocalDbConfig, error) {
	if err := envconfig.Process("", &c); err != nil {
		return c, err
	}
	return c, nil
}

func (a *Api) CorsNew(e *env.EnvConfig) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://cogito-ntnu.no"}
	return cors.New(config)
}

