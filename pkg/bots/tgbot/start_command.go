package tgbot

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

func (h *startCommandHandler) CommandHelp() string {
	return "开始与 Bot 的交互"
}

func (h *startCommandHandler) handle(c *Context) (Response, error) {
	for _, h := range h.startCommandHandlers {
		_, _ = h.Handle(c)
		if c.IsAborted() {
			return nil, nil
		}

		continue
	}

	return h.helpCommandHandler.handle(c)
}
