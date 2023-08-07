package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

// internal, had to be public to properly unmarshal
type CommonStanConfig struct {
	StanURL     string `env:"STREAMING_URL" env-description:"NATS-Streaming URL"        env-required:"true"`
	StanSubject string `env:"SUBJECT"       env-description:"publishing subject"        env-required:"true"`
	StanCluster string `env:"CLUSTER_ID"    env-description:"NATS-Streaming cluster id" env-required:"true"`
}

// config of the publisher service
type FillerConfig struct {
	CommonStanConfig        // StanURL, StanSubject, StanCluster
	StanClient       string `env:"PUBLISHER_ID" env-description:"NATS-Streaming client id"  env-required:"true"`
	DataDir          string `env:"DATA_DIR"     env-description:"path to messages"          env-required:"true"`
}

// config of the main service
type ServiceConfig struct {
	CommonStanConfig        // StanURL, StanSubject, StanCluster
	StanClient       string `env:"SUBSCRIBER_ID" env-description:"NATS-Streaming client id"         env-required:"true"`
	StanDurableName  string `env:"DURABLE_NAME"  env-description:"name of the durable subscription" env-required:"true"`
	DatabaseURL      string `env:"DATABASE_URL"  env-description:"postgres connection URL"          env-required:"true"`
}

// read config from environment variables
//
// returns invalid config on error
func Read[T ServiceConfig | FillerConfig]() (T, error) {
	var (
		cfg T
		err error
	)
	if err = cleanenv.ReadEnv(&cfg); err != nil {
		// TODO: open an issue at cleanenv: cleanenv.GetDescription can't handle generic types :\
		/*
			helpHeader := "Expected config:"
			if help, descErr := cleanenv.GetDescription(cfg, &helpHeader); descErr == nil {
				log.Println(help)
			}
		*/
		err = fmt.Errorf("could not read config: %w", err)
	}
	return cfg, err
}
