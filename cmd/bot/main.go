package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/logger"
	"github.com/tumarov/feeddy/pkg/logger/builtin"
	"github.com/tumarov/feeddy/pkg/telegram"
	l "log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		l.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		l.Fatal(err)
	}
	bot.Debug = cfg.Telegram.Debug

	log := initLogger(builtin.NewBuiltinLogger())

	telegramBot := telegram.NewBot(log, bot, cfg.Messages)

	if err := telegramBot.Start(); err != nil {
		log.Exception(err)
		l.Fatal(err)
	}
}

func initLogger(logger logger.Logger) logger.Logger {
	return logger
}
