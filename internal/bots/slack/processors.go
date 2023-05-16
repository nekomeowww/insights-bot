package slack

import (
	"context"
	"errors"
	"time"

	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/slack-go/slack"
)

func (b *Bot) smr(info smrRequestInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	summarization, err := b.smrModel.SummarizeInputURL(ctx, info.inputUrl)
	slackCfg := b.config.Slack
	slackCli := slackbot.NewSlackCli(nil, slackCfg.ClientID, slackCfg.ClientSecret, info.refreshToken, info.accessToken)
	tokenStoreFunc := b.newStoreFuncForRefresh(info.teamID)

	if err != nil {
		errMsg := ""
		if errors.Is(err, smr.ErrContentNotSupported) {
			errMsg = "暂时不支持量子速读这样的内容呢，可以换个别的链接试试。"
		} else if errors.Is(err, smr.ErrNetworkError) || errors.Is(err, smr.ErrRequestFailed) {
			errMsg = "量子速读的链接读取失败了哦。可以再试试？"
		} else {
			errMsg = "量子速读失败了。可以再试试？"
		}

		b.logger.WithField("error", err.Error()).Error("slack: summarization failed")

		_, _, _, err = slackCli.SendMessageWithTokenExpirationCheck(
			info.channelID,
			tokenStoreFunc,
			slack.MsgOptionText(errMsg, false),
		)
		if err != nil {
			b.logger.WithField("error", err.Error()).Warn("slack: failed to send error message")
		}

		return
	}

	_, _, _, err = slackCli.SendMessageWithTokenExpirationCheck(
		info.channelID,
		tokenStoreFunc,
		slack.MsgOptionText(summarization.FormatSummarizationAsSlackMarkdown(), false),
	)
	if err != nil {
		b.logger.WithField("error", err.Error()).Warn("slack: failed to send summarization")
	}
}

func (b *Bot) runSmr() {
	needToClose := false

	for {
		select {
		case <-b.closeChan:
			b.logger.WithField("last tasks", len(b.processChan)).Info("slack: received stop signal, waiting for all tasks done")
			needToClose = true
		case info := <-b.processChan:
			b.smr(info)
		}
		if needToClose && len(b.processChan) == 0 {
			b.logger.Info("slack: all tasks done")
			break
		}
	}
}
