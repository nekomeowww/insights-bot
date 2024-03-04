package recap

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/fo"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/nekomeowww/insights-bot/pkg/types/redis"
	"github.com/redis/rueidis"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type privateSubscriptionStartCommandContext struct {
	ChatID    int64  `json:"chat_id"`
	ChatTitle string `json:"chat_title"`
}

func (h *CommandHandler) setRecapForPrivateSubscriptionModeStartCommandContext(chatID int64, chatTitle string) (string, error) {
	hashSource := fmt.Sprintf("recap/private_subscription_mode/start_command_context/%d", chatID)
	hashKey := fmt.Sprintf("%x", sha256.Sum256([]byte(hashSource)))[0:8]

	setCmd := h.redis.Client.B().
		Set().
		Key(redis.RecapPrivateSubscriptionStartCommandContext1.Format(hashKey)).
		Value(string(lo.Must(json.Marshal(privateSubscriptionStartCommandContext{
			ChatID:    chatID,
			ChatTitle: chatTitle,
		})))).
		ExSeconds(24 * 60 * 60).
		Build()

	err := h.redis.Do(context.Background(), setCmd).Error()
	if err != nil {
		return hashKey, err
	}

	return hashKey, nil
}

func (h *CommandHandler) getRecapForPrivateSubscriptionModeStartCommandContext(hash string) (*privateSubscriptionStartCommandContext, error) {
	getCmd := h.redis.Client.B().
		Get().
		Key(redis.RecapPrivateSubscriptionStartCommandContext1.Format(hash)).
		Build()

	str, err := h.redis.Do(context.Background(), getCmd).ToString()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return nil, nil
		}

		return nil, err
	}
	if str == "" {
		return nil, nil
	}

	var data privateSubscriptionStartCommandContext

	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (h *CommandHandler) setSubscribeStartCommandContext(chatID int64, chatTitle string) (string, error) {
	hashSource := fmt.Sprintf("recap/subscribe_recap/start_command_context/%d", chatID)
	hashKey := fmt.Sprintf("%x", sha256.Sum256([]byte(hashSource)))[0:8]

	setCmd := h.redis.Client.B().
		Set().
		Key(redis.RecapSubscribeRecapStartCommandContext1.Format(hashKey)).
		Value(string(lo.Must(json.Marshal(privateSubscriptionStartCommandContext{
			ChatID:    chatID,
			ChatTitle: chatTitle,
		})))).
		ExSeconds(24 * 60 * 60).
		Build()

	err := h.redis.Do(context.Background(), setCmd).Error()
	if err != nil {
		return hashKey, err
	}

	return hashKey, nil
}

func (h *CommandHandler) getSubscribeStartCommandContext(hash string) (*privateSubscriptionStartCommandContext, error) {
	getCmd := h.redis.Client.B().
		Get().
		Key(redis.RecapSubscribeRecapStartCommandContext1.Format(hash)).
		Build()

	str, err := h.redis.Do(context.Background(), getCmd).ToString()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return nil, nil
		}

		return nil, err
	}
	if str == "" {
		return nil, nil
	}

	var data privateSubscriptionStartCommandContext

	err = json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func newRecapCommandWhenUserNeverStartedChat(bot *tgbot.Bot, hashKey string) string {
	return fmt.Sprintf(""+
		"抱歉，在给您发送引导您创建聊天回顾的消息时出现了问题，这似乎是因为您<b>从未</b>和本 Bot（@%s） "+
		"<b>发起过对话</b>导致的。\n\n"+
		"由于当前群组的聊天回顾功能已经被<b>群组创建者</b>设定为<b>私聊订阅模式</b>，Bot 需要通过私聊的方"+
		"式向您发送引导您创建聊天回顾的消息，届时，您需要完成以下任一一个操作后方可继续创建聊天回顾：\n"+
		"1. <b>点击链接</b> https://t.me/%s?start=%s 与 Bot 开始对话就能继续原先的 /recap 命令操作"+
		"；\n"+
		"2. 点击 Bot 头像并且开始对话，然后在群组内重新发送 /recap 命令来创建聊天回顾。"+
		"", bot.Self.UserName, bot.Self.UserName, hashKey)
}

