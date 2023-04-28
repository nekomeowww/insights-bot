package tgbot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _ CommandHandler = (*helpCommandHandler)(nil)

type helpCommandHandler struct {
	commands []CommandHandler
}

func newHelpCommandHandler() *helpCommandHandler {
	h := &helpCommandHandler{
		commands: make([]CommandHandler, 0),
	}

	return h
}

func (h helpCommandHandler) Command() string {
	return "help"
}

func (h helpCommandHandler) CommandHelp() string {
	return "获取帮助"
}

func (h helpCommandHandler) Handle(c *Context) error {
	helpMessage := strings.Builder{}
	helpMessage.WriteString("你好，欢迎使用 Insights Bot！\n\n")
	helpMessage.WriteString("我当前支持这些命令：\n")

	subCommandHelpMessages := make([]string, 0)
	for _, c := range h.commands {
		subCommandHelpMessage := strings.Builder{}
		subCommandHelpMessage.WriteString("/")
		subCommandHelpMessage.WriteString(c.Command())
		if c.CommandHelp() != "" {
			subCommandHelpMessage.WriteString(" ")
			subCommandHelpMessage.WriteString(c.CommandHelp())
		}

		subCommandHelpMessages = append(subCommandHelpMessages, subCommandHelpMessage.String())
	}
	helpMessage.WriteString(strings.Join(subCommandHelpMessages, "\n"))

	message := tgbotapi.NewMessage(c.Update.Message.Chat.ID, helpMessage.String())
	message.ReplyToMessageID = c.Update.Message.MessageID
	message.ParseMode = "HTML"
	_ = c.Bot.MustSend(message)
	return nil
}
