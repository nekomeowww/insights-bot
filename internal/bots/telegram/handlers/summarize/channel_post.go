package summarize

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/imroc/req/v3"
	"go.uber.org/fx"

	"github.com/nekomeowww/insights-bot/pkg/handler"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/openai"
)

type NewHandlerParam struct {
	fx.In

	Logger *logger.Logger
	OpenAI *openai.Client
}

type Handler struct {
	Logger *logger.Logger

	ReqClient *req.Client
	OpenAI    *openai.Client
}

func NewHandler() func(param NewHandlerParam) *Handler {
	return func(param NewHandlerParam) *Handler {
		handler := &Handler{
			Logger:    param.Logger,
			ReqClient: req.C(),
			OpenAI:    param.OpenAI,
		}
		return handler
	}
}

func (h *Handler) HandleChannelPost(c *handler.Context) {
	// 转发的消息不处理
	if c.Update.ChannelPost.ForwardFrom != nil {
		return
	}
	// 转发的消息不处理
	if c.Update.ChannelPost.ForwardFromChat != nil {
		return
	}
	// 若无 /s 命令则不处理
	if !strings.HasPrefix(c.Update.ChannelPost.Text, "/smr ") {
		return
	}

	urlString := strings.TrimSpace(strings.TrimPrefix(c.Update.ChannelPost.Text, "/smr "))
	summarization, err := h.summarizeInputURL(urlString)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	_, err = c.Bot.Request(tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    c.Update.ChannelPost.Chat.ID,
			MessageID: c.Update.ChannelPost.MessageID,
		},
		ParseMode: "HTML", // 这个结构体里面需要给 ParseMode 字段赋一个字面量为 HTML 的枚举值
		Text:      summarization,
	})
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram... %v", err)
		return
	}
}
