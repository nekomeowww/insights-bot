package discord

import (
	"context"
	"errors"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"net/url"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/samber/lo"
)

func (b *DiscordBot) smrCmd(event *events.ApplicationCommandInteractionCreate, data discord.SlashCommandInteractionData) {
	urlString := data.String("link")

	b.logger.Infof("discord: command received: /smr %s", urlString)

	var err error

	// url check
	if urlString == "" {
		err = errors.New("没有找到链接，可以发送一个有效的链接吗？用法：/smr <链接>")
	} else {
		var parsedURL *url.URL
		parsedURL, err = url.Parse(urlString)
		if err != nil {
			err = errors.New("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>")
		}
		if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
			err = errors.New("你发来的链接无法被理解，可以重新发一个试试。用法：/smr <链接>")
		}
	}

	if err != nil {
		err = event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(err.Error()).
			Build())
		if err != nil {
			b.logger.WithField("error", err.Error()).Warn("discord: failed to send error message")
		}

		return
	}

	// must reply the interaction as soon as possible
	err = event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("请稍等，量子速读中...").
		Build())
	if err != nil {
		b.logger.WithField("error", err.Error()).Warn("discord: failed to send error message")
	}

	output, err := b.smr.SummarizeInputURL(context.Background(), urlString)
	if err != nil {
		errMsg := ""
		if errors.Is(err, smr.ErrContentNotSupported) {
			errMsg = "暂时不支持量子速读这样的内容呢，可以换个别的链接试试。"
		} else if errors.Is(err, smr.ErrNetworkError) || errors.Is(err, smr.ErrRequestFailed) {
			errMsg = "量子速读的链接读取失败了哦。可以再试试？"
		} else {
			errMsg = "量子速读失败了。可以再试试？"
		}

		b.logger.WithField("error", err.Error()).Error("discord: summarization failed")

		_, err = b.botClient.Rest().CreateMessage(event.Channel().ID, discord.NewMessageCreateBuilder().
			SetContent(errMsg).
			Build())
		if err != nil {
			b.logger.WithField("error", err.Error()).Warn("discord: failed to send error message")
		}

		return
	}

	// cannot use event.CreateMessage again
	_, err = b.botClient.Rest().CreateMessage(event.Channel().ID, discord.NewMessageCreateBuilder().
		SetContent(output.FormatSummarizationAsDiscordMarkdown()).
		Build())
	if err != nil {
		b.logger.WithField("error", err.Error()).Warn("discord: failed to send summarization")
	}
}

func (b *DiscordBot) commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName() {
	case "smr":
		b.smrCmd(event, data)
	}
}
