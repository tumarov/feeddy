package config

import "github.com/spf13/viper"

type Telegram struct {
	Token  string
	BotURL string `mapstructure:"bot_url"`
	Debug  bool
}

func (t *Telegram) init(cfg *Config) error {
	if err := viper.Unmarshal(&cfg.Telegram); err != nil {
		return err
	}

	return nil
}

func (t *Telegram) parseEnv(cfg *Config) {
	cfg.Telegram.Token = viper.GetString("TELEGRAM_BOT_TOKEN")
	cfg.Telegram.Debug = cfg.Environment == "dev"
}
