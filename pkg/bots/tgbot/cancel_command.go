package tgbot

import "github.com/nekomeowww/fo"

type cancellableCommand struct {
	shouldCancelFunc func(c *Context) (bool, error)
	handler          Handler
}

type cancelCommandHandler struct {
	cancellableCommands []cancellableCommand
}

func newCancelCommandHandler() *cancelCommandHandler {
	h := &cancelCommandHandler{
		cancellableCommands: make([]cancellableCommand, 0),
	}

	return h
}

func (h *cancelCommandHandler) Command() string {
	return "cancel"
}

func (h *cancelCommandHandler) CommandHelp() string {
	return "取消正在进行的操作"
}

func (h *cancelCommandHandler) handle(c *Context) (Response, error) {
	may := fo.NewMay[bool]()

	for _, h := range h.cancellableCommands {
		should := may.Invoke(h.shouldCancelFunc(c))
		if should {
			return h.handler.Handle(c)
		}

		continue
	}

	err := may.CollectAsError()
	if err != nil {
		return nil, err
	}

	return c.NewMessageReplyTo("已经没有正在进行的操作了", c.Update.Message.MessageID), nil
}
