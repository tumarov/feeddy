package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
	"strings"
)

const (
	commandStart  = "start"
	commandList   = "list"
	commandRemove = "remove"
	commandHelp   = "help"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	case commandList:
		return b.handleListCommand(message)
	case commandRemove:
		return b.handleRemoveCommand(message)
	case commandHelp:
		return b.handleHelpCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	found, err := b.feedRepository.Get(message.Chat.ID)
	if err != nil {
		return nil
	}

	if found != nil {
		return errAlreadyAuthorized
	}

	if err := b.feedRepository.Create(message.Chat.ID); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Start)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleListCommand(message *tgbotapi.Message) error {
	found, err := b.feedRepository.Get(message.Chat.ID)
	if err != nil {
		return nil
	}

	if found == nil {
		return errUnauthorized
	}

	var response string

	for _, feed := range found.Feeds {
		response += feed.URL + "\n"
	}

	if response == "" {
		response = b.messages.ListEmpty
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	_, err = b.bot.Send(msg)

	return err
}

func (b *Bot) handleRemoveCommand(message *tgbotapi.Message) error {
	found, err := b.feedRepository.Get(message.Chat.ID)
	if err != nil {
		return nil
	}

	if found == nil {
		return errUnauthorized
	}

	if err = b.feedRepository.Remove(*found, strings.ReplaceAll(message.Text, "/"+message.Command()+" ", "")); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.RemovedSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Help)
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

	found, err := b.feedRepository.Get(message.Chat.ID)
	if err != nil {
		return nil
	}

	if found == nil {
		return errUnauthorized
	}

	if err = b.feedRepository.Add(*found, message.Text); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.SavedSuccessfully)
	_, err = b.bot.Send(msg)
	return err
}
