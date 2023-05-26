package tgbot

import "github.com/nekomeowww/fo"

var _ CommandHandler = (*cancelCommandHandler)(nil)

type cancelCommandHandler struct {
	commands []CancellableCommandHandler
}

func newCancelCommandHandler() *cancelCommandHandler {
	h := &cancelCommandHandler{
		commands: make([]CancellableCommandHandler, 0),
	}

	return h
}

func (h cancelCommandHandler) Command() string {
	return "cancel"
}

func (h cancelCommandHandler) CommandHelp() string {
	return "取消正在进行的操作"
}

func (h cancelCommandHandler) Handle(c *Context) (Response, error) {
	may := fo.NewMay[bool]()

	for _, h := range h.commands {
		should := may.Invoke(h.ShouldCancel(c))
		if should {
			return h.HandleCancel(c)
		}

		continue
	}

	err := may.CollectAsError()
	if err != nil {
		return nil, err
	}

	return c.NewMessageReplyTo("已经没有正在进行的操作了", c.Update.Message.MessageID), nil
}
