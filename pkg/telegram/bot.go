package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tumarov/feeddy/pkg/config"
	"github.com/tumarov/feeddy/pkg/logger"
	"github.com/tumarov/feeddy/pkg/repository"
)

type RSSEntry struct {
	Title       string
	Link        string
	Description string
	Published   string
}

type Bot struct {
	bot            *tgbotapi.BotAPI
	log            logger.Logger
	feedRepository repository.FeedRepository

	messages config.Messages
}

func NewBot(log logger.Logger, bot *tgbotapi.BotAPI, fr repository.FeedRepository, messages config.Messages) *Bot {
	return &Bot{log: log, bot: bot, feedRepository: fr, messages: messages}
}

func (b *Bot) Start() error {
	b.log.Debugf("Starting bot %s ...", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			b.log.Debugf("emtpy message from user")
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
}

func (b *Bot) SendFeeds(chatID int64, feeds []RSSEntry) {
	for _, feed := range feeds {
		msg := tgbotapi.NewMessage(chatID, feed.Link)
		_, err := b.bot.Send(msg)
		if err != nil {
			continue
		}
	}
}
