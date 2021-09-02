package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidURL   = errors.New("url is invalid")
	errUnauthorized = errors.New("user is not authorized")
	errUnableToSave = errors.New("unable to save")
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, b.messages.Default)

	switch err {
	case errInvalidURL:
		msg.Text = b.messages.InvalidURL
	case errUnauthorized:
		msg.Text = b.messages.Unauthorized
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
	default:
	}

	_, err = b.bot.Send(msg)
	if err != nil {
		b.log.Exception(err)
	}
}
