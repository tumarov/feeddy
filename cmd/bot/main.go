package main

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/logger"
	"github.com/tumarov/feeddy/pkg/logger/builtin"
	"github.com/tumarov/feeddy/pkg/repository/mongodb"
	"github.com/tumarov/feeddy/pkg/telegram"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	db, err := initDB(cfg)
	if err != nil {
		log.Exception(err)
		l.Fatal(err)
	}
	log.Debugf("Connected to MongoDB on %s", cfg.DBPath)

	feedRepository := mongodb.NewFeedRepository(db, cfg)

	telegramBot := telegram.NewBot(log, bot, feedRepository, cfg.Messages)

	if err := telegramBot.Start(); err != nil {
		log.Exception(err)
		l.Fatal(err)
	}
}

func initLogger(logger logger.Logger) logger.Logger {
	return logger
}

func initDB(cfg *config.Config) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.DBPath))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
