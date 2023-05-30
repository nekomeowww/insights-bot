package tgbot

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Command struct {
	Command     string
	HelpMessage string
	Handler     Handler
}

type commandGroup struct {
	name     string
	commands []Command
}

type helpCommandHandler struct {
	defaultGroup  commandGroup
	commandGroups []commandGroup
}

func newHelpCommandHandler() *helpCommandHandler {
	h := &helpCommandHandler{
		commandGroups: make([]commandGroup, 0),
	}

	return h
}

func (h *helpCommandHandler) Command() string {
	return "help"
}

func (h *helpCommandHandler) CommandHelp() string {
	return "获取帮助"
}

func (h *helpCommandHandler) handle(c *Context) (Response, error) {
	helpMessage := strings.Builder{}
	helpMessage.WriteString("你好，欢迎使用 Insights Bot！\n\n")
	helpMessage.WriteString("我当前支持这些命令：\n\n")

	commandGroupHelpMessages := make([]string, 0)

	if len(h.defaultGroup.commands) > 0 {
		h.commandGroups = append(h.commandGroups, h.defaultGroup)
	}

	for _, group := range h.commandGroups {
		commandHelpMessages := make([]string, 0)

		for _, c := range group.commands {
			commandHelpMessage := strings.Builder{}

			commandHelpMessage.WriteString("/")
			commandHelpMessage.WriteString(c.Command)

			if c.HelpMessage != "" {
				commandHelpMessage.WriteString(" - ")
				commandHelpMessage.WriteString(c.HelpMessage)
			}

			commandHelpMessages = append(commandHelpMessages, commandHelpMessage.String())
		}

		commandGroupHelpMessages = append(commandGroupHelpMessages, fmt.Sprintf("%s%s", lo.Ternary(group.name != "", fmt.Sprintf("<b>%s</b>\n\n", group.name), ""), strings.Join(commandHelpMessages, "\n")))
	}

	helpMessage.WriteString(strings.Join(commandGroupHelpMessages, "\n\n"))

	return c.NewMessageReplyTo(helpMessage.String(), c.Update.Message.MessageID).WithParseModeHTML(), nil
}
