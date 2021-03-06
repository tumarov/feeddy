package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Telegram      Telegram
	Sentry        Sentry
	DBFile        string
	DBPath        string
	DBName        string `mapstructure:"db_name"`
	Messages      Messages
	Environment   string
	ReaderTimeout int `mapstructure:"reader_timeout"`
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
	}

	viper.AutomaticEnv()

	cfg.Environment = viper.GetString("ENVIRONMENT")
	cfg.DBPath = viper.GetString("DB_PATH")
	cfg.DBFile = viper.GetString("DB_FILE")

	cfg.Telegram.parseEnv(cfg)
	cfg.Sentry.parseEnv(cfg)

	return nil
}
