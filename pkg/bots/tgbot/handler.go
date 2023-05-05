package tgbot

import (
	"fmt"
	"runtime/debug"

	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type Response interface{}

type HandlerGroup interface {
	Install(dispatcher *Dispatcher)
}

type CommandHandler interface {
	Handler
	Command() string
	CommandHelp() string
}

type MessageHandler interface {
	Handler
	Message() string
}

type CallbackQueryHandler interface {
	Handler
	CallbackQueryRoute() string
}

type Handler interface {
	Handle(c *Context) (Response, error)
}

type defaultHandler struct {
	handleFunc HandleFunc
}

func (h defaultHandler) Handle(c *Context) (Response, error) {
	return h.handleFunc(c)
}

type HandleFunc func(ctx *Context) (Response, error)

type MiddlewareFunc func(ctx *Context, next func())

func isErrorRepliable(updateType UpdateType) bool {
	return lo.Contains([]UpdateType{
		UpdateTypeMessage,
		UpdateTypeCallbackQuery,
	}, updateType)
}

func processError(ctx *Context, err error) Response {
	chatID := ctx.Update.FromChat().ID

	switch v := err.(type) {
	case MessageError:
		ctx.Logger.WithFields(logrus.Fields{
			"update_type":      ctx.UpdateType(),
			"chat_id":          ctx.Update.FromChat().ID,
			"update_id":        ctx.Update.UpdateID,
			"message":          v.message,
			"edit_message_id":  v.editMessageID,
			"reply_message_id": v.replyToMessageID,
		}).Infof("returned a message error: %s", v.message)
		if !isErrorRepliable(ctx.UpdateType()) {
			return nil
		}
		if chatID == 0 {
			return nil
		}
		if v.message == "" {
			return nil
		}
		if v.editMessageID != 0 {
			return NewEditMessageText(chatID, v.editMessageID, v.message)
		}
		if v.replyToMessageID != 0 {
			return NewMessageReplyTo(chatID, v.message, v.replyToMessageID)
		}

		return NewMessage(chatID, v.message)
	case ExceptionError:
		entry := logrus.NewEntry(ctx.Logger.Logger)
		logger.SetCallerFrameWithFileAndLine(entry, "insights-bot", v.callFrame.Function, v.callFrame.File, v.callFrame.Line)
		entry.WithFields(logrus.Fields{
			"update_type":      ctx.UpdateType(),
			"chat_id":          ctx.Update.FromChat().ID,
			"update_id":        ctx.Update.UpdateID,
			"message":          v.message,
			"error":            v.err,
			"edit_message_id":  v.editMessageID,
			"reply_message_id": v.replyToMessageID,
		}).Errorf("encountered an exception error: %v", v.err)
		if !isErrorRepliable(ctx.UpdateType()) || (v.replyToMessageID == 0 && v.editMessageID == 0) {
			return nil
		}
		if chatID == 0 {
			return nil
		}

		message := lo.Ternary(v.message != "", v.message, "发生了一些错误，请稍后再试")
		if v.editMessageID != 0 {
			return NewEditMessageText(chatID, v.editMessageID, message)
		}
		if v.replyToMessageID != 0 {
			return NewMessageReplyTo(chatID, message, v.replyToMessageID)
		}

		return NewMessage(chatID, message)
	default:
		ctx.Logger.WithFields(logrus.Fields{
			"update_type": ctx.UpdateType(),
			"chat_id":     ctx.Update.FromChat().ID,
			"update_id":   ctx.Update.UpdateID,
		}).Errorf("encountered unknown error: %v, stack: %s", err, debug.Stack())
		if isErrorRepliable(ctx.UpdateType()) {
			return NewMessage(chatID, "发生了一些错误，请稍后再试")
		}
	}

	return nil
}

func processResponse(ctx *Context, resp Response) {
	if resp == nil {
		return
	}

	switch v := resp.(type) {
	case MessageResponse:
		_ = ctx.Bot.MustSend(v.messageConfig)
	case EditMessageResponse:
		var err error
		if v.textConfig != nil {
			_, err = ctx.Bot.Request(v.textConfig)
		}
		if v.mediaConfig != nil {
			_, err = ctx.Bot.Request(v.mediaConfig)
		}
		if v.replyMarkupConfig != nil {
			_, err = ctx.Bot.Request(v.replyMarkupConfig)
		}
		if v.captionConfig != nil {
			_, err = ctx.Bot.Request(v.captionConfig)
		}
		if v.liveLocationConfig != nil {
			_, err = ctx.Bot.Request(v.liveLocationConfig)
		}
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{
				"chat_id": ctx.Update.FromChat().ID,
			}).Errorf("failed to edit message %v: %v", v, err)
		}
	default:
		ctx.Logger.WithFields(logrus.Fields{
			"chat_id": ctx.Update.FromChat().ID,
		}).Errorf("encountered unknown response %T: %v", v, resp)
	}
}

func NewHandler(h HandleFunc) Handler {
	wrapped := func(ctx *Context) (Response, error) {
		resp, err := h(ctx)
		if err != nil {
			resp = processError(ctx, err)
			fmt.Println(resp)
		}

		processResponse(ctx, resp)

		return nil, nil
	}

	return defaultHandler{
		handleFunc: wrapped,
	}
}
