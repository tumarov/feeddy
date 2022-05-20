package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/logger"
	"github.com/tumarov/feeddy/pkg/logger/builtin"
	rssParser "github.com/tumarov/feeddy/pkg/parser"
	"github.com/tumarov/feeddy/pkg/repository/boltdb"
	"github.com/tumarov/feeddy/pkg/telegram"
	bolt "go.etcd.io/bbolt"
	l "log"
	"time"
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

	db, err := bolt.Open(cfg.DBFile, 0666, &bolt.Options{Timeout: 2 * time.Second})
	defer db.Close()
	if err != nil {
		log.Exception(err)
		l.Fatal(err)
	}
	log.Debugf("Connected to database")

	feedRepository, err := boltdb.NewFeedRepository(db, cfg)
	if err != nil {
		log.Exception(err)
		l.Fatal(err)
	}

	telegramBot := telegram.NewBot(log, bot, feedRepository, cfg.Messages)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Exception(err)
			l.Fatal(err)
		}
	}()

	parser := rssParser.NewRSSReader(feedRepository, cfg, telegramBot)

	if err := parser.Start(); err != nil {
		log.Exception(err)
		l.Fatal(err)
	}
}

func initLogger(logger logger.Logger) logger.Logger {
	return logger
}
