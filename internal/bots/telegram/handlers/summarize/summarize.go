package summarize

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/imroc/req/v3"
	tokenizer "github.com/pandodao/tokenizer-go"
	"github.com/samber/lo"
	goopenai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
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

	c.Update.ChannelPost.Text = strings.TrimPrefix(c.Update.ChannelPost.Text, "/smr ")
	article, err := extractContentFromURL(c.Update.ChannelPost.Text)
	if err != nil {
		h.Logger.Errorf("failed to parse %s, %v\n", c.Update.ChannelPost.Text, err)
		return
	}

	textContent, err := truncateContentBasedOnTokens(article.TextContent)
	if err != nil {
		h.Logger.Errorf("failed to truncate content based on tokens... %v", err)
		return
	}

	h.Logger.WithFields(logrus.Fields{
		"title": article.Title,
		"url":   c.Update.ChannelPost.Text,
	}).Infof("summarizing article...")
	resp, err := h.OpenAI.SummarizeWithQuestionsAsSimplifiedChinese(
		article.Title,
		article.Byline,
		textContent,
	)
	if err != nil {
		h.Logger.Errorf("failed to create chat completion for summarizing... %v", err)
		return
	}

	h.Logger.WithFields(logrus.Fields{
		"title": article.Title,
		"url":   c.Update.ChannelPost.Text,
	}).Infof("summarizing article done")
	respMessages := lo.Map(resp.Choices, func(item goopenai.ChatCompletionChoice, _ int) string {
		return item.Message.Content
	})

	_, err = c.Bot.Request(tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    c.Update.ChannelPost.Chat.ID,
			MessageID: c.Update.ChannelPost.MessageID,
		},
		Text: fmt.Sprintf("%s\n\n摘要：\n%s", c.Update.ChannelPost.Text, strings.Join(respMessages, "\n")),
	})
	if err != nil {
		h.Logger.Errorf("failed to send message to telegram... %v", err)
		return
	}
}

func extractContentFromURL(urlString string) (*readability.Article, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	urlContent, err := readability.FromURL(parsedURL.String(), 1*time.Minute)
	if err != nil {
		return nil, err
	}

	return &urlContent, nil
}

// truncateContentBasedOnTokens 基于 token 计算的方式截断文本
func truncateContentBasedOnTokens(textContent string) (string, error) {
	tokens, err := tokenizer.CalToken(textContent)
	if err != nil {
		return "", err
	}
	if tokens > 3900 {
		return string([]rune(textContent)[:int(math.Min(3900, float64(len(textContent))))]), nil
	}

	return textContent, nil
}
