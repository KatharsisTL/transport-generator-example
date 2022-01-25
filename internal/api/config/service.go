package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type ServiceConfig struct {
	LogLevel   string `default:"debug" envconfig:"LOG_LEVEL"`
	Bind       string `default:":9000" envconfig:"BIND_ADDR"`
	HealthBind string `default:":9091" envconfig:"BIND_HEALTH"`
}

var service *ServiceConfig

func Server() ServiceConfig {
	if service != nil {
		return *service
	}

	service = &ServiceConfig{}
	if err := envconfig.Process("", service); err != nil {
		log.Fatal().Err(err).Msg("failed to parse service config")
	}

	return *service
}
