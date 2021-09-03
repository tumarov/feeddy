package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Telegram           Telegram
	Sentry             Sentry
	DBPath             string
	DBName             string `mapstructure:"db_name"`
	Messages           Messages
	Environment        string
	ApplicationRelease string
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := initMessages(&cfg); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	if err := initConfigs(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initConfigs(cfg *Config) error {
	if err := cfg.Telegram.init(cfg); err != nil {
		return err
	}

	return nil
}

func parseEnv(cfg *Config) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	viper.AutomaticEnv()

	cfg.Environment = viper.GetString("ENVIRONMENT")
	cfg.ApplicationRelease = viper.GetString("APPLICATION_RELEASE")
	cfg.DBPath = viper.GetString("DB_PATH")

	cfg.Telegram.parseEnv(cfg)
	cfg.Sentry.parseEnv(cfg)

	return nil
}
