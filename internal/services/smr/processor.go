package smr

import (
	"context"
	"errors"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/slack-go/slack"
	"time"
)

func (s *Service) processOutput(info types.TaskInfo, result *smr.URLSummarizationOutput) {
	switch info.Platform {
	case smr.FromPlatformTelegram:
		s.sendFinalResult(info, result.FormatSummarizationAsHTML(), true)
	case smr.FromPlatformSlack:
		s.sendFinalResult(info, result.FormatSummarizationAsSlackMarkdown(), true)
	case smr.FromPlatformDiscord:
		s.sendFinalResult(info, result.FormatSummarizationAsDiscordMarkdown(), true)
	}
}

func (s *Service) processError(info types.TaskInfo, err error) {
	errMsg := ""
	if errors.Is(err, smr.ErrContentNotSupported) {
		errMsg = "暂时不支持量子速读这样的内容呢，可以换个别的链接试试。"
	} else if errors.Is(err, smr.ErrNetworkError) || errors.Is(err, smr.ErrRequestFailed) {
		errMsg = "量子速读的链接读取失败了哦。可以再试试？"
	} else {
		errMsg = "量子速读失败了。可以再试试？"
	}

	s.sendFinalResult(info, errMsg, false)
}

func (s *Service) sendFinalResult(info types.TaskInfo, result string, ok bool) {
	switch info.Platform {
	case smr.FromPlatformTelegram:
		msgEdit := tgbotapi.EditMessageTextConfig{
			Text: result,
		}
		msgEdit.ChatID = info.ChatID
		msgEdit.MessageID = info.MessageID
		if ok {
			msgEdit.ParseMode = tgbotapi.ModeHTML
		}

		_, err := s.tgBot.Send(msgEdit)
		if err != nil {
			s.logger.WithError(err).WithField("platform", info.Platform).Warn("smr service: failed to send result message")
		}
	case smr.FromPlatformSlack:
		var token *ent.SlackOAuthCredentials
		token, err := s.ent.SlackOAuthCredentials.Query().
			Where(slackoauthcredentials.TeamID(info.TeamID)).
			First(context.Background())
		if err != nil {
			s.logger.WithError(err).WithField("platform", info.Platform).Warn("smr service: failed to get team's access token")
			return
		}

		slackCfg := s.config.Slack
		slackCli := slackbot.NewSlackCli(nil, slackCfg.ClientID, slackCfg.ClientSecret, token.RefreshToken, token.AccessToken)
		_, _, _, err = slackCli.SendMessageWithTokenExpirationCheck(
			info.ChannelID,
			s.slackBot.GetService().NewStoreFuncForRefresh(info.TeamID),
			slack.MsgOptionText(result, false),
		)
		if err != nil {
			s.logger.WithError(err).WithField("platform", info.Platform).Warn("smr service: failed to send result message")
		}
	case smr.FromPlatformDiscord:
		channelID, _ := snowflake.Parse(info.ChannelID)
		_, err := s.discordBot.Rest().
			CreateMessage(channelID, discord.NewMessageCreateBuilder().
				SetContent(result).
				Build(),
			)
		if err != nil {
			s.logger.WithError(err).WithField("platform", info.Platform).Warn("smr service: failed to send result message")
		}
	}
}

func (s *Service) botExists(platform smr.FromPlatform) bool {
	switch platform {
	case smr.FromPlatformTelegram:
		return s.tgBot != nil
	case smr.FromPlatformSlack:
		return s.slackBot != nil
	case smr.FromPlatformDiscord:
		return s.discordBot != nil
	}

	return false
}

func (s *Service) processor(info types.TaskInfo) {
	if !s.botExists(info.Platform) {
		s.logger.Errorf("received task from platform %v but instance not exists", info.Platform)
		// move back to queue
		err := s.AddTask(info)
		if err != nil {
			s.logger.WithError(err).Errorf("failed to move task back to queue")
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	result, err := s.model.SummarizeInputURL(ctx, info.Url, info.Platform)
	if err != nil {
		s.logger.WithError(err).Warn("smr service: summarization failed")
		s.processError(info, err)
		return
	}

	s.processOutput(info, result)
}
