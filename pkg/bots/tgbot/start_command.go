package tgbot

import (
	"fmt"

	"github.com/nekomeowww/insights-bot/pkg/types/telegram"
	"github.com/samber/lo"
)

type startCommandHandler struct {
	helpCommandHandler   *helpCommandHandler
	startCommandHandlers []Handler
}

func newStartCommandHandler() *startCommandHandler {
	h := &startCommandHandler{
		startCommandHandlers: make([]Handler, 0),
	}

	return h
}

func (h *startCommandHandler) Command() string {
	return "start"
}

func (h *startCommandHandler) CommandHelp(c *Context) string {
	return c.T("system.commands.groups.basic.commands.start.help")
}

func (h *startCommandHandler) handle(c *Context) (Response, error) {
	is, err := c.IsBotAdministrator()
	if err != nil {
		c.Logger.Error("failed to check if bot is administrator")
	}
	if is &&
		c.Update.Message != nil &&
		c.Update.Message.Chat != nil &&
		lo.Contains([]telegram.ChatType{telegram.ChatTypeGroup, telegram.ChatTypeSuperGroup}, telegram.ChatType(c.Update.Message.Chat.Type)) &&
		c.Update.Message.CommandWithAt() != fmt.Sprintf("%s@%s", h.Command(), c.Bot.Self.UserName) {
		return nil, nil
	}

	for _, h := range h.startCommandHandlers {
		_, _ = h.Handle(c)
		if c.IsAborted() {
			return nil, nil
		}

		continue
	}

	return h.helpCommandHandler.handle(c)
}
