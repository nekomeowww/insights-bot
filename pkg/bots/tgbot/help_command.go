package tgbot

import (
	"fmt"
	"strings"

	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
)

type Command struct {
	Command     string
	HelpMessage func(*Context) string
	Handler     Handler
}

type commandGroup struct {
	name     func(*Context) string
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

func (h *helpCommandHandler) CommandHelp(c *Context) string {
	return c.T("system.commands.groups.basic.commands.help.help")
}

func (h *helpCommandHandler) handle(c *Context) (Response, error) {
	is, err := c.IsBotAdministrator()
	if err != nil {
		c.Logger.Error("failed to check if bot is administrator")
	}
	if is &&
		c.Update.Message != nil &&
		c.Update.Message.Chat != nil &&
		lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, telegram.ChatType(c.Update.Message.Chat.Type)) &&
		!lo.Contains([]string{
			fmt.Sprintf("%s@%s", h.Command(), c.Bot.Self.UserName),
			fmt.Sprintf("%s@%s", "start", c.Bot.Self.UserName),
		}, c.Update.Message.CommandWithAt()) {
		return nil, nil
	}

	commandGroupHelpMessages := make([]string, 0)

	if len(h.defaultGroup.commands) > 0 {
		h.commandGroups = append(h.commandGroups, h.defaultGroup)
	}

	for _, group := range h.commandGroups {
		commandHelpMessages := make([]string, 0)

		for _, cmd := range group.commands {
			commandHelpMessage := strings.Builder{}

			commandHelpMessage.WriteString("/")
			commandHelpMessage.WriteString(cmd.Command)
			commandHelpMessage.WriteString("@")
			commandHelpMessage.WriteString(c.Bot.Self.UserName)

			if cmd.HelpMessage != nil {
				message := cmd.HelpMessage(c)
				if message != "" {
					commandHelpMessage.WriteString(" - ")
					commandHelpMessage.WriteString(message)
				}
			}

			commandHelpMessages = append(commandHelpMessages, commandHelpMessage.String())
		}

		commandGroupHelpMessages = append(commandGroupHelpMessages, fmt.Sprintf("%s%s", lo.Ternary(
			group.name(c) != "",
			fmt.Sprintf("<b>%s</b>\n\n", EscapeHTMLSymbols(group.name(c))), ""),
			strings.Join(commandHelpMessages, "\n"),
		))
	}

	return c.
		NewMessageReplyTo(
			c.T("system.commands.groups.basic.commands.help.message", i18n.M{
				"Commands": strings.Join(commandGroupHelpMessages, "\n\n"),
			}),
			c.Update.Message.MessageID).
		WithParseModeHTML(), nil
}
