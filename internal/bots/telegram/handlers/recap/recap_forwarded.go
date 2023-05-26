package recap

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/models/chathistories"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/redis"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

var (
	_ tgbot.CancellableCommandHandler = (*RecapForwardedStartCommandHandler)(nil)
	_ tgbot.CommandHandler            = (*RecapForwardedCommandHandler)(nil)
)

type NewRecapForwardedStartCommandHandlerParams struct {
	fx.In

	Redis         *datastore.Redis
	ChatHistories *chathistories.Model
	Logger        *logger.Logger
}

type RecapForwardedStartCommandHandler struct {
	redis         *datastore.Redis
	chathistories *chathistories.Model
	logger        *logger.Logger
}

func NewRecapForwardedStartCommandHandler() func(NewRecapForwardedStartCommandHandlerParams) *RecapForwardedStartCommandHandler {
	return func(param NewRecapForwardedStartCommandHandlerParams) *RecapForwardedStartCommandHandler {
		return &RecapForwardedStartCommandHandler{
			redis:         param.Redis,
			chathistories: param.ChatHistories,
			logger:        param.Logger,
		}
	}
}

func (h RecapForwardedStartCommandHandler) Command() string {
	return "recap_forwarded_start"
}

func (h RecapForwardedStartCommandHandler) CommandHelp() string {
	return "使 Bot 接收在私聊中转发给 Bot 的消息，并在发送 /recap_forwarded 后开始总结"
}

func (h *RecapForwardedStartCommandHandler) Handle(c *tgbot.Context) (tgbot.Response, error) {
	has, err := h.chathistories.HasOngoingRecapForwardedFromPrivateMessages(c.Update.Message.From.ID)
	if err != nil {
		return nil, err
	}

	if has {
		delCmd := h.redis.B().
			Del().
			Key(redis.RecapReplayFromPrivateMessageBatch1.Format(c.Update.Message.From.ID)).
			Build()

		err := h.redis.Do(context.Background(), delCmd).Error()
		if err != nil {
			return nil, err
		}
	}

	err = h.chathistories.EnabledRecapForwardedFromPrivateMessages(c.Update.Message.From.ID)
	if err != nil {
		return nil, err
	}

	return c.NewMessageReplyTo("没问题，请将你需要总结的消息在 2 小时内发给我吧。发送完毕后可以通过发送 /recap_forwarded 给我来开始总结哦！", c.Update.Message.MessageID), nil
}

func (h *RecapForwardedStartCommandHandler) ShouldCancel(c *tgbot.Context) (bool, error) {
	has, err := h.chathistories.HasOngoingRecapForwardedFromPrivateMessages(c.Update.Message.From.ID)
	if err != nil {
		return false, err
	}

	return has, nil
}

func (h *RecapForwardedStartCommandHandler) HandleCancel(c *tgbot.Context) (tgbot.Response, error) {
	err := h.chathistories.DisableRecapForwardedFromPrivateMessages(c.Update.Message.From.ID)
	if err != nil {
		return nil, err
	}

	return c.NewMessageReplyTo("好的，已经帮你把消息清理掉了，如果需要总结转发的消息，可以再次发送 /recap_forwarded_start 开始操作。", c.Update.Message.MessageID), nil
}

type NewRecapForwardedCommandHandlerParams struct {
	fx.In

	Redis         *datastore.Redis
	ChatHistories *chathistories.Model
	Logger        *logger.Logger
}

type RecapForwardedCommandHandler struct {
	redis         *datastore.Redis
	chathistories *chathistories.Model
	logger        *logger.Logger
}

func NewRecapForwardedCommandHandler() func(NewRecapForwardedCommandHandlerParams) *RecapForwardedCommandHandler {
	return func(param NewRecapForwardedCommandHandlerParams) *RecapForwardedCommandHandler {
		return &RecapForwardedCommandHandler{
			redis:         param.Redis,
			chathistories: param.ChatHistories,
			logger:        param.Logger,
		}
	}
}

func (h RecapForwardedCommandHandler) Command() string {
	return "recap_forwarded"
}

func (h RecapForwardedCommandHandler) CommandHelp() string {
	return "使 Bot 停止接收在私聊中转发给 Bot 的消息，对已经转发过的消息进行总结"
}

func (h *RecapForwardedCommandHandler) Handle(c *tgbot.Context) (tgbot.Response, error) {
	_, err := c.Bot.Send(tgbotapi.NewMessage(
		c.Update.Message.From.ID,
		"正在为已经接收到的聊天记录生成回顾，请稍等...",
	))
	if err != nil {
		h.logger.Error("failed to send message")
	}

	histories, err := h.chathistories.FindPrivateForwardedChatHistories(c.Update.Message.From.ID)
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("聊天记录回顾生成失败，请稍后再试！").WithReply(c.Update.Message)
	}
	if len(histories) < 5 {
		return nil, tgbot.NewMessageError("目前收到的聊天记录不足 5 条哦，要再多发送给我一些之后之后再试试吗？").WithReply(c.Update.Message)
	}

	summarizations, err := h.chathistories.SummarizePrivateForwardedChatHistories(c.Update.Message.From.ID, histories)
	if err != nil {
		return nil, tgbot.NewExceptionError(err).WithMessage("聊天记录回顾生成失败，请稍后再试！").WithReply(c.Update.Message)
	}

	summarizations = lo.Filter(summarizations, func(item string, _ int) bool { return item != "" })
	if len(summarizations) == 0 {
		return nil, tgbot.NewExceptionError(err).WithMessage("聊天记录回顾生成失败，请稍后再试！").WithReply(c.Update.Message)
	}

	for i, s := range summarizations {
		summarizations[i] = tgbot.ReplaceMarkdownTitlesToTelegramBoldElement(s)
	}

	summarizationBatches := tgbot.SplitMessagesAgainstLengthLimitIntoMessageGroups(summarizations)

	for i, s := range summarizationBatches {
		var content string
		if len(summarizationBatches) > 1 {
			content = fmt.Sprintf("%s\n\n(%d/%d)\n#recap\n<em>🤖️ Generated by chatGPT</em>", strings.Join(s, "\n\n"), i+1, len(summarizationBatches))
		} else {
			content = fmt.Sprintf("%s\n\n#recap\n<em>🤖️ Generated by chatGPT</em>", strings.Join(s, "\n\n"))
		}

		msg := tgbotapi.NewMessage(c.Update.Message.Chat.ID, content)
		msg.ParseMode = tgbotapi.ModeHTML
		msg.ReplyToMessageID = c.Update.Message.MessageID

		h.logger.Infof("sending chat histories recap for chat %d: %s", c.Update.Message.Chat.ID, msg.Text)

		c.Bot.MustSend(msg)
	}

	msg := tgbotapi.NewMessage(c.Update.Message.Chat.ID, "总结完成，如果你觉得不满意，可以再次发送 /recap_forwarded 重新生成哦！如果觉得满意，并且希望进行其他的总结操作，可以在开始前发送 /cancel 来清空当前已经接收到的消息，如果不进行操作，缓存的消息将会在 2 小时后被自动清理。")
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyToMessageID = c.Update.Message.MessageID

	return nil, nil
}
