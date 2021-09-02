package config

import "github.com/spf13/viper"

type Sentry struct {
	DSN string
}

func (s Sentry) parseEnv(cfg *Config) {
	cfg.Sentry.DSN = viper.GetString("SENTRY_DSN")
}
