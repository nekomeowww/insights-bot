package recap

import (
	"fmt"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/pkg/bots/tgbot"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
		fx.Provide(NewRecapCommandHandler()),
		fx.Provide(NewRecapCallbackQueryHandler()),
		fx.Provide(NewMessageHandler()),
		fx.Provide(NewEnableRecapCommandHandler()),
		fx.Provide(NewDisableRecapCommandHandler()),
	)
}

var (
	_ tgbot.CommandHandler = (*RecapCommandHandler)(nil)
)

type NewHandlersParams struct {
	fx.In

	RecapCommand        *RecapCommandHandler
	RecapCallbackQuery  *RecapCallbackQueryHandler
	RecordMessage       *MessageHandler
	EnableRecapCommand  *EnableRecapCommandHandler
	DisableRecapCommand *DisableRecapCommandHandler
}

type Handlers struct {
	recapCommand        *RecapCommandHandler
	recapCallbackQuery  *RecapCallbackQueryHandler
	recordMessage       *MessageHandler
	enableRecapCommand  *EnableRecapCommandHandler
	disableRecapCommand *DisableRecapCommandHandler
}

func NewHandlers() func(NewHandlersParams) *Handlers {
	return func(param NewHandlersParams) *Handlers {
		return &Handlers{
			recapCommand:        param.RecapCommand,
			recapCallbackQuery:  param.RecapCallbackQuery,
			recordMessage:       param.RecordMessage,
			enableRecapCommand:  param.EnableRecapCommand,
			disableRecapCommand: param.DisableRecapCommand,
		}
	}
}

func (h *Handlers) Install(dispatcher *tgbot.Dispatcher) {
	dispatcher.OnCommand(h.recapCommand)
	dispatcher.OnCallbackQuery(h.recapCallbackQuery)
	dispatcher.OnMessage(tgbot.NewHandler(h.recordMessage.HandleRecordMessage))
	dispatcher.OnCommand(h.enableRecapCommand)
	dispatcher.OnCommand(h.disableRecapCommand)
}

var (
	RecapSelectHourAvailables = []int64{
		1, 2, 4, 6, 12,
	}
	RecapSelectHourAvailableText = lo.SliceToMap(RecapSelectHourAvailables, func(item int64) (int64, string) {
		return item, fmt.Sprintf("%d 小时", item)
	})
	RecapSelectHourAvailableValues = lo.SliceToMap(RecapSelectHourAvailables, func(item int64) (int64, string) {
		return item, fmt.Sprintf("%d", item)
	})
)

type RecapCommandHandler struct{}

func NewRecapCommandHandler() func() *RecapCommandHandler {
	return func() *RecapCommandHandler {
		return &RecapCommandHandler{}
	}
}

func (h RecapCommandHandler) Command() string {
	return "recap"
}

func (h RecapCommandHandler) CommandHelp() string {
	return "总结过去的聊天记录并生成回顾快报"
}

func (h *RecapCommandHandler) Handle(c *tgbot.Context) error {
	chatID := c.Update.Message.Chat.ID
	message := tgbotapi.NewMessage(chatID, "要创建过去几个小时内的聊天回顾呢？")
	message.ReplyToMessageID = c.Update.Message.MessageID
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			lo.Map(RecapSelectHourAvailables, func(item int64, _ int) tgbotapi.InlineKeyboardButton {
				return tgbotapi.NewInlineKeyboardButtonData(
					RecapSelectHourAvailableText[item],
					tgbot.NewCallbackQueryData("recap", "select_hour", url.Values{"hour": []string{RecapSelectHourAvailableValues[item]}}),
				)
			})...,
		),
	)

	c.Bot.MustSend(message)
	return nil
}
