package config

import (
	"log"

	"github.com/caarlos0/env"
	configUtils "github.com/fubotv/fubotv-config-utils/v8/cfg"
	"github.com/fubotv/fubotv-metrics/v12/datadog/config"
)

type Cfg struct {
	AppCfg      AppCfg
	DataDogCfg  config.DataDogConfig
	DatabaseCfg DatabaseCfg
}

type AppCfg struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"keyplay-metadata-api"`
	ServerPort  int    `env:"SERVER_PORT" envDefault:"8080"`
	Env         string `env:"DEPLOY_ENV" envDefault:"local"`
	Version     configUtils.VersionConfig
	SwaggerDir  string `env:"SWAGGER_DIR" envDefault:"swaggerui"`
}

type DatabaseCfg struct {
	BTConnectionCredentials string `env:"BT_CONNECTION_CREDENTIALS" `
	BTInstanceId            string `env:"BT_INSTANCE_ID" envDefault:"fubotv-dev"`
	BTProjectId             string `env:"BT_PROJECT_ID" envDefault:"fubotv-dev"`
}

var cfg *Cfg

func init() {
	cfg = New()
}

func New() *Cfg {
	a := &AppCfg{}
	err := env.Parse(a)
	if err != nil {
		log.Fatal(err)
	}

	v := &configUtils.VersionConfig{}
	err = env.Parse(v)
	if err != nil {
		log.Fatal(err)
	}
	a.Version = *v

	d := &config.DataDogConfig{}
	err = env.Parse(d)
	if err != nil {
		log.Fatal(err)
	}

	db := &DatabaseCfg{}
	err = env.Parse(db)
	if err != nil {
		log.Fatal(err)
	}

	return &Cfg{AppCfg: *a, DataDogCfg: *d, DatabaseCfg: *db}
}

func GetConfig() *Cfg {
	return cfg
}
