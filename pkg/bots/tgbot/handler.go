package tgbot

import (
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type Response interface{}

type HandlerGroup interface {
	Install(dispatcher *Dispatcher)
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

func isErrorCanBeReplied(updateType UpdateType) bool {
	return lo.Contains([]UpdateType{
		UpdateTypeMessage,
		UpdateTypeCallbackQuery,
	}, updateType)
}

func processMessageError(ctx *Context, chatID int64, msgError MessageError) Response {
	if msgError.message == "" {
		return nil
	}

	if msgError.editMessage != nil && msgError.editMessage.MessageID != 0 {
		returns := NewEditMessageText(chatID, msgError.editMessage.MessageID, msgError.message)
		if msgError.parseMode == tgbotapi.ModeHTML {
			returns = returns.WithParseModeHTML()
		}
		if msgError.replyMarkup != nil {
			returns = returns.WithInlineReplyMarkup(*msgError.replyMarkup)
		}

		textIsTheSame := msgError.editMessage.Text == RemoveHTMLBlocksFromString(msgError.message) || msgError.editMessage.Caption == RemoveHTMLBlocksFromString(msgError.message)
		inlineKeyboardsAreTheSame := true

		if msgError.editMessage.ReplyMarkup.InlineKeyboard != nil && msgError.replyMarkup != nil &&
			len(msgError.editMessage.ReplyMarkup.InlineKeyboard) == len(msgError.replyMarkup.InlineKeyboard) {
			for i := range msgError.editMessage.ReplyMarkup.InlineKeyboard {
				diff1, diff2 := lo.Difference(msgError.editMessage.ReplyMarkup.InlineKeyboard[i], msgError.replyMarkup.InlineKeyboard[i])
				inlineKeyboardsAreTheSame = !(len(diff1) != 0 || len(diff2) != 0)

				break
			}
		}
		if textIsTheSame && inlineKeyboardsAreTheSame {
			return nil
		}

		ctx.Abort()

		return returns
	}
	if msgError.replyToMessageID != 0 {
		returns := NewMessageReplyTo(chatID, msgError.message, msgError.replyToMessageID)
		if msgError.parseMode == tgbotapi.ModeHTML {
			returns = returns.WithParseModeHTML()
		}
		if msgError.replyMarkup != nil {
			returns = returns.WithReplyMarkup(*msgError.replyMarkup)
		}
		if msgError.deleteLaterForUserID != 0 && msgError.deleteLaterChatID != 0 {
			returns = returns.WithDeleteLater(msgError.deleteLaterForUserID, msgError.deleteLaterChatID)
		}

		ctx.Abort()

		return returns
	}

	returns := NewMessage(chatID, msgError.message)
	if msgError.parseMode == tgbotapi.ModeHTML {
		returns = returns.WithParseModeHTML()
	}
	if msgError.replyMarkup != nil {
		returns = returns.WithReplyMarkup(*msgError.replyMarkup)
	}
	if msgError.deleteLaterForUserID != 0 && msgError.deleteLaterChatID != 0 {
		returns = returns.WithDeleteLater(msgError.deleteLaterForUserID, msgError.deleteLaterChatID)
	}

	ctx.Abort()

	return returns
}

func processExceptionError(ctx *Context, chatID int64, e ExceptionError) Response {
	entry := logrus.NewEntry(ctx.Logger.Logger)
	logger.SetCallerFrameWithFileAndLine(entry, "insights-bot", e.callFrame.Function, e.callFrame.File, e.callFrame.Line)
	entry.WithFields(logrus.Fields{
		"update_type":      ctx.UpdateType(),
		"chat_id":          ctx.Update.FromChat().ID,
		"update_id":        ctx.Update.UpdateID,
		"message":          e.message,
		"error":            e.err,
		"edit_message_id":  e.editMessage.MessageID,
		"reply_message_id": e.replyToMessageID,
	}).Errorf("encountered an exception error: %v", e.err)

	message := lo.Ternary(e.message != "", e.message, "发生了一些错误，请稍后再试")
	if e.editMessage != nil && e.editMessage.MessageID != 0 {
		return NewEditMessageText(chatID, e.editMessage.MessageID, message)
	}
	if e.replyToMessageID != 0 {
		returns := NewMessageReplyTo(chatID, message, e.replyToMessageID)
		if e.deleteLaterForUserID != 0 && e.deleteLaterChatID != 0 {
			returns = returns.WithDeleteLater(e.deleteLaterForUserID, e.deleteLaterChatID)
		}

		return returns
	}

	returns := NewMessage(chatID, message)
	if e.deleteLaterForUserID != 0 && e.deleteLaterChatID != 0 {
		returns = returns.WithDeleteLater(e.deleteLaterForUserID, e.deleteLaterChatID)
	}

	ctx.Abort()

	return returns
}

func processError(ctx *Context, err error) Response {
	if !isErrorCanBeReplied(ctx.UpdateType()) {
		ctx.Logger.WithFields(logrus.Fields{
			"update_type": ctx.UpdateType(),
			"chat_id":     ctx.Update.FromChat().ID,
			"update_id":   ctx.Update.UpdateID,
		}).Errorf("error occurred when handling response: %v", err)

		return nil
	}

	chatID := ctx.Update.FromChat().ID
	if chatID == 0 {
		return nil
	}

	switch v := err.(type) {
	case MessageError:
		return processMessageError(ctx, chatID, v)
	case ExceptionError:
		return processExceptionError(ctx, chatID, v)
	default:
		ctx.Logger.WithFields(logrus.Fields{
			"update_type": ctx.UpdateType(),
			"chat_id":     ctx.Update.FromChat().ID,
			"update_id":   ctx.Update.UpdateID,
		}).Errorf("encountered unknown error: %v, stack: %s", err, debug.Stack())

		ctx.Abort()

		return NewMessage(chatID, "发生了一些错误，请稍后再试")
	}
}

func processResponse(ctx *Context, resp Response) {
	if resp == nil {
		return
	}

	switch v := resp.(type) {
	case MessageResponse:
		ctx.Abort()

		msg := ctx.Bot.MaySend(v.messageConfig)
		if msg != nil && v.deleteLaterForUserID != 0 && v.deleteLaterChatID != 0 {
			err := ctx.Bot.PushOneDeleteLaterMessage(v.deleteLaterForUserID, v.deleteLaterChatID, msg.MessageID)
			if err != nil {
				ctx.Logger.Errorf("failed to push delete later message: %v", err)
			}
		}
	case EditMessageResponse:
		ctx.Abort()

		var err error
		if v.mediaConfig != nil {
			_, err = ctx.Bot.Request(v.mediaConfig)
		}
		if v.replyMarkupConfig != nil {
			_, err = ctx.Bot.Request(v.replyMarkupConfig)
		}
		if v.liveLocationConfig != nil {
			_, err = ctx.Bot.Request(v.liveLocationConfig)
		}
		if v.textConfig != nil {
			_, err = ctx.Bot.Request(v.textConfig)
		}
		if v.captionConfig != nil {
			_, err = ctx.Bot.Request(v.captionConfig)
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
		}

		processResponse(ctx, resp)

		return nil, nil
	}

	return defaultHandler{
		handleFunc: wrapped,
	}
}