func newSubscribeRecapCommandWhenUserNeverStartedChat(bot *tgbot.Bot, hashKey string) string {
	return fmt.Sprintf(""+
		"抱歉，在为您订阅本群组定时聊天回顾时出现了问题，这似乎是因为您<b>从未</b>和本 Bot（@%s） <b>发起"+
		"过对话</b>导致的。\n\n"+
		"订阅群组的聊天回顾需要 Bot 需要有权限通过私聊的方式向您定期发送聊天回顾，届时，您需要完成以下任一一"+
		"个操作后方可完成订阅：\n"+
		"1. <b>点击链接</b> https://t.me/%s?start=%s 与 Bot 开始对话；\n"+
		"2. 点击 Bot 头像并且开始对话，然后在群组内重新发送 /subscribe_recap 命令来订阅本群组的定时聊"+
		"天回顾。"+
		"", bot.Self.UserName, bot.Self.UserName, hashKey)
}

func newRecapCommandWhenUserBlockedMessage(bot *tgbot.Bot, hashKey string) string {
	return fmt.Sprintf(""+
		"抱歉，在给您发送引导您创建聊天回顾的消息时出现了问题，这似乎是因为您已将本 Bot（@%s）<b>停用</b>"+
		"或是添加到了<b>黑名单</b>中导致的。\n\n"+
		"由于当前群组的聊天回顾功能已经被<b>群组创建者</b>设定为<b>私聊订阅模式</b>，Bot 需要通过私聊的方"+
		"式向您发送引导您创建聊天回顾的消息，届时，您需要根据下面的提示进行操作：\n"+
		"1. 将 Bot 从<b>黑名单中移除</b>；\n"+
		"2. <b>点击链接</b> https://t.me/%s?start=%s 继续创建聊天回顾，或是在群组内重新发送 /recap "+
		"命令来创建聊天回顾。"+
		"", bot.Self.UserName, bot.Self.UserName, hashKey)
}

func newSubscribeRecapCommandWhenUserBlockedMessage(bot *tgbot.Bot, hashKey string) string {
	return fmt.Sprintf(""+
		"抱歉，在为您订阅本群组定时聊天回顾时出现了问题，这似乎是因为您已将本 Bot（@%s）<b>停用</b>或是添加"+
		"到了<b>黑名单</b>中导致的。\n\n"+
		"订阅群组的聊天回顾需要 Bot 需要有权限通过私聊的方式向您定期发送聊天回顾，届时，您需要根据下面的提示"+
		"进行操作：\n"+
		"1. 将 Bot 从<b>黑名单中移除</b>；\n"+
		"2. <b>点击链接</b> https://t.me/%s?start=%s 继续订阅本群组的定时聊天回顾操作，或是在群组内重新"+
		"发送 /subscribe_recap 命令来订阅本群组的定时聊天回顾。"+
		"", bot.Self.UserName, bot.Self.UserName, hashKey)
}

func (h *CommandHandler) handleUserNeverStartedChatOrBlockedErr(c *tgbot.Context, chatID int64, _ string, message string) (tgbot.Response, error) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = c.Update.Message.MessageID
	msg.ParseMode = tgbotapi.ModeHTML

	sentMsg := c.Bot.MaySend(msg)

	may := fo.NewMay0().Use(func(err error, messageArgs ...any) {
		h.logger.Error("failed to push one delete later message", zap.Error(err))
	})

	may.Invoke(c.Bot.PushOneDeleteLaterMessage(c.Update.Message.From.ID, chatID, c.Update.Message.MessageID))
	may.Invoke(c.Bot.PushOneDeleteLaterMessage(c.Update.Message.From.ID, chatID, sentMsg.MessageID))

	return nil, nil
}
