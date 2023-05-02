package slack

import (
	"errors"

	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/slack-go/slack"
)

func (b *SlackBot) smr(info recivedCommandInfo) {
	summarization, err := b.smrModel.SummarizeInputURL(info.Text)

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

		_, _, _, err = b.slackCli.SendMessage(info.ChannelId, slack.MsgOptionText(errMsg, false))
		if err != nil {
			b.logger.WithField("error", err.Error()).Warn("slack: failed to send error message")
		}
		return
	}

	_, _, _, err = b.slackCli.SendMessage(info.ChannelId, slack.MsgOptionText(summarization.FormatSummarizationAsSlackMarkdown(), false))
	if err != nil {
		b.logger.WithField("error", err.Error()).Warn("slack: failed to send summarization")
	}
}

func (b *SlackBot) runSmr() {
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
