package tgbot

import (
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type HandlerGroup interface {
	Install(dispatcher *Dispatcher)
}

type CommandHandler interface {
	Handle(c *Context) error
	Command() string
	CommandHelp() string
}

type CallbackQueryHandler interface {
	Handle(c *Context) error
	CallbackQueryRoute() string
}

type Handler interface {
	Handle(c *Context) error
}

type defaultHandler struct {
	handleFunc HandleFunc
}

func (h defaultHandler) Handle(c *Context) error {
	return h.handleFunc(c)
}

type HandleFunc func(ctx *Context) error

type MiddlewareFunc func(ctx *Context, next func()) HandleFunc

func isErrorRepliable(updateType UpdateType) bool {
	return lo.Contains([]UpdateType{
		UpdateTypeMessage,
	}, updateType)
}

func processError(ctx *Context, err error) {
	chatID := ctx.Update.FromChat().ID

	switch v := err.(type) {
	case MessageError:
		if isErrorRepliable(ctx.UpdateType()) {
			if chatID == 0 {
				return
			}

			message := tgbotapi.NewMessage(chatID, v.message)
			if v.editMessageID != 0 {
				ctx.Bot.MustEditMessageText(chatID, v.editMessageID, v.message)
				return
			}
			if v.replyToMessageID != 0 {
				message.ReplyToMessageID = v.replyToMessageID
			}

			_ = ctx.Bot.MustSend(message)
		}
	case ExceptionError:
		if (v.replyToMessageID != 0 || v.editMessageID != 0) && isErrorRepliable(ctx.UpdateType()) {
			if chatID == 0 {
				return
			}

			message := tgbotapi.NewMessage(ctx.Update.FromChat().ID, lo.Ternary(v.message != "", v.message, "发生了一些错误，请稍后再试"))
			if v.editMessageID != 0 {
				ctx.Bot.MustEditMessageText(chatID, v.editMessageID, v.message)
				return
			}
			if v.replyToMessageID != 0 {
				message.ReplyToMessageID = v.replyToMessageID
			}

			_ = ctx.Bot.MustSend(message)
		}

		entry := logrus.NewEntry(ctx.Logger.Logger)
		logger.SetCallerFrameWithFileAndLine(entry, "insights-bot", v.callFrame.Function, v.callFrame.File, v.callFrame.Line)
		entry.WithFields(logrus.Fields{
			"chat_id":   ctx.Update.FromChat().ID,
			"update_id": ctx.Update.UpdateID,
			"message":   v.message,
			"error":     v.err,
		}).Error(v.err)
	default:
		if isErrorRepliable(ctx.UpdateType()) {
			_ = ctx.Bot.MustSend(tgbotapi.NewMessage(ctx.Update.FromChat().ID, "发生了一些错误，请稍后再试"))
		}

		ctx.Logger.WithFields(logrus.Fields{
			"chat_id":   ctx.Update.FromChat().ID,
			"update_id": ctx.Update.UpdateID,
		}).Errorf("encountered unknown error: %v, stack: %s", err, debug.Stack())
	}
}

func NewHandler(h HandleFunc) Handler {
	wrapped := func(ctx *Context) error {
		err := h(ctx)
		if err != nil {
			processError(ctx, err)
		}

		return nil
	}

	return defaultHandler{
		handleFunc: wrapped,
	}
}
