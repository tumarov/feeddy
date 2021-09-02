package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
)

const (
	commandStart = "start"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	// TODO: some logic

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Start)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidURL
	}

	b.log.Debugf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}
